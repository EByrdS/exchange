package market_test

import (
	"exchange/engine/market"
	"exchange/engine/order"
	"exchange/engine/testutils"
	"testing"
)

func Benchmark_MatchTakerOrder(b *testing.B) {
	tracker := newEventsTracker(10)
	tracker.ignoreAll()

	pair := "USD/GBP"

	testCases := []struct {
		name        string
		orders      []*order.Order
		matchOrders []*order.Order
	}{
		{
			name:        "one_thousand_range_hundredth",
			orders:      testutils.OrdersDeterministic("USD/GBP", order.OrderBuy, 1, 10, 1, 10, 1_000),
			matchOrders: testutils.OrdersDeterministic("USD/GBP", order.OrderSell, 2, 11, 2, 11, 100),
		},
		{
			name:        "one_thousand_range_tenth",
			orders:      testutils.OrdersDeterministic("USD/GBP", order.OrderBuy, 1, 100, 1, 100, 1_000),
			matchOrders: testutils.OrdersDeterministic("USD/GBP", order.OrderSell, 2, 110, 2, 110, 100),
		},
		{
			name:        "ten_thousand_range_hundrendth",
			orders:      testutils.OrdersDeterministic("USD/GBP", order.OrderBuy, 1, 100, 1, 100, 10_000),
			matchOrders: testutils.OrdersDeterministic("USD/GBP", order.OrderSell, 2, 110, 2, 110, 1_000),
		},
		{
			name:        "ten_thousand_range_tenth",
			orders:      testutils.OrdersDeterministic("USD/GBP", order.OrderBuy, 1, 1_000, 1, 1_000, 10_000),
			matchOrders: testutils.OrdersDeterministic("USD/GBP", order.OrderSell, 2, 1_100, 2, 1_100, 1_000),
		},
		{
			name:        "hundred_thousand_range_hundrendth",
			orders:      testutils.OrdersDeterministic("USD/GBP", order.OrderBuy, 1, 1_000, 1, 1_000, 100_000),
			matchOrders: testutils.OrdersDeterministic("USD/GBP", order.OrderSell, 2, 1_100, 2, 1_100, 10_000),
		},
		{
			name:        "hundred_thousand_range_tenth",
			orders:      testutils.OrdersDeterministic("USD/GBP", order.OrderBuy, 1, 10_000, 1, 10_000, 100_000),
			matchOrders: testutils.OrdersDeterministic("USD/GBP", order.OrderSell, 2, 11_000, 2, 11_000, 10_000),
		},
		{
			name:        "million_range_hundrendth",
			orders:      testutils.OrdersDeterministic("USD/GBP", order.OrderBuy, 1, 10_000, 1, 10_000, 1_000_000),
			matchOrders: testutils.OrdersDeterministic("USD/GBP", order.OrderSell, 2, 11_000, 2, 11_000, 100_000),
		},
		{
			name:        "million_range_tenth",
			orders:      testutils.OrdersDeterministic("USD/GBP", order.OrderBuy, 1, 100_000, 1, 100_000, 1_000_000),
			matchOrders: testutils.OrdersDeterministic("USD/GBP", order.OrderSell, 2, 110_000, 2, 110_000, 100_000),
		},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			b.ReportAllocs()

			for range b.N {
				b.StopTimer()

				m := market.New(pair, tracker.orderEventsChan, tracker.volumeEventsChan, tracker.matchEventsChan)

				for _, o := range tc.orders {
					oCopy := &order.Order{ // the order is modified after matched, we need a copy
						ID:     o.ID,
						Pair:   o.Pair,
						Side:   o.Side,
						Volume: o.Volume,
						Price:  o.Price,
					}
					if err := m.InsertMakerOrder(oCopy); err != nil {
						b.Fatalf("InsertMakerOrder(%v) unexpected error: %v", o, err)
					}
				}

				b.StartTimer()

				for _, o := range tc.matchOrders {
					oCopy := &order.Order{ // the order is modified after matched, we need a copy
						ID:     o.ID,
						Pair:   o.Pair,
						Side:   o.Side,
						Volume: o.Volume,
						Price:  o.Price,
					}
					if err := m.MatchTakerOrder(oCopy); err != nil {
						b.Fatalf("Cancel(%v) unexpected error: %v", o, err)
					}
				}
			}
		})
	}
}
