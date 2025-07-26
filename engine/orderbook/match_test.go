package orderbook_test

import (
	"exchange/engine/order"
	"exchange/engine/orderbook"
	"exchange/engine/orderbook/rbtree"
	"sync"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func Test_Match(t *testing.T) {
	testCases := []struct {
		name                string
		side                order.OrderSide
		orders              []*order.Order
		matchVolume         uint64
		wantMatches         []*order.Match
		wantUnmatchedVolume uint64
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
		},
		{
			name: "match_less_than",
			side: order.OrderBuy,
			orders: []*order.Order{
				{ID: "1", Price: 1, Volume: 50, Side: order.OrderBuy},
			},
			matchVolume: 25,
			wantMatches: []*order.Match{
				{Type: order.OrderPartiallyFulfilled, MakerOrder: &order.Order{ID: "1", Price: 1, Volume: 25, Side: order.OrderBuy}, VolumeTaken: 25},
			},
		},
		{
			name: "match_equal",
			side: order.OrderBuy,
			orders: []*order.Order{
				{ID: "1", Price: 1, Volume: 50, Side: order.OrderBuy},
			},
			matchVolume: 50,
			wantMatches: []*order.Match{
				{Type: order.OrderFulfilled, MakerOrder: &order.Order{ID: "1", Price: 1, Volume: 0, Side: order.OrderBuy}, VolumeTaken: 50},
			},
		},
		{
			name: "match_more_than",
			side: order.OrderBuy,
			orders: []*order.Order{
				{ID: "1", Price: 1, Volume: 50, Side: order.OrderBuy},
			},
			matchVolume: 60,
			wantMatches: []*order.Match{
				{Type: order.OrderFulfilled, MakerOrder: &order.Order{ID: "1", Price: 1, Volume: 0, Side: order.OrderBuy}, VolumeTaken: 50},
			},
			wantUnmatchedVolume: 10,
		},
		{
			name: "match_two_less_than_complete",
			side: order.OrderBuy,
			orders: []*order.Order{
				{ID: "1", Price: 1, Volume: 25, Side: order.OrderBuy},
				{ID: "2", Price: 1, Volume: 25, Side: order.OrderBuy},
			},
			matchVolume: 25,
			wantMatches: []*order.Match{
				{Type: order.OrderFulfilled, MakerOrder: &order.Order{ID: "1", Price: 1, Volume: 0, Side: order.OrderBuy}, VolumeTaken: 25},
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
			wantMatches: []*order.Match{
				{Type: order.OrderFulfilled, MakerOrder: &order.Order{ID: "1", Price: 1, Volume: 0, Side: order.OrderBuy}, VolumeTaken: 25},
				{Type: order.OrderPartiallyFulfilled, MakerOrder: &order.Order{ID: "2", Price: 1, Volume: 15, Side: order.OrderBuy}, VolumeTaken: 10},
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
			wantMatches: []*order.Match{
				{Type: order.OrderFulfilled, MakerOrder: &order.Order{ID: "1", Price: 1, Volume: 0, Side: order.OrderBuy}, VolumeTaken: 25},
				{Type: order.OrderFulfilled, MakerOrder: &order.Order{ID: "2", Price: 1, Volume: 0, Side: order.OrderBuy}, VolumeTaken: 25},
			},
		},
		{
			name: "match_two_more_than",
			side: order.OrderBuy,
			orders: []*order.Order{
				{ID: "1", Price: 1, Volume: 25, Side: order.OrderBuy},
				{ID: "2", Price: 1, Volume: 25, Side: order.OrderBuy},
			},
			matchVolume: 60,
			wantMatches: []*order.Match{
				{Type: order.OrderFulfilled, MakerOrder: &order.Order{ID: "1", Price: 1, Volume: 0, Side: order.OrderBuy}, VolumeTaken: 25},
				{Type: order.OrderFulfilled, MakerOrder: &order.Order{ID: "2", Price: 1, Volume: 0, Side: order.OrderBuy}, VolumeTaken: 25},
			},
			wantUnmatchedVolume: 10,
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
			wantMatches: []*order.Match{
				{Type: order.OrderFulfilled, MakerOrder: &order.Order{ID: "1", Price: 1, Volume: 0, Side: order.OrderSell}, VolumeTaken: 25},
				{Type: order.OrderFulfilled, MakerOrder: &order.Order{ID: "2", Price: 1, Volume: 0, Side: order.OrderSell}, VolumeTaken: 25},
				{Type: order.OrderFulfilled, MakerOrder: &order.Order{ID: "4", Price: 1, Volume: 0, Side: order.OrderSell}, VolumeTaken: 25},
				{Type: order.OrderFulfilled, MakerOrder: &order.Order{ID: "6", Price: 1, Volume: 0, Side: order.OrderSell}, VolumeTaken: 25},
				{Type: order.OrderFulfilled, MakerOrder: &order.Order{ID: "3", Price: 2, Volume: 0, Side: order.OrderSell}, VolumeTaken: 25},
				{Type: order.OrderFulfilled, MakerOrder: &order.Order{ID: "5", Price: 2, Volume: 0, Side: order.OrderSell}, VolumeTaken: 25},
				{Type: order.OrderPartiallyFulfilled, MakerOrder: &order.Order{ID: "7", Price: 3, Volume: 10, Side: order.OrderSell}, VolumeTaken: 15},
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
			wantMatches: []*order.Match{
				{Type: order.OrderFulfilled, MakerOrder: &order.Order{ID: "8", Price: 4, Volume: 0, Side: order.OrderBuy}, VolumeTaken: 25},
				{Type: order.OrderFulfilled, MakerOrder: &order.Order{ID: "7", Price: 3, Volume: 0, Side: order.OrderBuy}, VolumeTaken: 25},
				{Type: order.OrderFulfilled, MakerOrder: &order.Order{ID: "9", Price: 3, Volume: 0, Side: order.OrderBuy}, VolumeTaken: 25},
				{Type: order.OrderFulfilled, MakerOrder: &order.Order{ID: "3", Price: 2, Volume: 0, Side: order.OrderBuy}, VolumeTaken: 25},
				{Type: order.OrderFulfilled, MakerOrder: &order.Order{ID: "5", Price: 2, Volume: 0, Side: order.OrderBuy}, VolumeTaken: 25},
				{Type: order.OrderFulfilled, MakerOrder: &order.Order{ID: "1", Price: 1, Volume: 0, Side: order.OrderBuy}, VolumeTaken: 25},
				{Type: order.OrderPartiallyFulfilled, MakerOrder: &order.Order{ID: "2", Price: 1, Volume: 10, Side: order.OrderBuy}, VolumeTaken: 15},
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

			gotMatches, gotUnmatchedVolume := b.MatchAndExtract(tc.matchVolume)

			opts := cmp.Options{
				cmpopts.EquateEmpty(),
			}
			if diff := cmp.Diff(tc.wantMatches, gotMatches, opts); diff != "" {
				t.Errorf("Match() matches diff (-want, +got):\n%s", diff)
			}

			if gotUnmatchedVolume != tc.wantUnmatchedVolume {
				t.Errorf("Match() unmatched volum diff, want: %d, got: %d", tc.wantUnmatchedVolume, gotUnmatchedVolume)
			}
		})
	}
}

func Test_Match_VolumeUpdate(t *testing.T) {
	callback := volumeCallbackTracker{}

	testCases := []struct {
		name        string
		side        order.OrderSide
		orders      []*order.Order
		matchVolume uint64
		wantUpdates []*volumeCallbackParams
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
		},
		{
			name: "match_less_than",
			side: order.OrderBuy,
			orders: []*order.Order{
				{ID: "1", Price: 1, Volume: 50, Side: order.OrderBuy},
			},
			matchVolume: 25,
			wantUpdates: []*volumeCallbackParams{
				{Price: 1, Volume: 25},
			},
		},
		{
			name: "match_equal",
			side: order.OrderBuy,
			orders: []*order.Order{
				{ID: "1", Price: 1, Volume: 50, Side: order.OrderBuy},
			},
			matchVolume: 50,
			wantUpdates: []*volumeCallbackParams{
				{Price: 1, Volume: 0},
			},
		},
		{
			name: "match_more_than",
			side: order.OrderBuy,
			orders: []*order.Order{
				{ID: "1", Price: 1, Volume: 50, Side: order.OrderBuy},
			},
			matchVolume: 60,
			wantUpdates: []*volumeCallbackParams{
				{Price: 1, Volume: 0},
			},
		},
		{
			name: "match_two_less_than_complete",
			side: order.OrderBuy,
			orders: []*order.Order{
				{ID: "1", Price: 1, Volume: 25, Side: order.OrderBuy},
				{ID: "2", Price: 1, Volume: 25, Side: order.OrderBuy},
			},
			matchVolume: 25,
			wantUpdates: []*volumeCallbackParams{
				{Price: 1, Volume: 25},
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
			wantUpdates: []*volumeCallbackParams{
				{Price: 1, Volume: 15},
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
			wantUpdates: []*volumeCallbackParams{
				{Price: 1, Volume: 0},
			},
		},
		{
			name: "match_two_more_than",
			side: order.OrderBuy,
			orders: []*order.Order{
				{ID: "1", Price: 1, Volume: 25, Side: order.OrderBuy},
				{ID: "2", Price: 1, Volume: 25, Side: order.OrderBuy},
			},
			matchVolume: 60,
			wantUpdates: []*volumeCallbackParams{
				{Price: 1, Volume: 0},
			},
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
			wantUpdates: []*volumeCallbackParams{
				{Price: 1, Volume: 0},
				{Price: 2, Volume: 0},
				{Price: 3, Volume: 35},
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
			wantUpdates: []*volumeCallbackParams{
				{Price: 4, Volume: 0},
				{Price: 3, Volume: 0},
				{Price: 2, Volume: 0},
				{Price: 1, Volume: 60},
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
			b := orderbook.New(tc.side, pool, callback.call)

			for _, o := range tc.orders {
				if err := b.Insert(o); err != nil {
					t.Fatalf("Insert(%v) unexpected error: %v", o, err)
				}
			}

			callback.reset()

			b.MatchAndExtract(tc.matchVolume)

			if diff := cmp.Diff(tc.wantUpdates, callback.history, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("Match() callback diff (-want, +got):\n%s", diff)
			}
		})
	}
}

func Test_Match_HeadPrice(t *testing.T) {
	testCases := []struct {
		name          string
		side          order.OrderSide
		orders        []*order.Order
		matchVolume   uint64
		wantHeadPrice uint64
	}{
		{
			name: "match_less_than",
			side: order.OrderBuy,
			orders: []*order.Order{
				{ID: "1", Price: 1, Volume: 50, Side: order.OrderBuy},
			},
			matchVolume:   25,
			wantHeadPrice: 1,
		},
		{
			name: "match_equal",
			side: order.OrderBuy,
			orders: []*order.Order{
				{ID: "1", Price: 1, Volume: 50, Side: order.OrderBuy},
			},
			matchVolume:   50,
			wantHeadPrice: 0,
		},
		{
			name: "match_more_than",
			side: order.OrderBuy,
			orders: []*order.Order{
				{ID: "1", Price: 1, Volume: 50, Side: order.OrderBuy},
			},
			matchVolume:   60,
			wantHeadPrice: 0,
		},
		{
			name: "match_two_less_than_complete",
			side: order.OrderBuy,
			orders: []*order.Order{
				{ID: "1", Price: 1, Volume: 25, Side: order.OrderBuy},
				{ID: "2", Price: 1, Volume: 25, Side: order.OrderBuy},
			},
			matchVolume:   25,
			wantHeadPrice: 1,
		},
		{
			name: "match_two_less_than_partial",
			side: order.OrderBuy,
			orders: []*order.Order{
				{ID: "1", Price: 1, Volume: 25, Side: order.OrderBuy},
				{ID: "2", Price: 1, Volume: 25, Side: order.OrderBuy},
			},
			matchVolume:   35,
			wantHeadPrice: 1,
		},
		{
			name: "match_two_equal",
			side: order.OrderBuy,
			orders: []*order.Order{
				{ID: "1", Price: 1, Volume: 25, Side: order.OrderBuy},
				{ID: "2", Price: 1, Volume: 25, Side: order.OrderBuy},
			},
			matchVolume:   50,
			wantHeadPrice: 0,
		},
		{
			name: "match_two_more_than",
			side: order.OrderBuy,
			orders: []*order.Order{
				{ID: "1", Price: 1, Volume: 25, Side: order.OrderBuy},
				{ID: "2", Price: 1, Volume: 25, Side: order.OrderBuy},
			},
			matchVolume:   60,
			wantHeadPrice: 0,
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
			matchVolume:   165,
			wantHeadPrice: 3,
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
			matchVolume:   165,
			wantHeadPrice: 1,
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

			gotHeadPrice := b.HeadPrice()

			if diff := cmp.Diff(tc.wantHeadPrice, gotHeadPrice); diff != "" {
				t.Errorf("Match() head price diff (-want, +got):\n%s", diff)
			}
		})
	}
}
