package pricelevel_test

import (
	"testing"

	"exchange/engine/order"
	"exchange/engine/orderbook/pricelevel"

	"github.com/google/go-cmp/cmp"
)

func Test_InsertDelete_Volume(t *testing.T) {
	testCases := []struct {
		name           string
		insertOrders   []*order.Order
		removeOrderIDs []string
		wantVolume     uint64
	}{
		{
			name: "empty",
		},
		{
			name: "insert_one",
			insertOrders: []*order.Order{
				{ID: "1", Volume: 50, Price: 1},
			},
			wantVolume: 50,
		},
		{
			name: "insert_two",
			insertOrders: []*order.Order{
				{ID: "1", Volume: 50, Price: 1},
				{ID: "2", Volume: 50, Price: 1},
			},
			wantVolume: 100,
		},
		{
			name: "insert_many",
			insertOrders: []*order.Order{
				{ID: "0", Volume: 1, Price: 1},
				{ID: "1", Volume: 2, Price: 1},
				{ID: "2", Volume: 3, Price: 1},
				{ID: "3", Volume: 4, Price: 1},
				{ID: "4", Volume: 5, Price: 1},
				{ID: "5", Volume: 6, Price: 1},
				{ID: "6", Volume: 7, Price: 1},
				{ID: "7", Volume: 8, Price: 1},
				{ID: "8", Volume: 9, Price: 1},
				{ID: "9", Volume: 10, Price: 1},
			},
			wantVolume: 55,
		},
		{
			name: "delete_back",
			insertOrders: []*order.Order{
				{ID: "1", Volume: 50, Price: 1},
				{ID: "2", Volume: 25, Price: 1},
			},
			removeOrderIDs: []string{"2"},
			wantVolume:     50,
		},
		{
			name: "delete_front",
			insertOrders: []*order.Order{
				{ID: "1", Volume: 50, Price: 1},
				{ID: "2", Volume: 25, Price: 1},
			},
			removeOrderIDs: []string{"1"},
			wantVolume:     25,
		},
		{
			name: "delete_middle",
			insertOrders: []*order.Order{
				{ID: "1", Volume: 50, Price: 1},
				{ID: "2", Volume: 25, Price: 1},
				{ID: "3", Volume: 15, Price: 1},
			},
			removeOrderIDs: []string{"2"},
			wantVolume:     65,
		},
		{
			name: "delete_many",
			insertOrders: []*order.Order{
				{ID: "0", Volume: 1, Price: 1},
				{ID: "1", Volume: 2, Price: 1},
				{ID: "2", Volume: 3, Price: 1},
				{ID: "3", Volume: 4, Price: 1},
				{ID: "4", Volume: 5, Price: 1},
			},
			removeOrderIDs: []string{"1", "0", "3"},
			wantVolume:     8,
		},
		{
			name: "delete_all",
			insertOrders: []*order.Order{
				{ID: "0", Volume: 1, Price: 1},
				{ID: "1", Volume: 2, Price: 1},
				{ID: "2", Volume: 3, Price: 1},
				{ID: "3", Volume: 4, Price: 1},
				{ID: "4", Volume: 5, Price: 1},
			},
			removeOrderIDs: []string{"1", "0", "3", "4", "2"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := pricelevel.New()

			for _, o := range tc.insertOrders {
				if err := p.Insert(o); err != nil {
					t.Fatalf("PriceLevel.Insert(%v) unexpected error: %v", o, err)
				}
			}

			for _, id := range tc.removeOrderIDs {
				if err := p.Remove(id); err != nil {
					t.Fatalf("PriceLevel.Remove(%s) unexpected error: %v", id, err)
				}
			}

			if tc.wantVolume != p.Volume() {
				t.Errorf("PriceLevel.Volume() want: %d, got: %d", tc.wantVolume, p.Volume())
			}
		})
	}
}

func Test_InsertDelete_Front(t *testing.T) {
	testCases := []struct {
		name           string
		price          uint64
		insertOrders   []*order.Order
		removeOrderIDs []string
		wantFront      *order.Order
	}{
		{
			name:  "empty",
			price: 1,
		},
		{
			name:  "insert_one",
			price: 1,
			insertOrders: []*order.Order{
				{ID: "1", Volume: 50, Price: 1},
			},
			wantFront: &order.Order{ID: "1", Volume: 50, Price: 1},
		},
		{
			name:  "insert_two",
			price: 1,
			insertOrders: []*order.Order{
				{ID: "1", Volume: 50, Price: 1},
				{ID: "2", Volume: 50, Price: 1},
			},
			wantFront: &order.Order{ID: "1", Volume: 50, Price: 1},
		},
		{
			name:  "insert_many",
			price: 1,
			insertOrders: []*order.Order{
				{ID: "0", Volume: 1, Price: 1},
				{ID: "1", Volume: 2, Price: 1},
				{ID: "2", Volume: 3, Price: 1},
				{ID: "3", Volume: 4, Price: 1},
				{ID: "4", Volume: 5, Price: 1},
				{ID: "5", Volume: 6, Price: 1},
				{ID: "6", Volume: 7, Price: 1},
				{ID: "7", Volume: 8, Price: 1},
				{ID: "8", Volume: 9, Price: 1},
				{ID: "9", Volume: 10, Price: 1},
			},
			wantFront: &order.Order{ID: "0", Volume: 1, Price: 1},
		},
		{
			name:  "delete_back",
			price: 1,
			insertOrders: []*order.Order{
				{ID: "1", Volume: 50, Price: 1},
				{ID: "2", Volume: 25, Price: 1},
			},
			removeOrderIDs: []string{"2"},
			wantFront:      &order.Order{ID: "1", Volume: 50, Price: 1},
		},
		{
			name:  "delete_front",
			price: 1,
			insertOrders: []*order.Order{
				{ID: "1", Volume: 50, Price: 1},
				{ID: "2", Volume: 25, Price: 1},
			},
			removeOrderIDs: []string{"1"},
			wantFront:      &order.Order{ID: "2", Volume: 25, Price: 1},
		},
		{
			name:  "delete_middle",
			price: 1,
			insertOrders: []*order.Order{
				{ID: "1", Volume: 50, Price: 1},
				{ID: "2", Volume: 25, Price: 1},
				{ID: "3", Volume: 15, Price: 1},
			},
			removeOrderIDs: []string{"2"},
			wantFront:      &order.Order{ID: "1", Volume: 50, Price: 1},
		},
		{
			name:  "delete_many",
			price: 1,
			insertOrders: []*order.Order{
				{ID: "0", Volume: 1, Price: 1},
				{ID: "1", Volume: 2, Price: 1},
				{ID: "2", Volume: 3, Price: 1},
				{ID: "3", Volume: 4, Price: 1},
				{ID: "4", Volume: 5, Price: 1},
			},
			removeOrderIDs: []string{"1", "0", "3"},
			wantFront:      &order.Order{ID: "2", Volume: 3, Price: 1},
		},
		{
			name:  "delete_all",
			price: 1,
			insertOrders: []*order.Order{
				{ID: "0", Volume: 1, Price: 1},
				{ID: "1", Volume: 2, Price: 1},
				{ID: "2", Volume: 3, Price: 1},
				{ID: "3", Volume: 4, Price: 1},
				{ID: "4", Volume: 5, Price: 1},
			},
			removeOrderIDs: []string{"1", "0", "3", "4", "2"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := pricelevel.New()

			for _, o := range tc.insertOrders {
				if err := p.Insert(o); err != nil {
					t.Fatalf("PriceLevel.Insert(%v) unexpected error: %v", o, err)
				}
			}

			for _, id := range tc.removeOrderIDs {
				if err := p.Remove(id); err != nil {
					t.Fatalf("PriceLevel.Remove(%s) unexpected error: %v", id, err)
				}
			}

			if diff := cmp.Diff(tc.wantFront, p.Front()); diff != "" {
				t.Errorf("PriceLevel.Front() diff (-want, +got):\n%s", diff)
			}
		})
	}
}

func Test_Match_Volume(t *testing.T) {
	testCases := []struct {
		name       string
		orders     []*order.Order
		extract    uint64
		wantVolume uint64
	}{
		{
			name: "empty_extract_zero",
		},
		{
			name:    "empty_extract_one",
			extract: 1,
		},
		{
			name: "one_extract_zero",
			orders: []*order.Order{
				{ID: "1", Volume: 1, Price: 1},
			},
			wantVolume: 1,
		},
		{
			name: "one_extract_one",
			orders: []*order.Order{
				{ID: "1", Volume: 1, Price: 1},
			},
			extract: 1,
		},
		{
			name: "one_extract_two",
			orders: []*order.Order{
				{ID: "1", Volume: 1, Price: 1},
			},
			extract: 2,
		},
		{
			name: "two_extract_one",
			orders: []*order.Order{
				{ID: "1", Volume: 1, Price: 1},
				{ID: "2", Volume: 1, Price: 1},
			},
			extract:    1,
			wantVolume: 1,
		},
		{
			name: "one_order_extract_partial",
			orders: []*order.Order{
				{ID: "1", Volume: 10, Price: 1},
			},
			extract:    5,
			wantVolume: 5,
		},
		{
			name: "two_orders_extract_partial",
			orders: []*order.Order{
				{ID: "1", Volume: 10, Price: 1},
				{ID: "2", Volume: 10, Price: 1},
			},
			extract:    15,
			wantVolume: 5,
		},
		{
			name: "two_orders_extract_first_one",
			orders: []*order.Order{
				{ID: "1", Volume: 10, Price: 1},
				{ID: "2", Volume: 10, Price: 1},
			},
			extract:    10,
			wantVolume: 10,
		},
		{
			name: "two_orders_extract_all",
			orders: []*order.Order{
				{ID: "1", Volume: 10, Price: 1},
				{ID: "2", Volume: 10, Price: 1},
			},
			extract: 20,
		},
		{
			name: "two_orders_extract_more",
			orders: []*order.Order{
				{ID: "1", Volume: 10, Price: 1},
				{ID: "2", Volume: 10, Price: 1},
			},
			extract: 30,
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
			extract:    12,
			wantVolume: 43,
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
			extract:    15,
			wantVolume: 40,
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

			p.MatchAndExtract(tc.extract)

			if tc.wantVolume != p.Volume() {
				t.Errorf("Volume() want %d, got %d", tc.wantVolume, p.Volume())
			}
		})
	}
}

func Test_Match_Front(t *testing.T) {
	testCases := []struct {
		name      string
		orders    []*order.Order
		extract   uint64
		wantFront *order.Order
	}{
		{
			name: "empty_extract_zero",
		},
		{
			name:    "empty_extract_one",
			extract: 1,
		},
		{
			name: "one_extract_zero",
			orders: []*order.Order{
				{ID: "1", Volume: 1, Price: 1},
			},
			wantFront: &order.Order{ID: "1", Volume: 1, Price: 1},
		},
		{
			name: "one_extract_one",
			orders: []*order.Order{
				{ID: "1", Volume: 1, Price: 1},
			},
			extract: 1,
		},
		{
			name: "one_extract_two",
			orders: []*order.Order{
				{ID: "1", Volume: 1, Price: 1},
			},
			extract: 2,
		},
		{
			name: "two_extract_one",
			orders: []*order.Order{
				{ID: "1", Volume: 1, Price: 1},
				{ID: "2", Volume: 1, Price: 1},
			},
			extract:   1,
			wantFront: &order.Order{ID: "2", Volume: 1, Price: 1},
		},
		{
			name: "one_order_extract_partial",
			orders: []*order.Order{
				{ID: "1", Volume: 10, Price: 1},
			},
			extract:   5,
			wantFront: &order.Order{ID: "1", Volume: 5, Price: 1},
		},
		{
			name: "two_orders_extract_partial",
			orders: []*order.Order{
				{ID: "1", Volume: 10, Price: 1},
				{ID: "2", Volume: 10, Price: 1},
			},
			extract:   15,
			wantFront: &order.Order{ID: "2", Volume: 5, Price: 1},
		},
		{
			name: "two_orders_extract_first_one",
			orders: []*order.Order{
				{ID: "1", Volume: 10, Price: 1},
				{ID: "2", Volume: 10, Price: 1},
			},
			extract:   10,
			wantFront: &order.Order{ID: "2", Volume: 10, Price: 1},
		},
		{
			name: "two_orders_extract_all",
			orders: []*order.Order{
				{ID: "1", Volume: 10, Price: 1},
				{ID: "2", Volume: 10, Price: 1},
			},
			extract: 20,
		},
		{
			name: "two_orders_extract_more",
			orders: []*order.Order{
				{ID: "1", Volume: 10, Price: 1},
				{ID: "2", Volume: 10, Price: 1},
			},
			extract: 30,
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
			extract:   12,
			wantFront: &order.Order{ID: "5", Volume: 3, Price: 1},
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
			extract:   15,
			wantFront: &order.Order{ID: "6", Volume: 6, Price: 1},
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

			p.MatchAndExtract(tc.extract)

			if diff := cmp.Diff(tc.wantFront, p.Front()); diff != "" {
				t.Errorf("PriceLevel.Front() diff (-want, +got):\n%s", diff)
			}
		})
	}
}
