package pricelevel_test

import (
	"fmt"
	"math/rand/v2"
	"testing"

	"exchange/engine/order"
	"exchange/engine/orderbook/pricelevel"
	"exchange/engine/testutils"
)

func Benchmark_Remove(b *testing.B) {
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

			for range b.N {
				b.StopTimer()

				orders := make([]*order.Order, 0, tc.size)
				orderIDs := make([]string, 0, tc.size)
				for _, num := range testutils.NumbersAscending(tc.size) {
					orderID := fmt.Sprintf("%d", num)
					orders = append(orders, &order.Order{ID: orderID, Price: price, Volume: num})
					orderIDs = append(orderIDs, orderID)
				}

				rand.Shuffle(len(orders), func(i, j int) {
					orderIDs[i], orderIDs[j] = orderIDs[j], orderIDs[i]
				})

				p := pricelevel.New()

				for _, o := range orders {
					p.Insert(o)
				}

				b.StartTimer()

				for _, id := range orderIDs {
					p.Remove(id)
				}
			}
		})
	}
}
