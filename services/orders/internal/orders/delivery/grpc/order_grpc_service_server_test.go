package grpc

import (
	"context"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/test"
	ordersService "github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/contracts/proto/service_clients"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/shared/test_fixtures/e2e"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
	"time"
)

type OrderGrpcServiceTests struct {
	*testing.T
	*e2e.E2ETestFixture
	*OrderGrpcServiceServer
}

func TestRunner(t *testing.T) {
	test.SkipCI(t)
	fixture := e2e.NewE2ETestFixture()

	//https://pkg.go.dev/testing@master#hdr-Subtests_and_Sub_benchmarks
	t.Run("GRPC", func(t *testing.T) {
		// Before running the tests
		orderGrpcService := NewOrderGrpcService(fixture.InfrastructureConfiguration)
		ordersService.RegisterOrdersServiceServer(fixture.GrpcServer.GetCurrentGrpcServer(), orderGrpcService)

		orderGrpcServiceTests := OrderGrpcServiceTests{
			T:                      t,
			E2ETestFixture:         fixture,
			OrderGrpcServiceServer: orderGrpcService,
		}

		// Run Tests
		orderGrpcServiceTests.Test_GetOrder_By_Id()
		orderGrpcServiceTests.Test_Create_Order()
		orderGrpcServiceTests.Test_GetOrders()

		// After running the tests
		fixture.GrpcServer.GracefulShutdown()
		fixture.Cleanup()
	})
}

func (p *OrderGrpcServiceTests) Test_Create_Order() {
	req := &ordersService.CreateOrderReq{
		AccountEmail:    gofakeit.Email(),
		DeliveryAddress: gofakeit.Address().Address,
		DeliveryTime:    timestamppb.New(time.Now()),
		ShopItems: []*ordersService.ShopItem{
			{
				Quantity:    uint64(gofakeit.Number(1, 10)),
				Description: gofakeit.AdjectiveDescriptive(),
				Price:       gofakeit.Price(100, 10000),
				Title:       gofakeit.Name(),
			},
		},
	}

	res, err := p.CreateOrder(context.Background(), req)
	assert.NoError(p.T, err)
	assert.NotZero(p.T, res.OrderId)
}

func (p *OrderGrpcServiceTests) Test_GetOrder_By_Id() {
	res, err := p.GetOrderByID(context.Background(), &ordersService.GetOrderByIDReq{Id: "1b4b0599-bc3c-4c1d-94af-fd1895713620"})
	assert.NoError(p.T, err)
	assert.NotNil(p.T, res.Order)
	assert.Equal(p.T, res.Order.Id, "1b4b0599-bc3c-4c1d-94af-fd1895713620")
}

func (p *OrderGrpcServiceTests) Test_GetOrders() {
	res, err := p.GetOrders(context.Background(), &ordersService.GetOrdersReq{Size: 10, Page: 1})
	assert.NoError(p.T, err)
	assert.NotNil(p.T, res.Orders)
	assert.NotEmpty(p.T, res.Orders)
}
