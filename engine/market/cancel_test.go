package market_test

import (
	"exchange/engine/market"
	"exchange/engine/order"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func Test_Cancel(t *testing.T) {
	tracker := newEventsTracker(10)

	pair := "USD/GBP"

	testCases := []struct {
		name             string
		setup            []*order.Order
		cancel           *order.Order
		wantVolumeEvents []*market.VolumeEvent
		wantOrderEvents  []*market.OrderEvent
		wantErr          bool
	}{
		{
			name: "cancel_buy",
			setup: []*order.Order{
				{Pair: pair, ID: "100", Price: 10, Side: order.OrderBuy, Volume: 15},
			},
			cancel: &order.Order{Pair: pair, ID: "100", Price: 10, Side: order.OrderBuy, Volume: 15},
			wantVolumeEvents: []*market.VolumeEvent{
				{Pair: pair, Side: order.OrderBuy, Price: 10, Volume: 0, Timestamp: time.Now()},
			},
			wantOrderEvents: []*market.OrderEvent{
				{Type: market.OrderCancelled, OrderID: "100", Timestamp: time.Now()},
			},
		},
		{
			name: "cancel_sell",
			setup: []*order.Order{
				{Pair: pair, ID: "100", Price: 10, Side: order.OrderSell, Volume: 15},
			},
			cancel: &order.Order{Pair: pair, ID: "100", Price: 10, Side: order.OrderSell, Volume: 15},
			wantVolumeEvents: []*market.VolumeEvent{
				{Pair: pair, Side: order.OrderSell, Price: 10, Volume: 0, Timestamp: time.Now()},
			},
			wantOrderEvents: []*market.OrderEvent{
				{Type: market.OrderCancelled, OrderID: "100", Timestamp: time.Now()},
			},
		},
		{
			name: "cancel_decreases_volume",
			setup: []*order.Order{
				{Pair: pair, ID: "100", Price: 10, Side: order.OrderBuy, Volume: 15},
				{Pair: pair, ID: "101", Price: 10, Side: order.OrderBuy, Volume: 15},
				{Pair: pair, ID: "102", Price: 10, Side: order.OrderBuy, Volume: 15},
			},
			cancel: &order.Order{Pair: pair, ID: "100", Price: 10, Side: order.OrderBuy, Volume: 15},
			wantVolumeEvents: []*market.VolumeEvent{
				{Pair: pair, Side: order.OrderBuy, Price: 10, Volume: 30, Timestamp: time.Now()},
			},
			wantOrderEvents: []*market.OrderEvent{
				{Type: market.OrderCancelled, OrderID: "100", Timestamp: time.Now()},
			},
		},
		{
			name: "removes_volume_from_book",
			setup: []*order.Order{
				{Pair: pair, ID: "100", Price: 10, Side: order.OrderBuy, Volume: 15},
			},
			cancel: &order.Order{Pair: pair, ID: "100", Price: 10, Side: order.OrderBuy, Volume: 1}, // Less than volume in book
			wantVolumeEvents: []*market.VolumeEvent{
				{Pair: pair, Side: order.OrderBuy, Price: 10, Volume: 0, Timestamp: time.Now()},
			},
			wantOrderEvents: []*market.OrderEvent{
				{Type: market.OrderCancelled, OrderID: "100", Timestamp: time.Now()},
			},
		},
		{
			name: "unknown_price",
			setup: []*order.Order{
				{Pair: pair, ID: "100", Price: 10, Side: order.OrderBuy, Volume: 15},
			},
			cancel:  &order.Order{Pair: pair, ID: "1", Price: 8, Side: order.OrderBuy, Volume: 10},
			wantErr: true,
		},
		{
			name:    "nil_order",
			cancel:  nil,
			wantErr: true,
		},
		{
			name:    "no_order_id",
			cancel:  &order.Order{Pair: pair, Price: 10, Side: order.OrderBuy, Volume: 10},
			wantErr: true,
			wantOrderEvents: []*market.OrderEvent{
				{Type: market.OrderRejected, Timestamp: time.Now()},
			},
		},
		{
			name:    "different_pair",
			cancel:  &order.Order{Pair: "USD/BTC", ID: "1", Price: 10, Side: order.OrderBuy, Volume: 10},
			wantErr: true,
			wantOrderEvents: []*market.OrderEvent{
				{Type: market.OrderRejected, OrderID: "1", Timestamp: time.Now()},
			},
		},
		{
			name:    "zero_price",
			cancel:  &order.Order{Pair: pair, ID: "1", Price: 0, Side: order.OrderBuy, Volume: 10},
			wantErr: true,
			wantOrderEvents: []*market.OrderEvent{
				{Type: market.OrderRejected, OrderID: "1", Timestamp: time.Now()},
			},
		},
		{
			name:    "zero_volume",
			cancel:  &order.Order{Pair: pair, ID: "1", Price: 10, Side: order.OrderBuy, Volume: 0},
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

			err := m.Cancel(tc.cancel)
			if err != nil && !tc.wantErr {
				t.Errorf("Cancel(%v) unexpected error, want nil, got: %v", tc.cancel, err)
			}
			if err == nil && tc.wantErr {
				t.Errorf("Cancel(%v) unexpected error, want %v, got inl", tc.cancel, tc.wantErr)
			}

			tracker.flush()

			opts := cmp.Options{
				cmpopts.EquateApproxTime(30 * time.Second),
				cmpopts.EquateEmpty(),
			}

			if diff := cmp.Diff(tc.wantOrderEvents, tracker.orderEvents, opts); diff != "" {
				t.Errorf("Cancel order events diff:\n%s", diff)
			}

			if diff := cmp.Diff(tc.wantVolumeEvents, tracker.volumeEvents, opts); diff != "" {
				t.Errorf("Cancel volume events diff:\n%s", diff)
			}
		})
	}
}
