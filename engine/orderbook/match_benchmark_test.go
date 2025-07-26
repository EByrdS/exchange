package orderbook_test

import (
	"exchange/engine/order"
	"exchange/engine/orderbook"
	"exchange/engine/orderbook/rbtree"
	"exchange/engine/testutils"
	"sync"
	"testing"
)

func Benchmark_Match(b *testing.B) {
	testCases := []struct {
		name       string
		orderCount uint64
		maxPrice   uint64
		maxVolume  uint64
	}{
		{
			name:       "one_thousand_range_hundredth",
			orderCount: 1_000,
			maxPrice:   10,
			maxVolume:  10,
		},
		{
			name:       "one_thousand_range_tenth",
			orderCount: 1_000,
			maxPrice:   100,
			maxVolume:  100,
		},
		{
			name:       "ten_thousand_range_hundrendth",
			orderCount: 10_000,
			maxPrice:   100,
			maxVolume:  100,
		},
		{
			name:       "ten_thousand_range_tenth",
			orderCount: 10_000,
			maxPrice:   1_000,
			maxVolume:  1_000,
		},
		{
			name:       "hundred_thousand_range_hundrendth",
			orderCount: 100_000,
			maxPrice:   1_000,
			maxVolume:  1_000,
		},
		{
			name:       "hundred_thousand_range_tenth",
			orderCount: 1000_000,
			maxPrice:   10_000,
			maxVolume:  10_000,
		},
		{
			name:       "million_range_hundrendth",
			orderCount: 1_000_000,
			maxPrice:   10_000,
			maxVolume:  10_000,
		},
		{
			name:       "million_range_tenth",
			orderCount: 1_000_000,
			maxPrice:   100_000,
			maxVolume:  100_000,
		},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			b.ReportAllocs()

			for range b.N {
				b.StopTimer()
				pool := &sync.Pool{
					New: func() any {
						return rbtree.NewNode()
					},
				}
				book := orderbook.New(order.OrderBuy, pool, func(uint64, uint64) {})

				orders := testutils.OrdersDeterministic("USD/GBP", order.OrderBuy, 1, tc.maxPrice, 1, tc.maxVolume, tc.orderCount)
				totalVolume := (tc.maxVolume * (tc.maxVolume + 1)) * (tc.orderCount / tc.maxVolume)
				extractVolume := uint64(float64(totalVolume) * 0.85)

				for _, o := range orders {
					if err := book.Insert(o); err != nil {
						b.Fatalf("Insert() unexpected error: %v", err)
					}
				}

				b.StartTimer()

				book.MatchAndExtract(extractVolume)
			}
		})
	}
}
