package pricelevel_test

import (
	"testing"

	"exchange/engine/order"
	"exchange/engine/orderbook/pricelevel"
)

func Test_Remove(t *testing.T) {
	testCases := []struct {
		name           string
		insertOrders   []*order.Order
		removeOrderIDs []string
		wantErr        bool
	}{
		{
			name:           "remove_empty",
			removeOrderIDs: []string{"1"},
			wantErr:        true,
		},
		{
			name: "remove_unique",
			insertOrders: []*order.Order{
				{ID: "1", Volume: 10, Price: 1},
			},
			removeOrderIDs: []string{"1"},
		},
		{
			name: "remove_front",
			insertOrders: []*order.Order{
				{ID: "1", Volume: 10, Price: 1},
				{ID: "2", Volume: 10, Price: 1},
				{ID: "3", Volume: 10, Price: 1},
			},
			removeOrderIDs: []string{"1"},
		},
		{
			name: "remove_middle",
			insertOrders: []*order.Order{
				{ID: "1", Volume: 10, Price: 1},
				{ID: "2", Volume: 10, Price: 1},
				{ID: "3", Volume: 10, Price: 1},
			},
			removeOrderIDs: []string{"2"},
		},
		{
			name: "remove_back",
			insertOrders: []*order.Order{
				{ID: "1", Volume: 10, Price: 1},
				{ID: "2", Volume: 10, Price: 1},
				{ID: "3", Volume: 10, Price: 1},
			},
			removeOrderIDs: []string{"3"},
		},
		{
			name: "remove_unknown",
			insertOrders: []*order.Order{
				{ID: "1", Volume: 10, Price: 1},
				{ID: "2", Volume: 10, Price: 1},
				{ID: "3", Volume: 10, Price: 1},
			},
			removeOrderIDs: []string{"10"},
			wantErr:        true,
		},
		{
			name: "remove_twice",
			insertOrders: []*order.Order{
				{ID: "1", Volume: 10, Price: 1},
				{ID: "2", Volume: 10, Price: 1},
				{ID: "3", Volume: 10, Price: 1},
			},
			removeOrderIDs: []string{"2", "2"},
			wantErr:        true,
		},
		{
			name: "remove_all",
			insertOrders: []*order.Order{
				{ID: "1", Volume: 10, Price: 1},
				{ID: "2", Volume: 10, Price: 1},
				{ID: "3", Volume: 10, Price: 1},
			},
			removeOrderIDs: []string{"1", "2", "3"},
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

			var lastError error = nil
			for _, id := range tc.removeOrderIDs {
				if err := p.Remove(id); err != nil {
					lastError = err
				}
			}

			if lastError != nil && !tc.wantErr {
				t.Errorf("Remove() unexpected error: %v", lastError)
			}

			if lastError == nil && tc.wantErr {
				t.Error("Remove() expected error, got nil")
			}
		})
	}
}
