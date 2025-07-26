package market_test

import (
	"errors"
	"exchange/engine/market"
	"exchange/engine/order"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func Test_InsertMakerOrder(t *testing.T) {
	tracker := newEventsTracker(10)

	pair := "USD/GBP"

	testCases := []struct {
		name             string
		setup            []*order.Order
		insert           *order.Order
		wantErr          error
		wantOrderEvents  []*market.OrderEvent
		wantVolumeEvents []*market.VolumeEvent
		wantMatchEvents  []*market.MatchEvent
	}{
		{
			name:   "insert_from_empty",
			insert: &order.Order{Pair: pair, ID: "1", Price: 10, Side: order.OrderBuy, Volume: 10},
			wantOrderEvents: []*market.OrderEvent{
				{Type: market.MakerOrderInserted, OrderID: "1", Timestamp: time.Now()},
			},
			wantVolumeEvents: []*market.VolumeEvent{
				{Pair: pair, Side: order.OrderBuy, Price: 10, Volume: 10, Timestamp: time.Now()},
			},
		},
		{
			name: "increase_volume",
			setup: []*order.Order{
				{Pair: pair, ID: "100", Price: 10, Side: order.OrderBuy, Volume: 15},
			},
			insert: &order.Order{Pair: pair, ID: "1", Price: 10, Side: order.OrderBuy, Volume: 10},
			wantVolumeEvents: []*market.VolumeEvent{
				{Pair: pair, Side: order.OrderBuy, Price: 10, Volume: 25, Timestamp: time.Now()},
			},
			wantOrderEvents: []*market.OrderEvent{
				{Type: market.MakerOrderInserted, OrderID: "1", Timestamp: time.Now()},
			},
		},
		{
			name: "sell_lower_than_buy_boundary",
			setup: []*order.Order{
				{Pair: pair, ID: "100", Price: 10, Side: order.OrderBuy, Volume: 15},
			},
			insert: &order.Order{Pair: pair, ID: "1", Price: 8, Side: order.OrderSell, Volume: 10},
			wantVolumeEvents: []*market.VolumeEvent{
				{Pair: pair, Side: order.OrderBuy, Price: 10, Volume: 5, Timestamp: time.Now()},
			},
			wantMatchEvents: []*market.MatchEvent{
				{Pair: pair, TakerOrderID: "1", TakerMatchType: order.OrderFulfilled, MakerOrderID: "100", MakerMatchType: order.OrderPartiallyFulfilled, SettlementPrice: 10, MatchedVolume: 10, Timestamp: time.Now()},
			},
		},
		{
			name: "buy_higher_than_sell_boundary",
			setup: []*order.Order{
				{Pair: pair, ID: "100", Price: 10, Side: order.OrderSell, Volume: 15},
			},
			insert: &order.Order{Pair: pair, ID: "1", Price: 12, Side: order.OrderBuy, Volume: 10},
			wantVolumeEvents: []*market.VolumeEvent{
				{Pair: pair, Side: order.OrderSell, Price: 10, Volume: 5, Timestamp: time.Now()},
			},
			wantMatchEvents: []*market.MatchEvent{
				{Pair: pair, TakerOrderID: "1", TakerMatchType: order.OrderFulfilled, MakerOrderID: "100", MakerMatchType: order.OrderPartiallyFulfilled, SettlementPrice: 10, MatchedVolume: 10, Timestamp: time.Now()},
			},
		},
		{
			name: "sell_lower_than_buy_boundary_multi_match",
			setup: []*order.Order{
				{Pair: pair, ID: "100", Price: 10, Side: order.OrderBuy, Volume: 15},
				{Pair: pair, ID: "101", Price: 9, Side: order.OrderBuy, Volume: 15},
				{Pair: pair, ID: "102", Price: 8, Side: order.OrderBuy, Volume: 15},
			},
			insert: &order.Order{Pair: pair, ID: "1", Price: 7, Side: order.OrderSell, Volume: 40},
			wantVolumeEvents: []*market.VolumeEvent{
				{Pair: pair, Side: order.OrderBuy, Price: 10, Volume: 0, Timestamp: time.Now()},
				{Pair: pair, Side: order.OrderBuy, Price: 9, Volume: 0, Timestamp: time.Now()},
				{Pair: pair, Side: order.OrderBuy, Price: 8, Volume: 5, Timestamp: time.Now()},
			},
			wantMatchEvents: []*market.MatchEvent{
				{Pair: pair, TakerOrderID: "1", TakerMatchType: order.OrderPartiallyFulfilled, MakerOrderID: "100", MakerMatchType: order.OrderFulfilled, SettlementPrice: 10, MatchedVolume: 15, Timestamp: time.Now()},
				{Pair: pair, TakerOrderID: "1", TakerMatchType: order.OrderPartiallyFulfilled, MakerOrderID: "101", MakerMatchType: order.OrderFulfilled, SettlementPrice: 9, MatchedVolume: 15, Timestamp: time.Now()},
				{Pair: pair, TakerOrderID: "1", TakerMatchType: order.OrderFulfilled, MakerOrderID: "102", MakerMatchType: order.OrderPartiallyFulfilled, SettlementPrice: 8, MatchedVolume: 10, Timestamp: time.Now()},
			},
		},
		{
			name: "buy_higher_than_sell_boundary_multi_match",
			setup: []*order.Order{
				{Pair: pair, ID: "100", Price: 10, Side: order.OrderSell, Volume: 15},
				{Pair: pair, ID: "101", Price: 11, Side: order.OrderSell, Volume: 15},
				{Pair: pair, ID: "102", Price: 12, Side: order.OrderSell, Volume: 15},
			},
			insert: &order.Order{Pair: pair, ID: "1", Price: 13, Side: order.OrderBuy, Volume: 40},
			wantVolumeEvents: []*market.VolumeEvent{
				{Pair: pair, Side: order.OrderSell, Price: 10, Volume: 0, Timestamp: time.Now()},
				{Pair: pair, Side: order.OrderSell, Price: 11, Volume: 0, Timestamp: time.Now()},
				{Pair: pair, Side: order.OrderSell, Price: 12, Volume: 5, Timestamp: time.Now()},
			},
			wantMatchEvents: []*market.MatchEvent{
				{Pair: pair, TakerOrderID: "1", TakerMatchType: order.OrderPartiallyFulfilled, MakerOrderID: "100", MakerMatchType: order.OrderFulfilled, SettlementPrice: 10, MatchedVolume: 15, Timestamp: time.Now()},
				{Pair: pair, TakerOrderID: "1", TakerMatchType: order.OrderPartiallyFulfilled, MakerOrderID: "101", MakerMatchType: order.OrderFulfilled, SettlementPrice: 11, MatchedVolume: 15, Timestamp: time.Now()},
				{Pair: pair, TakerOrderID: "1", TakerMatchType: order.OrderFulfilled, MakerOrderID: "102", MakerMatchType: order.OrderPartiallyFulfilled, SettlementPrice: 12, MatchedVolume: 10, Timestamp: time.Now()},
			},
		},
		{
			name:    "nil_order",
			insert:  nil,
			wantErr: market.InvalidOrderErr,
		},
		{
			name:    "no_order_id",
			insert:  &order.Order{Pair: pair, Price: 10, Side: order.OrderBuy, Volume: 10},
			wantErr: market.InvalidOrderErr,
			wantOrderEvents: []*market.OrderEvent{
				{Type: market.OrderRejected, Timestamp: time.Now()},
			},
		},
		{
			name:    "different_pair",
			insert:  &order.Order{Pair: "USD/BTC", ID: "1", Price: 10, Side: order.OrderBuy, Volume: 10},
			wantErr: market.InvalidOrderErr,
			wantOrderEvents: []*market.OrderEvent{
				{Type: market.OrderRejected, OrderID: "1", Timestamp: time.Now()},
			},
		},
		{
			name:    "zero_price",
			insert:  &order.Order{Pair: pair, ID: "1", Price: 0, Side: order.OrderBuy, Volume: 10},
			wantErr: market.InvalidOrderErr,
			wantOrderEvents: []*market.OrderEvent{
				{Type: market.OrderRejected, OrderID: "1", Timestamp: time.Now()},
			},
		},
		{
			name:    "zero_volume",
			insert:  &order.Order{Pair: pair, ID: "1", Price: 10, Side: order.OrderBuy, Volume: 0},
			wantErr: market.InvalidOrderErr,
			wantOrderEvents: []*market.OrderEvent{
				{Type: market.OrderRejected, OrderID: "1", Timestamp: time.Now()},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := market.New(pair, tracker.orderEventsChan, tracker.volumeEventsChan, tracker.matchEventsChan)

			for _, o := range tc.setup {
				if err := m.InsertMakerOrder(o); err != nil {
					t.Fatalf("InsertMakerOrder(%v) unexpected error: %v", o, err)
				}
			}

			tracker.reset()

			err := m.InsertMakerOrder(tc.insert)
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("InsertMakerOrder unexpected error, want: %v, got: %v", tc.wantErr, err)
			}

			tracker.flush()

			opts := cmp.Options{
				cmpopts.EquateApproxTime(30 * time.Second),
				cmpopts.EquateEmpty(),
			}

			if diff := cmp.Diff(tc.wantOrderEvents, tracker.orderEvents, opts); diff != "" {
				t.Errorf("InsertMakerOrder order events diff:\n%s", diff)
			}

			if diff := cmp.Diff(tc.wantVolumeEvents, tracker.volumeEvents, opts); diff != "" {
				t.Errorf("InsertMakerOrder volume events diff:\n%s", diff)
			}

			if diff := cmp.Diff(tc.wantMatchEvents, tracker.matchEvents, opts); diff != "" {
				t.Errorf("InsertMakerOrder match events diff:\n%s", diff)
			}
		})
	}
}
