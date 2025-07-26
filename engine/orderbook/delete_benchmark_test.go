package orderbook_test

import (
	"exchange/engine/order"
	"exchange/engine/orderbook"
	"exchange/engine/orderbook/rbtree"
	"exchange/engine/testutils"
	"sync"
	"testing"
)

func Benchmark_Delete(b *testing.B) {
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
				pool := &sync.Pool{
					New: func() any {
						return rbtree.NewNode()
					},
				}
				book := orderbook.New(order.OrderBuy, pool, func(uint64, uint64) {})

				for _, o := range tc.orders {
					if err := book.Insert(o); err != nil {
						b.Fatalf("Insert() unexpected error: %v", err)
					}
				}

				b.StartTimer()

				for _, o := range tc.orders {
					if err := book.Delete(o); err != nil {
						b.Fatalf("Delete() unexpected error: %v", err)
					}
				}
			}
		})
	}
}
