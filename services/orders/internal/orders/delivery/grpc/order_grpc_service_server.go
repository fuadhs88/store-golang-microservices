package grpc

import (
	"context"
	"emperror.dev/errors"
	"fmt"
	"github.com/mehdihadeli/go-mediatr"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/grpc/grpcErrors"
	customErrors "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/http_errors/custom_errors"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/mapper"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/tracing"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/utils"
	grpcOrderService "github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/contracts/proto/service_clients"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/dtos"
	creatingOrderCommandV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/features/creating_order/commands/v1"
	orderDtos "github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/features/creating_order/dtos"
	gettingOrderByIdDtos "github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/features/getting_order_by_id/dtos"
	gettingOrderByIdQueryV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/features/getting_order_by_id/queries/v1"
	gettingOrdersDtos "github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/features/getting_orders/dtos"
	gettingOrdersQueryV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/features/getting_orders/queryies/v1"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/shared/configurations/infrastructure"
	"github.com/opentracing/opentracing-go/log"
	uuid "github.com/satori/go.uuid"
)

type OrderGrpcServiceServer struct {
	*infrastructure.InfrastructureConfiguration
}

func NewOrderGrpcService(infra *infrastructure.InfrastructureConfiguration) *OrderGrpcServiceServer {
	return &OrderGrpcServiceServer{InfrastructureConfiguration: infra}
}

func (o OrderGrpcServiceServer) CreateOrder(ctx context.Context, req *grpcOrderService.CreateOrderReq) (*grpcOrderService.CreateOrderRes, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "OrderGrpcServiceServer.CreateOrder")
	span.LogFields(log.Object("Request", req))
	o.Metrics.CreateOrderGrpcRequests.Inc()
	defer span.Finish()

	shopItemsDtos, err := mapper.Map[[]*dtos.ShopItemDto](req.GetShopItems())
	if err != nil {
		return nil, err
	}

	command := creatingOrderCommandV1.NewCreateOrder(shopItemsDtos, req.AccountEmail, req.DeliveryAddress, req.DeliveryTime.AsTime())
	if err := o.Validator.StructCtx(ctx, command); err != nil {
		validationErr := customErrors.NewValidationErrorWrap(err, "[OrderGrpcServiceServer_CreateOrder.StructCtx] command validation failed")
		o.Log.Errorf(fmt.Sprintf("[OrderGrpcServiceServer_CreateOrder.StructCtx] err: %v", tracing.TraceWithErr(span, validationErr)))
		return nil, grpcErrors.ErrGrpcResponse(validationErr)
	}

	result, err := mediatr.Send[*creatingOrderCommandV1.CreateOrder, *orderDtos.CreateOrderResponseDto](ctx, command)

	if err != nil {
		err = errors.WithMessage(err, "[ProductGrpcServiceServer_CreateOrder.Send] error in sending CreateOrder")
		o.Log.Errorw(fmt.Sprintf("[ProductGrpcServiceServer_CreateOrder.Send] id: {%s}, err: %v", command.OrderId, tracing.TraceWithErr(span, err)), logger.Fields{"Id": command.OrderId})
		return nil, grpcErrors.ErrGrpcResponse(err)
	}

	o.Metrics.SuccessGrpcRequests.Inc()
	return &grpcOrderService.CreateOrderRes{OrderId: result.OrderId.String()}, nil
}

func (o OrderGrpcServiceServer) GetOrderByID(ctx context.Context, req *grpcOrderService.GetOrderByIDReq) (*grpcOrderService.GetOrderByIDRes, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "OrderGrpcServiceServer.GetOrderByID")
	span.LogFields(log.Object("Request", req))
	o.Metrics.GetOrderByIdGrpcRequests.Inc()
	defer span.Finish()

	orderIdUUID, err := uuid.FromString(req.Id)
	if err != nil {
		badRequestErr := customErrors.NewBadRequestErrorWrap(err, "[OrderGrpcServiceServer_GetOrderByID.uuid.FromString] error in converting uuid")
		o.Log.Errorf(fmt.Sprintf("[OrderGrpcServiceServer_GetOrderByID.uuid.FromString] err: %v", tracing.TraceWithErr(span, badRequestErr)))
		return nil, grpcErrors.ErrGrpcResponse(badRequestErr)
	}

	query := gettingOrderByIdQueryV1.NewGetOrderById(orderIdUUID)
	if err := o.Validator.StructCtx(ctx, query); err != nil {
		validationErr := customErrors.NewValidationErrorWrap(err, "[OrderGrpcServiceServer_GetOrderByID.StructCtx] query validation failed")
		o.Log.Errorf(fmt.Sprintf("[OrderGrpcServiceServer_GetOrderByID.StructCtx] err: %v", tracing.TraceWithErr(span, validationErr)))
		return nil, grpcErrors.ErrGrpcResponse(validationErr)
	}

	queryResult, err := mediatr.Send[*gettingOrderByIdQueryV1.GetOrderById, *gettingOrderByIdDtos.GetOrderByIdResponseDto](ctx, query)
	if err != nil {
		err = errors.WithMessage(err, "[OrderGrpcServiceServer_GetOrderByID.Send] error in sending GetOrderById")
		o.Log.Errorw(fmt.Sprintf("[OrderGrpcServiceServer_GetOrderByID.Send] id: {%s}, err: %v", query.Id, tracing.TraceWithErr(span, err)), logger.Fields{"Id": query.Id})
		return nil, grpcErrors.ErrGrpcResponse(err)
	}

	q := queryResult.Order
	order, err := mapper.Map[*grpcOrderService.OrderReadModel](q)
	if err != nil {
		err = errors.WithMessage(err, "[OrderGrpcServiceServer_GetOrderByID.Map] error in mapping order")
		return nil, grpcErrors.ErrGrpcResponse(tracing.TraceWithErr(span, err))
	}

	o.Metrics.SuccessGrpcRequests.Inc()

	return &grpcOrderService.GetOrderByIDRes{Order: order}, nil
}

func (o OrderGrpcServiceServer) SubmitOrder(ctx context.Context, req *grpcOrderService.SubmitOrderReq) (*grpcOrderService.SubmitOrderRes, error) {
	//TODO implement me
	panic("implement me")
}

func (o OrderGrpcServiceServer) UpdateShoppingCart(ctx context.Context, req *grpcOrderService.UpdateShoppingCartReq) (*grpcOrderService.UpdateShoppingCartRes, error) {
	//TODO implement me
	panic("implement me")
}

func (o OrderGrpcServiceServer) GetOrders(ctx context.Context, req *grpcOrderService.GetOrdersReq) (*grpcOrderService.GetOrdersRes, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "OrderGrpcServiceServer.GetOrders")
	span.LogFields(log.Object("Request", req))
	o.Metrics.GetOrdersGrpcRequests.Inc()
	defer span.Finish()

	query := gettingOrdersQueryV1.NewGetOrders(&utils.ListQuery{Page: int(req.Page), Size: int(req.Size)})

	queryResult, err := mediatr.Send[*gettingOrdersQueryV1.GetOrders, *gettingOrdersDtos.GetOrdersResponseDto](ctx, query)

	if err != nil {
		err = errors.WithMessage(err, "[OrderGrpcServiceServer_GetOrders.Send] error in sending GetOrders")
		o.Log.Error(fmt.Sprintf("[OrderGrpcServiceServer_GetOrders.Send] err: {%v}", tracing.TraceWithErr(span, err)))
		return nil, grpcErrors.ErrGrpcResponse(err)
	}

	ordersResponse, err := mapper.Map[*grpcOrderService.GetOrdersRes](queryResult.Orders)
	if err != nil {
		err = errors.WithMessage(err, "[OrderGrpcServiceServer_GetOrders.Map] error in mapping orders")
		return nil, grpcErrors.ErrGrpcResponse(tracing.TraceWithErr(span, err))
	}

	o.Metrics.SuccessGrpcRequests.Inc()

	return ordersResponse, nil
}
