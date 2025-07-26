package pricelevel_test

import (
	"exchange/engine/order"
	"exchange/engine/orderbook/pricelevel"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func Test_MatchAndExtract(t *testing.T) {
	testCases := []struct {
		name                string
		orders              []*order.Order
		extract             uint64
		wantMatches         []*order.Match
		wantUnmatchedVolume uint64
	}{
		{
			name: "empty_extract_zero",
		},
		{
			name:                "empty_extract_one",
			extract:             1,
			wantUnmatchedVolume: 1,
		},
		{
			name: "one_extract_zero",
			orders: []*order.Order{
				{ID: "1", Volume: 1, Price: 1},
			},
		},
		{
			name: "one_extract_one",
			orders: []*order.Order{
				{ID: "1", Volume: 1, Price: 1},
			},
			extract: 1,
			wantMatches: []*order.Match{
				{
					Type:        order.OrderFulfilled,
					MakerOrder:  &order.Order{ID: "1", Volume: 0, Price: 1},
					VolumeTaken: 1,
				},
			},
		},
		{
			name: "two_extract_one",
			orders: []*order.Order{
				{ID: "1", Volume: 1, Price: 1},
				{ID: "2", Volume: 1, Price: 1},
			},
			extract: 1,
			wantMatches: []*order.Match{
				{
					Type:        order.OrderFulfilled,
					MakerOrder:  &order.Order{ID: "1", Volume: 0, Price: 1},
					VolumeTaken: 1,
				},
			},
		},
		{
			name: "one_extract_two",
			orders: []*order.Order{
				{ID: "1", Volume: 1, Price: 1},
			},
			extract: 2,
			wantMatches: []*order.Match{
				{
					Type:        order.OrderFulfilled,
					MakerOrder:  &order.Order{ID: "1", Volume: 0, Price: 1},
					VolumeTaken: 1,
				},
			},
			wantUnmatchedVolume: 1,
		},
		{
			name: "one_order_extract_partial",
			orders: []*order.Order{
				{ID: "1", Volume: 10, Price: 1},
			},
			extract: 5,
			wantMatches: []*order.Match{
				{
					Type:        order.OrderPartiallyFulfilled,
					MakerOrder:  &order.Order{ID: "1", Volume: 5, Price: 1},
					VolumeTaken: 5,
				},
			},
		},
		{
			name: "two_orders_extract_partial",
			orders: []*order.Order{
				{ID: "1", Volume: 10, Price: 1},
				{ID: "2", Volume: 10, Price: 1},
			},
			extract: 15,
			wantMatches: []*order.Match{
				{
					Type:        order.OrderFulfilled,
					MakerOrder:  &order.Order{ID: "1", Volume: 0, Price: 1},
					VolumeTaken: 10,
				},
				{
					Type:        order.OrderPartiallyFulfilled,
					MakerOrder:  &order.Order{ID: "2", Volume: 5, Price: 1},
					VolumeTaken: 5,
				},
			},
		},
		{
			name: "two_orders_extract_first_one",
			orders: []*order.Order{
				{ID: "1", Volume: 10, Price: 1},
				{ID: "2", Volume: 10, Price: 1},
			},
			extract: 10,
			wantMatches: []*order.Match{
				{
					Type:        order.OrderFulfilled,
					MakerOrder:  &order.Order{ID: "1", Volume: 0, Price: 1},
					VolumeTaken: 10,
				},
			},
		},
		{
			name: "two_orders_extract_all",
			orders: []*order.Order{
				{ID: "1", Volume: 10, Price: 1},
				{ID: "2", Volume: 10, Price: 1},
			},
			extract: 20,
			wantMatches: []*order.Match{
				{
					Type:        order.OrderFulfilled,
					MakerOrder:  &order.Order{ID: "1", Volume: 0, Price: 1},
					VolumeTaken: 10,
				},
				{
					Type:        order.OrderFulfilled,
					MakerOrder:  &order.Order{ID: "2", Volume: 0, Price: 1},
					VolumeTaken: 10,
				},
			},
		},
		{
			name: "two_orders_extract_more",
			orders: []*order.Order{
				{ID: "1", Volume: 10, Price: 1},
				{ID: "2", Volume: 10, Price: 1},
			},
			extract: 30,
			wantMatches: []*order.Match{
				{
					Type:        order.OrderFulfilled,
					MakerOrder:  &order.Order{ID: "1", Volume: 0, Price: 1},
					VolumeTaken: 10,
				},
				{
					Type:        order.OrderFulfilled,
					MakerOrder:  &order.Order{ID: "2", Volume: 0, Price: 1},
					VolumeTaken: 10,
				},
			},
			wantUnmatchedVolume: 10,
		},
		{
			name: "many_extract_partial_less_than",
			orders: []*order.Order{
				{ID: "1", Volume: 1, Price: 1},
				{ID: "2", Volume: 2, Price: 1},
				{ID: "3", Volume: 3, Price: 1},
				{ID: "4", Volume: 4, Price: 1},
				{ID: "5", Volume: 5, Price: 1},
				{ID: "6", Volume: 6, Price: 1},
				{ID: "7", Volume: 7, Price: 1},
				{ID: "8", Volume: 8, Price: 1},
				{ID: "9", Volume: 9, Price: 1},
				{ID: "10", Volume: 10, Price: 1},
			},
			extract: 12,
			wantMatches: []*order.Match{
				{Type: order.OrderFulfilled, MakerOrder: &order.Order{ID: "1", Volume: 0, Price: 1}, VolumeTaken: 1},
				{Type: order.OrderFulfilled, MakerOrder: &order.Order{ID: "2", Volume: 0, Price: 1}, VolumeTaken: 2},
				{Type: order.OrderFulfilled, MakerOrder: &order.Order{ID: "3", Volume: 0, Price: 1}, VolumeTaken: 3},
				{Type: order.OrderFulfilled, MakerOrder: &order.Order{ID: "4", Volume: 0, Price: 1}, VolumeTaken: 4},
				{Type: order.OrderPartiallyFulfilled, MakerOrder: &order.Order{ID: "5", Volume: 3, Price: 1}, VolumeTaken: 2},
			},
		},
		{
			name: "many_extract_full_orders_less_than",
			orders: []*order.Order{
				{ID: "1", Volume: 1, Price: 1},
				{ID: "2", Volume: 2, Price: 1},
				{ID: "3", Volume: 3, Price: 1},
				{ID: "4", Volume: 4, Price: 1},
				{ID: "5", Volume: 5, Price: 1},
				{ID: "6", Volume: 6, Price: 1},
				{ID: "7", Volume: 7, Price: 1},
				{ID: "8", Volume: 8, Price: 1},
				{ID: "9", Volume: 9, Price: 1},
				{ID: "10", Volume: 10, Price: 1},
			},
			extract: 15,
			wantMatches: []*order.Match{
				{Type: order.OrderFulfilled, MakerOrder: &order.Order{ID: "1", Volume: 0, Price: 1}, VolumeTaken: 1},
				{Type: order.OrderFulfilled, MakerOrder: &order.Order{ID: "2", Volume: 0, Price: 1}, VolumeTaken: 2},
				{Type: order.OrderFulfilled, MakerOrder: &order.Order{ID: "3", Volume: 0, Price: 1}, VolumeTaken: 3},
				{Type: order.OrderFulfilled, MakerOrder: &order.Order{ID: "4", Volume: 0, Price: 1}, VolumeTaken: 4},
				{Type: order.OrderFulfilled, MakerOrder: &order.Order{ID: "5", Volume: 0, Price: 1}, VolumeTaken: 5},
			},
		},
		{
			name: "many_extract_equal",
			orders: []*order.Order{
				{ID: "1", Volume: 1, Price: 1},
				{ID: "2", Volume: 2, Price: 1},
				{ID: "3", Volume: 3, Price: 1},
				{ID: "4", Volume: 4, Price: 1},
				{ID: "5", Volume: 5, Price: 1},
				{ID: "6", Volume: 6, Price: 1},
				{ID: "7", Volume: 7, Price: 1},
				{ID: "8", Volume: 8, Price: 1},
				{ID: "9", Volume: 9, Price: 1},
				{ID: "10", Volume: 10, Price: 1},
			},
			extract: 55,
			wantMatches: []*order.Match{
				{Type: order.OrderFulfilled, MakerOrder: &order.Order{ID: "1", Volume: 0, Price: 1}, VolumeTaken: 1},
				{Type: order.OrderFulfilled, MakerOrder: &order.Order{ID: "2", Volume: 0, Price: 1}, VolumeTaken: 2},
				{Type: order.OrderFulfilled, MakerOrder: &order.Order{ID: "3", Volume: 0, Price: 1}, VolumeTaken: 3},
				{Type: order.OrderFulfilled, MakerOrder: &order.Order{ID: "4", Volume: 0, Price: 1}, VolumeTaken: 4},
				{Type: order.OrderFulfilled, MakerOrder: &order.Order{ID: "5", Volume: 0, Price: 1}, VolumeTaken: 5},
				{Type: order.OrderFulfilled, MakerOrder: &order.Order{ID: "6", Volume: 0, Price: 1}, VolumeTaken: 6},
				{Type: order.OrderFulfilled, MakerOrder: &order.Order{ID: "7", Volume: 0, Price: 1}, VolumeTaken: 7},
				{Type: order.OrderFulfilled, MakerOrder: &order.Order{ID: "8", Volume: 0, Price: 1}, VolumeTaken: 8},
				{Type: order.OrderFulfilled, MakerOrder: &order.Order{ID: "9", Volume: 0, Price: 1}, VolumeTaken: 9},
				{Type: order.OrderFulfilled, MakerOrder: &order.Order{ID: "10", Volume: 0, Price: 1}, VolumeTaken: 10},
			},
		},
		{
			name: "many_extract_more_than",
			orders: []*order.Order{
				{ID: "1", Volume: 1, Price: 1},
				{ID: "2", Volume: 2, Price: 1},
				{ID: "3", Volume: 3, Price: 1},
				{ID: "4", Volume: 4, Price: 1},
				{ID: "5", Volume: 5, Price: 1},
				{ID: "6", Volume: 6, Price: 1},
				{ID: "7", Volume: 7, Price: 1},
				{ID: "8", Volume: 8, Price: 1},
				{ID: "9", Volume: 9, Price: 1},
				{ID: "10", Volume: 10, Price: 1},
			},
			extract: 65,
			wantMatches: []*order.Match{
				{Type: order.OrderFulfilled, MakerOrder: &order.Order{ID: "1", Volume: 0, Price: 1}, VolumeTaken: 1},
				{Type: order.OrderFulfilled, MakerOrder: &order.Order{ID: "2", Volume: 0, Price: 1}, VolumeTaken: 2},
				{Type: order.OrderFulfilled, MakerOrder: &order.Order{ID: "3", Volume: 0, Price: 1}, VolumeTaken: 3},
				{Type: order.OrderFulfilled, MakerOrder: &order.Order{ID: "4", Volume: 0, Price: 1}, VolumeTaken: 4},
				{Type: order.OrderFulfilled, MakerOrder: &order.Order{ID: "5", Volume: 0, Price: 1}, VolumeTaken: 5},
				{Type: order.OrderFulfilled, MakerOrder: &order.Order{ID: "6", Volume: 0, Price: 1}, VolumeTaken: 6},
				{Type: order.OrderFulfilled, MakerOrder: &order.Order{ID: "7", Volume: 0, Price: 1}, VolumeTaken: 7},
				{Type: order.OrderFulfilled, MakerOrder: &order.Order{ID: "8", Volume: 0, Price: 1}, VolumeTaken: 8},
				{Type: order.OrderFulfilled, MakerOrder: &order.Order{ID: "9", Volume: 0, Price: 1}, VolumeTaken: 9},
				{Type: order.OrderFulfilled, MakerOrder: &order.Order{ID: "10", Volume: 0, Price: 1}, VolumeTaken: 10},
			},
			wantUnmatchedVolume: 10,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := pricelevel.New()

			for _, o := range tc.orders {
				if err := p.Insert(o); err != nil {
					t.Fatalf("Insert(%v) unexpected error: %v", o, err)
				}
			}

			gotMatches, gotUnmatchedVolume := p.MatchAndExtract(tc.extract)

			opts := cmp.Options{
				cmpopts.EquateEmpty(),
			}
			if diff := cmp.Diff(tc.wantMatches, gotMatches, opts); diff != "" {
				t.Errorf("PriceLevel.MatchAndExtract(%d) Matches diff (-want, +got):\n%s", tc.extract, diff)
			}

			if gotUnmatchedVolume != tc.wantUnmatchedVolume {
				t.Errorf("PriceLevel.MatchAndExtract(%d) unmatched volume want %d, got %d", tc.extract, tc.wantUnmatchedVolume, gotUnmatchedVolume)
			}
		})
	}
}
