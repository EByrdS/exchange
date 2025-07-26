package market_test

import (
	"exchange/engine/market"
	"exchange/engine/order"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func Test_MatchTakerOrder(t *testing.T) {
	tracker := newEventsTracker(10)

	pair := "USD/GBP"

	testCases := []struct {
		name             string
		setup            []*order.Order
		match            *order.Order
		wantVolumeEvents []*market.VolumeEvent
		wantOrderEvents  []*market.OrderEvent
		wantMatchEvents  []*market.MatchEvent
		wantErr          bool
	}{
		{
			name: "match_buy",
			setup: []*order.Order{
				{Pair: pair, ID: "100", Price: 10, Side: order.OrderSell, Volume: 15},
				{Pair: pair, ID: "101", Price: 10, Side: order.OrderSell, Volume: 15},
				{Pair: pair, ID: "102", Price: 11, Side: order.OrderSell, Volume: 15},
				{Pair: pair, ID: "103", Price: 11, Side: order.OrderSell, Volume: 15},
				{Pair: pair, ID: "104", Price: 12, Side: order.OrderSell, Volume: 15},
			},
			match: &order.Order{Pair: pair, ID: "1", Side: order.OrderBuy, Volume: 70},
			wantVolumeEvents: []*market.VolumeEvent{
				{Pair: pair, Side: order.OrderSell, Price: 10, Volume: 0, Timestamp: time.Now()},
				{Pair: pair, Side: order.OrderSell, Price: 11, Volume: 0, Timestamp: time.Now()},
				{Pair: pair, Side: order.OrderSell, Price: 12, Volume: 5, Timestamp: time.Now()},
			},
			wantMatchEvents: []*market.MatchEvent{
				{Pair: pair, TakerOrderID: "1", TakerMatchType: order.OrderPartiallyFulfilled, MakerOrderID: "100", MakerMatchType: order.OrderFulfilled, SettlementPrice: 10, MatchedVolume: 15, Timestamp: time.Now()},
				{Pair: pair, TakerOrderID: "1", TakerMatchType: order.OrderPartiallyFulfilled, MakerOrderID: "101", MakerMatchType: order.OrderFulfilled, SettlementPrice: 10, MatchedVolume: 15, Timestamp: time.Now()},
				{Pair: pair, TakerOrderID: "1", TakerMatchType: order.OrderPartiallyFulfilled, MakerOrderID: "102", MakerMatchType: order.OrderFulfilled, SettlementPrice: 11, MatchedVolume: 15, Timestamp: time.Now()},
				{Pair: pair, TakerOrderID: "1", TakerMatchType: order.OrderPartiallyFulfilled, MakerOrderID: "103", MakerMatchType: order.OrderFulfilled, SettlementPrice: 11, MatchedVolume: 15, Timestamp: time.Now()},
				{Pair: pair, TakerOrderID: "1", TakerMatchType: order.OrderFulfilled, MakerOrderID: "104", MakerMatchType: order.OrderPartiallyFulfilled, SettlementPrice: 12, MatchedVolume: 10, Timestamp: time.Now()},
			},
		},
		{
			name: "match_sell",
			setup: []*order.Order{
				{Pair: pair, ID: "100", Price: 10, Side: order.OrderBuy, Volume: 15},
				{Pair: pair, ID: "101", Price: 10, Side: order.OrderBuy, Volume: 15},
				{Pair: pair, ID: "102", Price: 9, Side: order.OrderBuy, Volume: 15},
				{Pair: pair, ID: "103", Price: 9, Side: order.OrderBuy, Volume: 15},
				{Pair: pair, ID: "104", Price: 8, Side: order.OrderBuy, Volume: 15},
			},
			match: &order.Order{Pair: pair, ID: "1", Side: order.OrderSell, Volume: 70},
			wantVolumeEvents: []*market.VolumeEvent{
				{Pair: pair, Side: order.OrderBuy, Price: 10, Volume: 0, Timestamp: time.Now()},
				{Pair: pair, Side: order.OrderBuy, Price: 9, Volume: 0, Timestamp: time.Now()},
				{Pair: pair, Side: order.OrderBuy, Price: 8, Volume: 5, Timestamp: time.Now()},
			},
			wantMatchEvents: []*market.MatchEvent{
				{Pair: pair, TakerOrderID: "1", TakerMatchType: order.OrderPartiallyFulfilled, MakerOrderID: "100", MakerMatchType: order.OrderFulfilled, SettlementPrice: 10, MatchedVolume: 15, Timestamp: time.Now()},
				{Pair: pair, TakerOrderID: "1", TakerMatchType: order.OrderPartiallyFulfilled, MakerOrderID: "101", MakerMatchType: order.OrderFulfilled, SettlementPrice: 10, MatchedVolume: 15, Timestamp: time.Now()},
				{Pair: pair, TakerOrderID: "1", TakerMatchType: order.OrderPartiallyFulfilled, MakerOrderID: "102", MakerMatchType: order.OrderFulfilled, SettlementPrice: 9, MatchedVolume: 15, Timestamp: time.Now()},
				{Pair: pair, TakerOrderID: "1", TakerMatchType: order.OrderPartiallyFulfilled, MakerOrderID: "103", MakerMatchType: order.OrderFulfilled, SettlementPrice: 9, MatchedVolume: 15, Timestamp: time.Now()},
				{Pair: pair, TakerOrderID: "1", TakerMatchType: order.OrderFulfilled, MakerOrderID: "104", MakerMatchType: order.OrderPartiallyFulfilled, SettlementPrice: 8, MatchedVolume: 10, Timestamp: time.Now()},
			},
		},
		{
			name: "not_enough_volume",
			setup: []*order.Order{
				{Pair: pair, ID: "100", Price: 10, Side: order.OrderBuy, Volume: 15},
			},
			match: &order.Order{Pair: pair, ID: "1", Side: order.OrderSell, Volume: 20},
			wantVolumeEvents: []*market.VolumeEvent{
				{Pair: pair, Side: order.OrderBuy, Price: 10, Volume: 0, Timestamp: time.Now()},
			},
			wantOrderEvents: []*market.OrderEvent{
				{Type: market.TakerOrderUnfulfilled, OrderID: "1", Timestamp: time.Now()},
			},
			wantMatchEvents: []*market.MatchEvent{
				{Pair: pair, TakerOrderID: "1", TakerMatchType: order.OrderPartiallyFulfilled, MakerOrderID: "100", MakerMatchType: order.OrderFulfilled, SettlementPrice: 10, MatchedVolume: 15, Timestamp: time.Now()},
			},
		},
		{
			name:    "nil_order",
			match:   nil,
			wantErr: true,
		},
		{
			name:    "no_order_id",
			match:   &order.Order{Pair: pair, Side: order.OrderBuy, Volume: 10},
			wantErr: true,
			wantOrderEvents: []*market.OrderEvent{
				{Type: market.OrderRejected, Timestamp: time.Now()},
			},
		},
		{
			name:    "different_pair",
			match:   &order.Order{Pair: "USD/BTC", ID: "1", Side: order.OrderBuy, Volume: 10},
			wantErr: true,
			wantOrderEvents: []*market.OrderEvent{
				{Type: market.OrderRejected, OrderID: "1", Timestamp: time.Now()},
			},
		},
		{
			name:    "zero_volume",
			match:   &order.Order{Pair: pair, ID: "1", Side: order.OrderBuy, Volume: 0},
			wantErr: true,
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

			err := m.MatchTakerOrder(tc.match)
			if err != nil && !tc.wantErr {
				t.Errorf("MatchTakerOrder(%v) unexpected error, want nil, got: %v", tc.match, err)
			}
			if err == nil && tc.wantErr {
				t.Errorf("MatchTakerOrder(%v) unexpected error, want %v, got inl", tc.match, tc.wantErr)
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
