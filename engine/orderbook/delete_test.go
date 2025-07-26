package orderbook_test

import (
	"exchange/engine/order"
	"exchange/engine/orderbook"
	"exchange/engine/orderbook/rbtree"
	"sync"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_Delete(t *testing.T) {
	testCases := []struct {
		name       string
		insertions []*order.Order
		deletions  []*order.Order
		wantError  bool
	}{
		{
			name: "same_side",
			insertions: []*order.Order{
				{ID: "1", Price: 1, Volume: 1, Side: order.OrderBuy},
			},
			deletions: []*order.Order{
				{ID: "1", Price: 1, Volume: 1, Side: order.OrderBuy},
			},
		},
		{
			name: "two_same_price",
			insertions: []*order.Order{
				{ID: "1", Price: 1, Volume: 1, Side: order.OrderBuy},
				{ID: "2", Price: 1, Volume: 1, Side: order.OrderBuy},
			},
			deletions: []*order.Order{
				{ID: "1", Price: 1, Volume: 1, Side: order.OrderBuy},
				{ID: "2", Price: 1, Volume: 1, Side: order.OrderBuy},
			},
		},
		{
			name: "two_different_price",
			insertions: []*order.Order{
				{ID: "1", Price: 1, Volume: 1, Side: order.OrderBuy},
				{ID: "2", Price: 2, Volume: 1, Side: order.OrderBuy},
			},
			deletions: []*order.Order{
				{ID: "1", Price: 1, Volume: 1, Side: order.OrderBuy},
				{ID: "2", Price: 2, Volume: 1, Side: order.OrderBuy},
			},
		},
		{
			name: "multiple_same_price",
			insertions: []*order.Order{
				{ID: "1", Price: 1, Volume: 1, Side: order.OrderBuy},
				{ID: "2", Price: 1, Volume: 1, Side: order.OrderBuy},
				{ID: "3", Price: 1, Volume: 1, Side: order.OrderBuy},
				{ID: "4", Price: 1, Volume: 1, Side: order.OrderBuy},
				{ID: "5", Price: 1, Volume: 1, Side: order.OrderBuy},
			},
			deletions: []*order.Order{
				{ID: "1", Price: 1, Volume: 1, Side: order.OrderBuy},
				{ID: "2", Price: 1, Volume: 1, Side: order.OrderBuy},
				{ID: "3", Price: 1, Volume: 1, Side: order.OrderBuy},
				{ID: "4", Price: 1, Volume: 1, Side: order.OrderBuy},
				{ID: "5", Price: 1, Volume: 1, Side: order.OrderBuy},
			},
		},
		{
			name: "multiple_different_price",
			insertions: []*order.Order{
				{ID: "1", Price: 1, Volume: 1, Side: order.OrderBuy},
				{ID: "2", Price: 2, Volume: 1, Side: order.OrderBuy},
				{ID: "3", Price: 3, Volume: 1, Side: order.OrderBuy},
				{ID: "4", Price: 4, Volume: 1, Side: order.OrderBuy},
				{ID: "5", Price: 5, Volume: 1, Side: order.OrderBuy},
			},
			deletions: []*order.Order{
				{ID: "1", Price: 1, Volume: 1, Side: order.OrderBuy},
				{ID: "2", Price: 2, Volume: 1, Side: order.OrderBuy},
				{ID: "3", Price: 3, Volume: 1, Side: order.OrderBuy},
				{ID: "4", Price: 4, Volume: 1, Side: order.OrderBuy},
				{ID: "5", Price: 5, Volume: 1, Side: order.OrderBuy},
			},
		},
		{
			name: "multiple_same_price_not_all",
			insertions: []*order.Order{
				{ID: "1", Price: 1, Volume: 1, Side: order.OrderBuy},
				{ID: "2", Price: 1, Volume: 1, Side: order.OrderBuy},
				{ID: "3", Price: 1, Volume: 1, Side: order.OrderBuy},
				{ID: "4", Price: 1, Volume: 1, Side: order.OrderBuy},
				{ID: "5", Price: 1, Volume: 1, Side: order.OrderBuy},
			},
			deletions: []*order.Order{
				{ID: "1", Price: 1, Volume: 1, Side: order.OrderBuy},
				{ID: "2", Price: 1, Volume: 1, Side: order.OrderBuy},
				{ID: "3", Price: 1, Volume: 1, Side: order.OrderBuy},
			},
		},
		{
			name: "multiple_different_price_not_al",
			insertions: []*order.Order{
				{ID: "1", Price: 1, Volume: 1, Side: order.OrderBuy},
				{ID: "2", Price: 1, Volume: 1, Side: order.OrderBuy},
				{ID: "3", Price: 1, Volume: 1, Side: order.OrderBuy},
				{ID: "4", Price: 2, Volume: 1, Side: order.OrderBuy},
				{ID: "5", Price: 2, Volume: 1, Side: order.OrderBuy},
			},
			deletions: []*order.Order{
				{ID: "1", Price: 1, Volume: 1, Side: order.OrderBuy},
				{ID: "5", Price: 2, Volume: 1, Side: order.OrderBuy},
			},
		},
		{
			name: "empty",
			deletions: []*order.Order{
				{ID: "1", Price: 1, Volume: 1, Side: order.OrderBuy},
			},
			wantError: true,
		},
		{
			name: "different_side",
			insertions: []*order.Order{
				{ID: "1", Price: 1, Volume: 1, Side: order.OrderBuy},
			},
			deletions: []*order.Order{
				{ID: "1", Price: 1, Volume: 1, Side: order.OrderSell},
			},
			wantError: true,
		},
		{
			name: "repeated",
			insertions: []*order.Order{
				{ID: "1", Price: 1, Volume: 1, Side: order.OrderBuy},
			},
			deletions: []*order.Order{
				{ID: "1", Price: 1, Volume: 1, Side: order.OrderBuy},
				{ID: "1", Price: 1, Volume: 1, Side: order.OrderBuy},
			},
			wantError: true,
		},
		{
			name: "unknown_id",
			insertions: []*order.Order{
				{ID: "1", Price: 1, Volume: 1, Side: order.OrderBuy},
			},
			deletions: []*order.Order{
				{ID: "2", Price: 1, Volume: 1, Side: order.OrderBuy},
			},
			wantError: true,
		},
		{
			name: "id_exists_wrong_price",
			insertions: []*order.Order{
				{ID: "1", Price: 1, Volume: 1, Side: order.OrderBuy},
			},
			deletions: []*order.Order{
				{ID: "1", Price: 2, Volume: 1, Side: order.OrderBuy},
			},
			wantError: true,
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

			for _, o := range tc.insertions {
				if err := b.Insert(o); err != nil {
					t.Fatalf("Insert() unexpected error: %v", err)
				}
			}

			var lastError error
			for _, o := range tc.deletions {
				lastError = b.Delete(o)
			}

			if lastError != nil && !tc.wantError {
				t.Errorf("Delete() unexpected error: %v", lastError)
			}

			if lastError == nil && tc.wantError {
				t.Error("Delete() expected error, got nil")
			}
		})
	}
}

func Test_Delete_VolumeUpdate(t *testing.T) {
	callback := volumeCallbackTracker{}

	testCases := []struct {
		name        string
		insertions  []*order.Order
		deletions   []*order.Order
		wantUpdates []*volumeCallbackParams
	}{
		{
			name: "reduce_volume",
			insertions: []*order.Order{
				{ID: "1", Price: 1, Volume: 2, Side: order.OrderBuy},
				{ID: "2", Price: 1, Volume: 3, Side: order.OrderBuy},
				{ID: "3", Price: 1, Volume: 4, Side: order.OrderBuy},
				{ID: "4", Price: 1, Volume: 5, Side: order.OrderBuy},
			},
			deletions: []*order.Order{
				{ID: "1", Price: 1, Volume: 2, Side: order.OrderBuy},
			},
			wantUpdates: []*volumeCallbackParams{
				{Price: 1, Volume: 12},
			},
		},
		{
			name: "multiple_updates",
			insertions: []*order.Order{
				{ID: "1", Price: 1, Volume: 2, Side: order.OrderBuy},
				{ID: "2", Price: 1, Volume: 3, Side: order.OrderBuy},
				{ID: "3", Price: 1, Volume: 4, Side: order.OrderBuy},
				{ID: "4", Price: 1, Volume: 5, Side: order.OrderBuy},
			},
			deletions: []*order.Order{
				{ID: "1", Price: 1, Volume: 2, Side: order.OrderBuy},
				{ID: "4", Price: 1, Volume: 5, Side: order.OrderBuy},
			},
			wantUpdates: []*volumeCallbackParams{
				{Price: 1, Volume: 12},
				{Price: 1, Volume: 7},
			},
		},
		{
			name: "update_to_zero",
			insertions: []*order.Order{
				{ID: "1", Price: 1, Volume: 1, Side: order.OrderBuy},
			},
			deletions: []*order.Order{
				{ID: "1", Price: 1, Volume: 1, Side: order.OrderBuy},
			},
			wantUpdates: []*volumeCallbackParams{
				{Price: 1, Volume: 0},
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

			callback.reset()

			for _, o := range tc.deletions {
				if err := b.Delete(o); err != nil {
					t.Fatalf("Delete() unexpected error: %v", err)
				}
			}

			if diff := cmp.Diff(tc.wantUpdates, callback.history); diff != "" {
				t.Errorf("Delete() updates diff (-want, +got):\n%s", diff)
			}
		})
	}
}

func Test_Delete_HeadPrice(t *testing.T) {
	testCases := []struct {
		name          string
		side          order.OrderSide
		insertions    []*order.Order
		deletions     []*order.Order
		wantHeadPrice uint64
	}{
		{
			name: "buy_delete_all",
			side: order.OrderBuy,
			insertions: []*order.Order{
				{ID: "1", Price: 1, Volume: 1, Side: order.OrderBuy},
			},
			deletions: []*order.Order{
				{ID: "1", Price: 1, Volume: 1, Side: order.OrderBuy},
			},
		},
		{
			name: "sell_delete_all",
			side: order.OrderSell,
			insertions: []*order.Order{
				{ID: "1", Price: 1, Volume: 1, Side: order.OrderSell},
			},
			deletions: []*order.Order{
				{ID: "1", Price: 1, Volume: 1, Side: order.OrderSell},
			},
		},
		{
			name: "buy_delete_head",
			side: order.OrderBuy,
			insertions: []*order.Order{
				{ID: "1", Price: 10, Volume: 1, Side: order.OrderBuy},
				{ID: "2", Price: 12, Volume: 1, Side: order.OrderBuy},
			},
			deletions: []*order.Order{
				{ID: "2", Price: 12, Volume: 1, Side: order.OrderBuy},
			},
			wantHeadPrice: 10,
		},
		{
			name: "sell_delete_head",
			side: order.OrderSell,
			insertions: []*order.Order{
				{ID: "1", Price: 10, Volume: 1, Side: order.OrderSell},
				{ID: "2", Price: 9, Volume: 1, Side: order.OrderSell},
			},
			deletions: []*order.Order{
				{ID: "2", Price: 9, Volume: 1, Side: order.OrderSell},
			},
			wantHeadPrice: 10,
		},
		{
			name: "buy_delete_back",
			side: order.OrderBuy,
			insertions: []*order.Order{
				{ID: "1", Price: 10, Volume: 1, Side: order.OrderBuy},
				{ID: "2", Price: 12, Volume: 1, Side: order.OrderBuy},
			},
			deletions: []*order.Order{
				{ID: "1", Price: 10, Volume: 1, Side: order.OrderBuy},
			},
			wantHeadPrice: 12,
		},
		{
			name: "sell_delete_back",
			side: order.OrderSell,
			insertions: []*order.Order{
				{ID: "1", Price: 10, Volume: 1, Side: order.OrderSell},
				{ID: "2", Price: 9, Volume: 1, Side: order.OrderSell},
			},
			deletions: []*order.Order{
				{ID: "1", Price: 10, Volume: 1, Side: order.OrderSell},
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

			for _, o := range tc.deletions {
				if err := b.Delete(o); err != nil {
					t.Fatalf("Delete() unexpected error: %v", err)
				}
			}

			gotHeadPrice := b.HeadPrice()

			if diff := cmp.Diff(tc.wantHeadPrice, gotHeadPrice); diff != "" {
				t.Errorf("Delete() head price diff (-want, +got):\n%s", diff)
			}
		})
	}
}
