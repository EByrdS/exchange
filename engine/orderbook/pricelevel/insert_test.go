package pricelevel_test

import (
	"testing"

	"exchange/engine/order"
	"exchange/engine/orderbook/pricelevel"
)

func Test_Insert(t *testing.T) {
	testCases := []struct {
		name         string
		insertOrders []*order.Order
		wantErr      bool
	}{
		{
			name: "empty",
		},
		{
			name: "insert_one",
			insertOrders: []*order.Order{
				{ID: "1", Volume: 50, Price: 1},
			},
		},
		{
			name: "insert_two",
			insertOrders: []*order.Order{
				{ID: "1", Volume: 50, Price: 1},
				{ID: "2", Volume: 50, Price: 1},
			},
		},
		{
			name: "insert_repeated_id",
			insertOrders: []*order.Order{
				{ID: "1", Volume: 50, Price: 1},
				{ID: "1", Volume: 15, Price: 1},
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := pricelevel.New()

			var lastError error = nil
			for _, order := range tc.insertOrders {
				if err := p.Insert(order); err != nil {
					lastError = err
				}
			}

			if lastError != nil && !tc.wantErr {
				t.Errorf("Insert() unexpected error: %v", lastError)
			}

			if lastError == nil && tc.wantErr {
				t.Error("Insert() expected error, got nil")
			}
		})
	}
}
