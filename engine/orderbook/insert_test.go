package orderbook_test

import (
	"exchange/engine/order"
	"exchange/engine/orderbook"
	"exchange/engine/orderbook/rbtree"
	"sync"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_Insert(t *testing.T) {
	testCases := []struct {
		name       string
		insertions []*order.Order
		wantError  bool
	}{
		{
			name: "different_side",
			insertions: []*order.Order{
				{ID: "1", Price: 1, Volume: 1, Side: order.OrderSell},
			},
			wantError: true,
		},
		{
			name: "same_side",
			insertions: []*order.Order{
				{ID: "1", Price: 1, Volume: 1, Side: order.OrderBuy},
			},
		},
		{
			name: "duplicated_id",
			insertions: []*order.Order{
				{ID: "1", Price: 1, Volume: 1, Side: order.OrderBuy},
				{ID: "1", Price: 1, Volume: 1, Side: order.OrderBuy},
			},
			wantError: true,
		},
		{
			name: "different_prices",
			insertions: []*order.Order{
				{ID: "1", Price: 1, Volume: 1, Side: order.OrderBuy},
				{ID: "2", Price: 2, Volume: 1, Side: order.OrderBuy},
			},
		},
		{
			name: "same_prices",
			insertions: []*order.Order{
				{ID: "1", Price: 1, Volume: 1, Side: order.OrderBuy},
				{ID: "2", Price: 1, Volume: 1, Side: order.OrderBuy},
			},
		},
		{
			name: "many",
			insertions: []*order.Order{
				{ID: "1", Price: 1, Volume: 1, Side: order.OrderBuy},
				{ID: "2", Price: 2, Volume: 1, Side: order.OrderBuy},
				{ID: "3", Price: 3, Volume: 1, Side: order.OrderBuy},
				{ID: "4", Price: 1, Volume: 1, Side: order.OrderBuy},
				{ID: "5", Price: 2, Volume: 1, Side: order.OrderBuy},
				{ID: "6", Price: 3, Volume: 1, Side: order.OrderBuy},
				{ID: "7", Price: 1, Volume: 1, Side: order.OrderBuy},
				{ID: "8", Price: 2, Volume: 1, Side: order.OrderBuy},
				{ID: "9", Price: 3, Volume: 1, Side: order.OrderBuy},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			pool := &sync.Pool{
				New: func() any {
					return rbtree.NewNode()
				},
			}
			b := orderbook.New(order.OrderBuy, pool, func(uint64, uint64) {})

			var lastError error
			for _, o := range tc.insertions {
				lastError = b.Insert(o)
			}

			if lastError != nil && !tc.wantError {
				t.Errorf("Insert() unexpected error: %v", lastError)
			}

			if lastError == nil && tc.wantError {
				t.Error("Insert() expected error, got nil")
			}
		})
	}
}

func Test_Insert_VolumeUpdate(t *testing.T) {
	callback := volumeCallbackTracker{}

	testCases := []struct {
		name        string
		insertions  []*order.Order
		wantUpdates []*volumeCallbackParams
	}{
		{
			name: "single_order",
			insertions: []*order.Order{
				{ID: "1", Price: 1, Volume: 2, Side: order.OrderBuy},
			},
			wantUpdates: []*volumeCallbackParams{
				{Price: 1, Volume: 2},
			},
		},
		{
			name: "multiple_orders",
			insertions: []*order.Order{
				{ID: "1", Price: 1, Volume: 2, Side: order.OrderBuy},
				{ID: "2", Price: 2, Volume: 3, Side: order.OrderBuy},
				{ID: "3", Price: 1, Volume: 4, Side: order.OrderBuy},
				{ID: "4", Price: 2, Volume: 5, Side: order.OrderBuy},
			},
			wantUpdates: []*volumeCallbackParams{
				{Price: 1, Volume: 2},
				{Price: 2, Volume: 3},
				{Price: 1, Volume: 6},
				{Price: 2, Volume: 8},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			callback.reset()

			pool := &sync.Pool{
				New: func() any {
					return rbtree.NewNode()
				},
			}

			b := orderbook.New(order.OrderBuy, pool, callback.call)

			for _, o := range tc.insertions {
				if err := b.Insert(o); err != nil {
					t.Fatalf("Insert() unexpected error: %v", err)
				}
			}

			if diff := cmp.Diff(tc.wantUpdates, callback.history); diff != "" {
				t.Errorf("Delete() updates diff (-want, +got):\n%s", diff)
			}
		})
	}
}

func Test_Insert_HeadPrice(t *testing.T) {
	testCases := []struct {
		name          string
		side          order.OrderSide
		insertions    []*order.Order
		wantHeadPrice uint64
	}{
		{
			name: "buy_empty",
			side: order.OrderBuy,
		},
		{
			name: "sell_empty",
			side: order.OrderSell,
		},
		{
			name: "buy_single",
			side: order.OrderBuy,
			insertions: []*order.Order{
				{ID: "1", Side: order.OrderBuy, Price: 10, Volume: 10},
			},
			wantHeadPrice: 10,
		},
		{
			name: "sell_single",
			side: order.OrderSell,
			insertions: []*order.Order{
				{ID: "1", Side: order.OrderSell, Price: 10, Volume: 10},
			},
			wantHeadPrice: 10,
		},
		{
			name: "buy_more_prices",
			side: order.OrderBuy,
			insertions: []*order.Order{
				{ID: "1", Side: order.OrderBuy, Price: 10, Volume: 10},
				{ID: "2", Side: order.OrderBuy, Price: 9, Volume: 10},
			},
			wantHeadPrice: 10,
		},
		{
			name: "sell_more_prices",
			side: order.OrderSell,
			insertions: []*order.Order{
				{ID: "1", Side: order.OrderSell, Price: 10, Volume: 10},
				{ID: "2", Side: order.OrderSell, Price: 12, Volume: 10},
			},
			wantHeadPrice: 10,
		},
		{
			name: "buy_overwrite",
			side: order.OrderBuy,
			insertions: []*order.Order{
				{ID: "1", Side: order.OrderBuy, Price: 10, Volume: 10},
				{ID: "2", Side: order.OrderBuy, Price: 12, Volume: 10},
			},
			wantHeadPrice: 12,
		},
		{
			name: "sell_more_prices",
			side: order.OrderSell,
			insertions: []*order.Order{
				{ID: "1", Side: order.OrderSell, Price: 10, Volume: 10},
				{ID: "2", Side: order.OrderSell, Price: 9, Volume: 10},
			},
			wantHeadPrice: 9,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			pool := &sync.Pool{
				New: func() any {
					return rbtree.NewNode()
				},
			}

			b := orderbook.New(tc.side, pool, func(uint64, uint64) {})

			for _, o := range tc.insertions {
				if err := b.Insert(o); err != nil {
					t.Fatalf("Insert() unexpected error: %v", err)
				}
			}

			gotHeadPrice := b.HeadPrice()

			if diff := cmp.Diff(tc.wantHeadPrice, gotHeadPrice); diff != "" {
				t.Errorf("Insert() head price diff (-want, +got):\n%s", diff)
			}
		})
	}
}
