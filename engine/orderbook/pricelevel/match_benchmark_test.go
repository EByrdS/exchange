package pricelevel_test

import (
	"fmt"
	"testing"

	"exchange/engine/order"
	"exchange/engine/orderbook/pricelevel"
	"exchange/engine/testutils"
)

func Benchmark_MatchAndExtract(b *testing.B) {
	// Gauss sumation
	sumation := func(n uint64) uint64 { return uint64(n * (n + 1) / 2) }

	testCases := []struct {
		name    string
		size    uint64
		extract uint64
	}{
		{
			name:    "one_thousand",
			size:    1_000,
			extract: sumation(1_000) - 100,
		},
		{
			name:    "ten_thousand",
			size:    10_000,
			extract: sumation(10_000) - 1_000,
		},
		{
			name:    "hundred_thousand",
			size:    100_000,
			extract: sumation(100_000) - 10_000,
		},
		{
			name:    "million",
			size:    1_000_000,
			extract: sumation(1_000_000) - 100_000,
		},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			b.ReportAllocs()

			orders := make([]order.Order, 0, tc.size)
			for _, num := range testutils.NumbersAscending(tc.size) {
				orders = append(orders, order.Order{Price: 1, ID: fmt.Sprintf("%d", num), Volume: num})
			}

			for range b.N {
				b.StopTimer()

				ordersCopy := make([]order.Order, len(orders))
				copy(ordersCopy, orders)

				p := pricelevel.New()

				for _, o := range ordersCopy {
					if err := p.Insert(&o); err != nil {
						b.Fatalf("Insert(%+v) unexpected error: %v", o, err)
					}
				}

				b.StartTimer()

				p.MatchAndExtract(tc.extract)
			}
		})
	}
}
