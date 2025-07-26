package ordersservice

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/twmb/franz-go/pkg/kgo"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"

	exchangepb "exchange/api/v1"
	enginepb "exchange/engine/api/v1"
)

type Service struct {
	exchangepb.UnimplementedOrdersServiceServer

	// temporary, use db instead
	orders map[string]*exchangepb.Order

	kafka *kgo.Client
}

func (s *Service) CreateOrder(ctx context.Context, req *exchangepb.CreateOrderRequest) (*exchangepb.Order, error) {
	fmt.Printf("CreateOrder: %+v\n", req.Order)

	parts := strings.Split(req.Order.Pair, "/")

	var orderType enginepb.OrderRequest_Type
	if req.Order.Type == exchangepb.Order_LIMIT {
		orderType = enginepb.OrderRequest_LIMIT
	} else if req.Order.Type == exchangepb.Order_MARKET {
		orderType = enginepb.OrderRequest_MARKET
	} else {
		return nil, fmt.Errorf("order type not supported: %v: %w", req.Order.Type, errors.New("Bad request"))
	}

	requestPB := &enginepb.OrderRequest{
		Type:  orderType,
		Order: req.Order,
	}

	msg, err := proto.Marshal(requestPB)
	if err != nil {
		return nil, fmt.Errorf("error serializing proto: %w", err)
	}

	r := &kgo.Record{
		Topic: fmt.Sprintf("engine.%s.%s", parts[0], parts[1]),
		Value: msg,
	}

	if err := s.kafka.ProduceSync(ctx, r).FirstErr(); err != nil {
		return nil, fmt.Errorf("error producing record: %w", err)
	}

	s.orders[req.Order.Id] = req.Order

	return req.Order, nil
}

func (s *Service) DeleteOrder(ctx context.Context, req *exchangepb.DeleteOrderRequest) (*emptypb.Empty, error) {
	fmt.Printf("DeleteOrder: %q\n", req.OrderId)

	order, ok := s.orders[req.OrderId]
	if !ok {
		return nil, fmt.Errorf("order with id %q not found: %w", req.OrderId, errors.New("not found"))
	}

	requestPB := &enginepb.OrderRequest{
		Type:  enginepb.OrderRequest_CANCEL,
		Order: order,
	}

	msg, err := proto.Marshal(requestPB)
	if err != nil {
		return nil, fmt.Errorf("error serializing proto: %w", err)
	}

	parts := strings.Split(order.Pair, "/")
	r := &kgo.Record{
		Topic: fmt.Sprintf("engine.%s.%s", parts[0], parts[1]),
		Value: msg,
	}

	if err := s.kafka.ProduceSync(ctx, r).FirstErr(); err != nil {
		return nil, fmt.Errorf("error producing record: %w", err)
	}

	delete(s.orders, order.Id)

	return &emptypb.Empty{}, nil
}

func New() (*Service, error) {
	cl, err := kgo.NewClient(
		kgo.SeedBrokers("localhost:9092"),
	)
	if err != nil {
		return nil, err
	}

	s := &Service{
		kafka:  cl,
		orders: map[string]*exchangepb.Order{},
	}

	return s, nil
}
