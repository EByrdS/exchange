package market_test

import (
	"exchange/engine/market"
	"exchange/engine/order"
	"exchange/engine/testutils"
	"testing"
)

func Benchmark_Cancel(b *testing.B) {
	tracker := newEventsTracker(10)
	tracker.ignoreAll()

	pair := "USD/GBP"

	testCases := []struct {
		name   string
		orders []*order.Order
	}{
		{
			name:   "one_thousand_range_hundredth",
			orders: testutils.OrdersDeterministic("USD/GBP", order.OrderBuy, 1, 10, 1, 10, 1_000),
		},
		{
			name:   "one_thousand_range_tenth",
			orders: testutils.OrdersDeterministic("USD/GBP", order.OrderBuy, 1, 100, 1, 100, 1_000),
		},
		{
			name:   "ten_thousand_range_hundrendth",
			orders: testutils.OrdersDeterministic("USD/GBP", order.OrderBuy, 1, 100, 1, 100, 10_000),
		},
		{
			name:   "ten_thousand_range_tenth",
			orders: testutils.OrdersDeterministic("USD/GBP", order.OrderBuy, 1, 1_000, 1, 1_000, 10_000),
		},
		{
			name:   "hundred_thousand_range_hundrendth",
			orders: testutils.OrdersDeterministic("USD/GBP", order.OrderBuy, 1, 1_000, 1, 1_000, 100_000),
		},
		{
			name:   "hundred_thousand_range_tenth",
			orders: testutils.OrdersDeterministic("USD/GBP", order.OrderBuy, 1, 10_000, 1, 10_000, 100_000),
		},
		{
			name:   "million_range_hundrendth",
			orders: testutils.OrdersDeterministic("USD/GBP", order.OrderBuy, 1, 10_000, 1, 10_000, 1_000_000),
		},
		{
			name:   "million_range_tenth",
			orders: testutils.OrdersDeterministic("USD/GBP", order.OrderBuy, 1, 100_000, 1, 100_000, 1_000_000),
		},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			b.ReportAllocs()

			for range b.N {
				b.StopTimer()

				m := market.New(pair, tracker.orderEventsChan, tracker.volumeEventsChan, tracker.matchEventsChan)

				for _, o := range tc.orders {
					if err := m.InsertMakerOrder(o); err != nil {
						b.Fatalf("InsertMakerOrder(%v) unexpected error: %v", o, err)
					}
				}

				b.StartTimer()

				for _, o := range tc.orders {
					if err := m.Cancel(o); err != nil {
						b.Fatalf("Cancel(%v) unexpected error: %v", o, err)
					}
				}
			}
		})
	}
}
