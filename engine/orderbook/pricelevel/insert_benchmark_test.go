package pricelevel_test

import (
	"fmt"
	"testing"

	"exchange/engine/order"
	"exchange/engine/orderbook/pricelevel"
	"exchange/engine/testutils"
)

func Benchmark_Insert(b *testing.B) {
	const price uint64 = 1

	testCases := []struct {
		name string
		size uint64
	}{
		{
			name: "one_thousand",
			size: 1_000,
		},
		{
			name: "ten_thousand",
			size: 10_000,
		},
		{
			name: "hundred_thousand",
			size: 100_000,
		},
		{
			name: "million",
			size: 1_000_000,
		},
	}
	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			b.ReportAllocs()

			orders := make([]*order.Order, 0, tc.size)
			for _, num := range testutils.NumbersAscending(tc.size) {
				orders = append(orders, &order.Order{ID: fmt.Sprintf("%d", num), Price: price, Volume: num})
			}

			b.ResetTimer()
			for range b.N {
				b.StopTimer()
				p := pricelevel.New()
				b.StartTimer()

				for _, o := range orders {
					p.Insert(o)
				}
			}
		})
	}
}
