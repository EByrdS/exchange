package orderbook_test

import (
	"sync"
	"testing"

	"exchange/engine/order"
	"exchange/engine/orderbook"
	"exchange/engine/orderbook/rbtree"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func Test_Insert_Snapshot(t *testing.T) {
	testCases := []struct {
		name         string
		insertions   []*order.Order
		wantSnapshot map[uint64]uint64
	}{
		{
			name: "one_order",
			insertions: []*order.Order{
				{ID: "1", Price: 1, Volume: 3, Side: order.OrderBuy},
			},
			wantSnapshot: map[uint64]uint64{
				1: 3,
			},
		},
		{
			name: "two_orders_same_price",
			insertions: []*order.Order{
				{ID: "1", Price: 1, Volume: 3, Side: order.OrderBuy},
				{ID: "2", Price: 1, Volume: 2, Side: order.OrderBuy},
			},
			wantSnapshot: map[uint64]uint64{
				1: 5,
			},
		},
		{
			name: "two_orders_different_price",
			insertions: []*order.Order{
				{ID: "1", Price: 1, Volume: 3, Side: order.OrderBuy},
				{ID: "2", Price: 2, Volume: 3, Side: order.OrderBuy},
			},
			wantSnapshot: map[uint64]uint64{
				1: 3,
				2: 3,
			},
		},
		{
			name: "multiple_orders_same_price",
			insertions: []*order.Order{
				{ID: "1", Price: 1, Volume: 1, Side: order.OrderBuy},
				{ID: "2", Price: 1, Volume: 2, Side: order.OrderBuy},
				{ID: "3", Price: 1, Volume: 3, Side: order.OrderBuy},
				{ID: "4", Price: 1, Volume: 4, Side: order.OrderBuy},
				{ID: "5", Price: 1, Volume: 5, Side: order.OrderBuy},
			},
			wantSnapshot: map[uint64]uint64{
				1: 15,
			},
		},
		{
			name: "multiple_orders_different_prices",
			insertions: []*order.Order{
				{ID: "1", Price: 1, Volume: 1, Side: order.OrderBuy},
				{ID: "2", Price: 2, Volume: 2, Side: order.OrderBuy},
				{ID: "3", Price: 3, Volume: 3, Side: order.OrderBuy},
				{ID: "4", Price: 4, Volume: 4, Side: order.OrderBuy},
				{ID: "5", Price: 5, Volume: 5, Side: order.OrderBuy},
			},
			wantSnapshot: map[uint64]uint64{
				1: 1,
				2: 2,
				3: 3,
				4: 4,
				5: 5,
			},
		},
		{
			name: "combinations",
			insertions: []*order.Order{
				{ID: "1", Price: 1, Volume: 1, Side: order.OrderBuy},
				{ID: "2", Price: 1, Volume: 2, Side: order.OrderBuy},
				{ID: "3", Price: 1, Volume: 3, Side: order.OrderBuy},
				{ID: "4", Price: 2, Volume: 1, Side: order.OrderBuy},
				{ID: "5", Price: 2, Volume: 2, Side: order.OrderBuy},
				{ID: "6", Price: 2, Volume: 3, Side: order.OrderBuy},
				{ID: "7", Price: 3, Volume: 1, Side: order.OrderBuy},
				{ID: "8", Price: 3, Volume: 2, Side: order.OrderBuy},
				{ID: "9", Price: 3, Volume: 3, Side: order.OrderBuy},
				{ID: "0", Price: 4, Volume: 1, Side: order.OrderBuy},
			},
			wantSnapshot: map[uint64]uint64{
				1: 6,
				2: 6,
				3: 6,
				4: 1,
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

			for _, o := range tc.insertions {
				if err := b.Insert(o); err != nil {
					t.Fatalf("Insert(%v) unexpected error: %v", o, err)
				}
			}

			if diff := cmp.Diff(tc.wantSnapshot, b.Snapshot()); diff != "" {
				t.Errorf("Snapshort() diff (-want, +got):\n%s", diff)
			}
		})
	}
}

func Test_Delete_Snapshot(t *testing.T) {
	testCases := []struct {
		name         string
		insertions   []*order.Order
		deletions    []*order.Order
		wantSnapshot map[uint64]uint64
	}{
		{
			name: "empty",
			insertions: []*order.Order{
				{ID: "1", Price: 1, Volume: 3, Side: order.OrderBuy},
			},
			deletions: []*order.Order{
				{ID: "1", Price: 1, Volume: 3, Side: order.OrderBuy},
			},
		},
		{
			name: "one_order",
			insertions: []*order.Order{
				{ID: "1", Price: 1, Volume: 3, Side: order.OrderBuy},
				{ID: "2", Price: 1, Volume: 3, Side: order.OrderBuy},
			},
			deletions: []*order.Order{
				{ID: "2", Price: 1, Volume: 3, Side: order.OrderBuy},
			},
			wantSnapshot: map[uint64]uint64{
				1: 3,
			},
		},
		{
			name: "two_orders_same_price",
			insertions: []*order.Order{
				{ID: "1", Price: 1, Volume: 3, Side: order.OrderBuy},
				{ID: "2", Price: 1, Volume: 2, Side: order.OrderBuy},
				{ID: "3", Price: 1, Volume: 1, Side: order.OrderBuy},
				{ID: "4", Price: 1, Volume: 5, Side: order.OrderBuy},
			},
			deletions: []*order.Order{
				{ID: "1", Price: 1, Volume: 3, Side: order.OrderBuy},
				{ID: "4", Price: 1, Volume: 5, Side: order.OrderBuy},
			},
			wantSnapshot: map[uint64]uint64{
				1: 3,
			},
		},
		{
			name: "two_orders_different_price",
			insertions: []*order.Order{
				{ID: "1", Price: 1, Volume: 3, Side: order.OrderBuy},
				{ID: "2", Price: 2, Volume: 3, Side: order.OrderBuy},
				{ID: "3", Price: 3, Volume: 3, Side: order.OrderBuy},
				{ID: "4", Price: 4, Volume: 3, Side: order.OrderBuy},
			},
			deletions: []*order.Order{
				{ID: "1", Price: 1, Volume: 3, Side: order.OrderBuy},
				{ID: "4", Price: 4, Volume: 3, Side: order.OrderBuy},
			},
			wantSnapshot: map[uint64]uint64{
				2: 3,
				3: 3,
			},
		},
		{
			name: "multiple_orders_same_price",
			insertions: []*order.Order{
				{ID: "1", Price: 1, Volume: 1, Side: order.OrderBuy},
				{ID: "2", Price: 1, Volume: 2, Side: order.OrderBuy},
				{ID: "3", Price: 1, Volume: 3, Side: order.OrderBuy},
				{ID: "4", Price: 1, Volume: 4, Side: order.OrderBuy},
				{ID: "5", Price: 1, Volume: 5, Side: order.OrderBuy},
				{ID: "6", Price: 1, Volume: 1, Side: order.OrderBuy},
				{ID: "7", Price: 1, Volume: 2, Side: order.OrderBuy},
				{ID: "8", Price: 1, Volume: 3, Side: order.OrderBuy},
				{ID: "9", Price: 1, Volume: 4, Side: order.OrderBuy},
				{ID: "0", Price: 1, Volume: 5, Side: order.OrderBuy},
			},
			deletions: []*order.Order{
				{ID: "1", Price: 1, Volume: 1, Side: order.OrderBuy},
				{ID: "3", Price: 1, Volume: 3, Side: order.OrderBuy},
				{ID: "5", Price: 1, Volume: 5, Side: order.OrderBuy},
				{ID: "7", Price: 1, Volume: 2, Side: order.OrderBuy},
				{ID: "9", Price: 1, Volume: 4, Side: order.OrderBuy},
			},
			wantSnapshot: map[uint64]uint64{
				1: 15,
			},
		},
		{
			name: "multiple_orders_different_prices",
			insertions: []*order.Order{
				{ID: "1", Price: 1, Volume: 1, Side: order.OrderBuy},
				{ID: "2", Price: 2, Volume: 2, Side: order.OrderBuy},
				{ID: "3", Price: 3, Volume: 3, Side: order.OrderBuy},
				{ID: "4", Price: 4, Volume: 4, Side: order.OrderBuy},
				{ID: "5", Price: 5, Volume: 5, Side: order.OrderBuy},
				{ID: "6", Price: 1, Volume: 1, Side: order.OrderBuy},
				{ID: "7", Price: 2, Volume: 2, Side: order.OrderBuy},
				{ID: "8", Price: 3, Volume: 3, Side: order.OrderBuy},
				{ID: "9", Price: 4, Volume: 4, Side: order.OrderBuy},
				{ID: "0", Price: 5, Volume: 5, Side: order.OrderBuy},
			},
			deletions: []*order.Order{
				{ID: "1", Price: 1, Volume: 1, Side: order.OrderBuy},
				{ID: "3", Price: 3, Volume: 3, Side: order.OrderBuy},
				{ID: "5", Price: 5, Volume: 5, Side: order.OrderBuy},
				{ID: "7", Price: 2, Volume: 2, Side: order.OrderBuy},
				{ID: "9", Price: 4, Volume: 4, Side: order.OrderBuy},
			},
			wantSnapshot: map[uint64]uint64{
				1: 1,
				2: 2,
				3: 3,
				4: 4,
				5: 5,
			},
		},
		{
			name: "combinations",
			insertions: []*order.Order{
				{ID: "01", Price: 1, Volume: 1, Side: order.OrderBuy},
				{ID: "02", Price: 1, Volume: 2, Side: order.OrderBuy},
				{ID: "03", Price: 1, Volume: 3, Side: order.OrderBuy},
				{ID: "04", Price: 2, Volume: 1, Side: order.OrderBuy},
				{ID: "05", Price: 2, Volume: 2, Side: order.OrderBuy},
				{ID: "06", Price: 2, Volume: 3, Side: order.OrderBuy},
				{ID: "07", Price: 3, Volume: 1, Side: order.OrderBuy},
				{ID: "08", Price: 3, Volume: 2, Side: order.OrderBuy},
				{ID: "09", Price: 3, Volume: 3, Side: order.OrderBuy},
				{ID: "10", Price: 4, Volume: 1, Side: order.OrderBuy},
				{ID: "11", Price: 1, Volume: 1, Side: order.OrderBuy},
				{ID: "12", Price: 1, Volume: 2, Side: order.OrderBuy},
				{ID: "13", Price: 1, Volume: 3, Side: order.OrderBuy},
				{ID: "14", Price: 2, Volume: 1, Side: order.OrderBuy},
				{ID: "15", Price: 2, Volume: 2, Side: order.OrderBuy},
				{ID: "16", Price: 2, Volume: 3, Side: order.OrderBuy},
				{ID: "17", Price: 3, Volume: 1, Side: order.OrderBuy},
				{ID: "18", Price: 3, Volume: 2, Side: order.OrderBuy},
				{ID: "19", Price: 3, Volume: 3, Side: order.OrderBuy},
				{ID: "20", Price: 4, Volume: 1, Side: order.OrderBuy},
			},
			deletions: []*order.Order{
				{ID: "02", Price: 1, Volume: 2, Side: order.OrderBuy},
				{ID: "04", Price: 2, Volume: 1, Side: order.OrderBuy},
				{ID: "05", Price: 2, Volume: 2, Side: order.OrderBuy},
				{ID: "07", Price: 3, Volume: 1, Side: order.OrderBuy},
				{ID: "08", Price: 3, Volume: 2, Side: order.OrderBuy},
				{ID: "10", Price: 4, Volume: 1, Side: order.OrderBuy},
				{ID: "12", Price: 1, Volume: 2, Side: order.OrderBuy},
				{ID: "13", Price: 1, Volume: 3, Side: order.OrderBuy},
				{ID: "15", Price: 2, Volume: 2, Side: order.OrderBuy},
				{ID: "17", Price: 3, Volume: 1, Side: order.OrderBuy},
				{ID: "19", Price: 3, Volume: 3, Side: order.OrderBuy},
			},
			wantSnapshot: map[uint64]uint64{
				1: 5,
				2: 7,
				3: 5,
				4: 1,
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

			for _, o := range tc.insertions {
				if err := b.Insert(o); err != nil {
					t.Fatalf("Insert(%v) unexpected error: %v", o, err)
				}
			}

			for _, o := range tc.deletions {
				if err := b.Delete(o); err != nil {
					t.Fatalf("Insert(%v) unexpected error: %v", o, err)
				}
			}

			opts := cmp.Options{
				cmpopts.EquateEmpty(),
			}
			if diff := cmp.Diff(tc.wantSnapshot, b.Snapshot(), opts); diff != "" {
				t.Errorf("Snapshort() diff (-want, +got):\n%s", diff)
			}
		})
	}
}

func Test_Match_Snapshot(t *testing.T) {
	testCases := []struct {
		name         string
		side         order.OrderSide
		orders       []*order.Order
		matchVolume  uint64
		wantSnapshot map[uint64]uint64
	}{
		{
			name: "empty_match_0",
			side: order.OrderBuy,
		},
		{
			name: "non_empty_match_0",
			side: order.OrderBuy,
			orders: []*order.Order{
				{ID: "1", Price: 1, Volume: 50, Side: order.OrderBuy},
			},
			matchVolume: 0,
			wantSnapshot: map[uint64]uint64{
				1: 50,
			},
		},
		{
			name: "match_less_than",
			side: order.OrderBuy,
			orders: []*order.Order{
				{ID: "1", Price: 1, Volume: 50, Side: order.OrderBuy},
			},
			matchVolume: 25,
			wantSnapshot: map[uint64]uint64{
				1: 25,
			},
		},
		{
			name: "match_equal",
			side: order.OrderBuy,
			orders: []*order.Order{
				{ID: "1", Price: 1, Volume: 50, Side: order.OrderBuy},
			},
			matchVolume: 50,
		},
		{
			name: "match_more_than",
			side: order.OrderBuy,
			orders: []*order.Order{
				{ID: "1", Price: 1, Volume: 50, Side: order.OrderBuy},
			},
			matchVolume: 60,
		},
		{
			name: "match_two_less_than_complete",
			side: order.OrderBuy,
			orders: []*order.Order{
				{ID: "1", Price: 1, Volume: 25, Side: order.OrderBuy},
				{ID: "2", Price: 1, Volume: 25, Side: order.OrderBuy},
			},
			matchVolume: 25,
			wantSnapshot: map[uint64]uint64{
				1: 25,
			},
		},
		{
			name: "match_two_less_than_partial",
			side: order.OrderBuy,
			orders: []*order.Order{
				{ID: "1", Price: 1, Volume: 25, Side: order.OrderBuy},
				{ID: "2", Price: 1, Volume: 25, Side: order.OrderBuy},
			},
			matchVolume: 35,
			wantSnapshot: map[uint64]uint64{
				1: 15,
			},
		},
		{
			name: "match_two_equal",
			side: order.OrderBuy,
			orders: []*order.Order{
				{ID: "1", Price: 1, Volume: 25, Side: order.OrderBuy},
				{ID: "2", Price: 1, Volume: 25, Side: order.OrderBuy},
			},
			matchVolume: 50,
		},
		{
			name: "match_two_more_than",
			side: order.OrderBuy,
			orders: []*order.Order{
				{ID: "1", Price: 1, Volume: 25, Side: order.OrderBuy},
				{ID: "2", Price: 1, Volume: 25, Side: order.OrderBuy},
			},
			matchVolume: 60,
		},
		{
			name: "match_multiple_prices_sell_side",
			side: order.OrderSell,
			orders: []*order.Order{
				{ID: "1", Price: 1, Volume: 25, Side: order.OrderSell},
				{ID: "2", Price: 1, Volume: 25, Side: order.OrderSell},
				{ID: "3", Price: 2, Volume: 25, Side: order.OrderSell},
				{ID: "4", Price: 1, Volume: 25, Side: order.OrderSell},
				{ID: "5", Price: 2, Volume: 25, Side: order.OrderSell},
				{ID: "6", Price: 1, Volume: 25, Side: order.OrderSell},
				{ID: "7", Price: 3, Volume: 25, Side: order.OrderSell},
				{ID: "8", Price: 4, Volume: 25, Side: order.OrderSell},
				{ID: "9", Price: 3, Volume: 25, Side: order.OrderSell},
			},
			matchVolume: 165,
			wantSnapshot: map[uint64]uint64{
				3: 35,
				4: 25,
			},
		},
		{
			name: "match_multiple_prices_buy_side",
			side: order.OrderBuy,
			orders: []*order.Order{
				{ID: "1", Price: 1, Volume: 25, Side: order.OrderBuy},
				{ID: "2", Price: 1, Volume: 25, Side: order.OrderBuy},
				{ID: "3", Price: 2, Volume: 25, Side: order.OrderBuy},
				{ID: "4", Price: 1, Volume: 25, Side: order.OrderBuy},
				{ID: "5", Price: 2, Volume: 25, Side: order.OrderBuy},
				{ID: "6", Price: 1, Volume: 25, Side: order.OrderBuy},
				{ID: "7", Price: 3, Volume: 25, Side: order.OrderBuy},
				{ID: "8", Price: 4, Volume: 25, Side: order.OrderBuy},
				{ID: "9", Price: 3, Volume: 25, Side: order.OrderBuy},
			},
			matchVolume: 165,
			wantSnapshot: map[uint64]uint64{
				1: 60,
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
			b := orderbook.New(tc.side, pool, func(uint64, uint64) {})

			for _, o := range tc.orders {
				if err := b.Insert(o); err != nil {
					t.Fatalf("Insert(%v) unexpected error: %v", o, err)
				}
			}

			b.MatchAndExtract(tc.matchVolume)

			opts := cmp.Options{
				cmpopts.EquateEmpty(),
			}
			if diff := cmp.Diff(tc.wantSnapshot, b.Snapshot(), opts); diff != "" {
				t.Errorf("Snapshort() diff (-want, +got):\n%s", diff)
			}
		})
	}
}
