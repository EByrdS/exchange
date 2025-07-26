package engineserver

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/twmb/franz-go/pkg/kgo"
	"google.golang.org/protobuf/proto"

	"exchange/engine/market"
	"exchange/engine/order"

	exchangepb "exchange/api/v1"
	enginepb "exchange/engine/api/v1"
)

func (e *Engine) processOrderRequest(record *kgo.Record) error {
	m, ok := e.pairs.Load(record.Topic)
	if !ok {
		return errors.New("market not found")
	}

	msg := &enginepb.OrderRequest{}
	err := proto.Unmarshal(record.Value, msg)
	if err != nil {
		return err
	}

	orderSide := order.OrderBuy
	if msg.Order.Side == exchangepb.Side_SELL {
		orderSide = order.OrderSell
	}

	o := &order.Order{
		ID:     msg.Order.Id,
		Pair:   msg.Order.Pair,
		Side:   orderSide,
		Price:  msg.Order.Price,
		Volume: msg.Order.Volume,
	}

	market := m.(*market.Market)
	switch msg.Type {
	case enginepb.OrderRequest_LIMIT:
		if err := market.InsertMakerOrder(o); err != nil {
			return err
		}
	case enginepb.OrderRequest_MARKET:
		if err := market.MatchTakerOrder(o); err != nil {
			return err
		}
	case enginepb.OrderRequest_CANCEL:
		if err := market.Cancel(o); err != nil {
			return err
		}
	default:
		return errors.New(fmt.Sprintf("unhandled order request type %v, from %+v", msg.Type, msg))
	}

	return nil
}

func (e *Engine) Listen(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			fetches := e.kafka.PollFetches(ctx)
			if errs := fetches.Errors(); len(errs) > 0 {
				// Log errors but continue processing
				for _, err := range errs {
					log.Printf("Error polling Kafka: %v", err)
				}
				continue
			}

			fetches.EachRecord(func(record *kgo.Record) {
				if err := e.processOrderRequest(record); err != nil {
					log.Printf("Error processing record: %v", err)
				}
			})
		}
	}
}
