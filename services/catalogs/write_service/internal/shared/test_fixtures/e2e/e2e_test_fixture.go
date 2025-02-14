package e2e

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/constants"
	grpcServer "github.com/mehdihadeli/store-golang-microservice-sample/pkg/grpc"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger/defaultLogger"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/configurations/mappings"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/configurations/mediatr"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/data/repositories"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/shared/configurations/infrastructure"
	"net/http/httptest"
)

type E2ETestFixture struct {
	Echo *echo.Echo
	*infrastructure.InfrastructureConfiguration
	V1         *V1Groups
	GrpcServer grpcServer.GrpcServer
	HttpServer *httptest.Server
	ctx        context.Context
	cancel     context.CancelFunc
	Cleanup    func()
}

type V1Groups struct {
	ProductsGroup *echo.Group
}

func NewE2ETestFixture() *E2ETestFixture {
	cfg, _ := config.InitConfig(constants.Test)

	ctx, cancel := context.WithCancel(context.Background())
	c := infrastructure.NewInfrastructureConfigurator(defaultLogger.Logger, cfg)
	infrastructures, _, cleanup := c.ConfigInfrastructures(context.Background())
	echo := echo.New()

	v1Group := echo.Group("/api/v1")
	productsV1 := v1Group.Group("/products")

	v1Groups := &V1Groups{ProductsGroup: productsV1}

	productRep := repositories.NewPostgresProductRepository(infrastructures.Log, cfg, infrastructures.Gorm.DB)

	err := mediatr.ConfigProductsMediator(productRep, infrastructures)
	if err != nil {
		cancel()
		return nil
	}

	err = mappings.ConfigureMappings()
	if err != nil {
		cancel()
		return nil
	}

	grpcServer := grpcServer.NewGrpcServer(cfg.GRPC, defaultLogger.Logger)
	httpServer := httptest.NewServer(echo)

	return &E2ETestFixture{
		Cleanup: func() {
			cancel()
			cleanup()
			grpcServer.GracefulShutdown()
			echo.Shutdown(ctx)
			httpServer.Close()
		},
		InfrastructureConfiguration: infrastructures,
		Echo:                        echo,
		V1:                          v1Groups,
		GrpcServer:                  grpcServer,
		HttpServer:                  httpServer,
		ctx:                         ctx,
		cancel:                      cancel,
	}
}

func (e *E2ETestFixture) Run() {
	go func() {
		if err := e.GrpcServer.RunGrpcServer(nil); err != nil {
			e.cancel()
			e.Log.Errorf("(s.RunGrpcServer) err: %v", err)
		}
	}()
}
