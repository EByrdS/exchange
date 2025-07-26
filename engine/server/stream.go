package engineserver

import (
	"context"
	"fmt"

	"github.com/twmb/franz-go/pkg/kgo"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"exchange/engine/market"
	"exchange/engine/order"

	exchangepb "exchange/api/v1"
	enginepb "exchange/engine/api/v1"
)

func (e *Engine) produce(ctx context.Context, topic string, msg []byte) {
	r := &kgo.Record{Topic: topic, Value: msg}

	e.kafka.Produce(ctx, r, func(_ *kgo.Record, err error) {
		if err != nil {
			fmt.Printf("record had a produce error: %v\n", err)
		}
	})
}

func (e *Engine) Stream(ctx context.Context) {
	for _, ms := range e.marketSymbols {
		fmt.Printf("Streaming for market: %v\n", ms)

		go func(topic string, orderEventsChan chan *market.OrderEvent) {
			for {
				ev, ok := <-orderEventsChan
				if !ok {
					return
				}

				eventType := enginepb.OrderEvent_UNDEFINED
				switch ev.Type {
				case market.OrderCancelled:
					eventType = enginepb.OrderEvent_ORDER_CANCELLED
				case market.MakerOrderInserted:
					eventType = enginepb.OrderEvent_MAKER_ORDER_INSERTED
				case market.TakerOrderUnfulfilled:
					eventType = enginepb.OrderEvent_TAKER_ORDER_UNFULFILLED
				case market.OrderRejected:
					eventType = enginepb.OrderEvent_ORDER_REJECTED
				}

				eventPB := &enginepb.OrderEvent{
					Type:    eventType,
					OrderId: ev.OrderID,
					Time:    timestamppb.New(ev.Timestamp),
				}

				msg, err := proto.Marshal(eventPB)
				if err != nil {
					fmt.Printf("Error marshalling proto: %v", err)
					continue
				}

				e.produce(ctx, topic, msg)
			}
		}(ms.Topic()+".orders", e.orderEventsChans[ms.Topic()])

		go func(topic string, volumeEventsChan chan *market.VolumeEvent) {
			for {
				ev, ok := <-volumeEventsChan
				if !ok {
					return
				}

				eventSide := exchangepb.Side_BUY
				if ev.Side == order.OrderSell {
					eventSide = exchangepb.Side_SELL
				}

				eventPB := &enginepb.VolumeEvent{
					Pair:   ev.Pair,
					Side:   eventSide,
					Price:  ev.Price,
					Volume: ev.Volume,
					Time:   timestamppb.New(ev.Timestamp),
				}

				msg, err := proto.Marshal(eventPB)
				if err != nil {
					fmt.Printf("Error marshalling proto: %v", err)
					continue
				}

				e.produce(ctx, topic, msg)
			}
		}(ms.Topic()+".volumes", e.volumeEventsChans[ms.Topic()])

		go func(topic string, matchEventsChan chan *market.MatchEvent) {
			for {
				ev, ok := <-matchEventsChan
				if !ok {
					return
				}

				takerMatchType := enginepb.MatchType_ORDER_FULFILLED
				if ev.TakerMatchType == order.OrderPartiallyFulfilled {
					takerMatchType = enginepb.MatchType_ORDER_PARTIALLY_FULFILLED
				}

				makerMatchType := enginepb.MatchType_ORDER_FULFILLED
				if ev.MakerMatchType == order.OrderPartiallyFulfilled {
					makerMatchType = enginepb.MatchType_ORDER_PARTIALLY_FULFILLED
				}

				eventPB := &enginepb.MatchEvent{
					Pair:            ev.Pair,
					TakerOrderId:    ev.TakerOrderID,
					TakerMatchType:  takerMatchType,
					MakerOrderId:    ev.MakerOrderID,
					MakerMatchType:  makerMatchType,
					MatchedVolume:   ev.MatchedVolume,
					SettlementPrice: ev.SettlementPrice,
					Time:            timestamppb.New(ev.Timestamp),
				}

				msg, err := proto.Marshal(eventPB)
				if err != nil {
					fmt.Printf("Error marshalling proto: %v", err)
					continue
				}

				e.produce(ctx, topic, msg)
			}
		}(ms.Topic()+".matches", e.matchEventsChans[ms.Topic()])
	}
}
