package rbtree_test

import (
	"runtime"
	"testing"

	"exchange/engine/testutils"
)

func Benchmark_Insert_Delete(b *testing.B) {
	testCases := []struct {
		name   string
		size   uint64
		cycles int
		insert func(uint64) []uint64
		delete func(uint64) []uint64
	}{
		{
			name:   "insert_random_delete_random_10_cycles_100",
			size:   10,
			cycles: 100,
			insert: testutils.NumbersRandom,
			delete: testutils.NumbersRandom,
		},
		{
			name:   "insert_random_delete_random_100_cycles_100",
			size:   100,
			cycles: 100,
			insert: testutils.NumbersRandom,
			delete: testutils.NumbersRandom,
		},
		{
			name:   "insert_random_delete_random_1000_cycles_100",
			size:   1000,
			cycles: 100,
			insert: testutils.NumbersRandom,
			delete: testutils.NumbersRandom,
		},
		{
			name:   "insert_random_delete_random_10000_cycles_100",
			size:   10000,
			cycles: 100,
			insert: testutils.NumbersRandom,
			delete: testutils.NumbersRandom,
		},
		{
			name:   "insert_random_delete_random_10_cycles_1000",
			size:   10,
			cycles: 1000,
			insert: testutils.NumbersRandom,
			delete: testutils.NumbersRandom,
		},
		{
			name:   "insert_random_delete_random_100_cycles_1000",
			size:   100,
			cycles: 1000,
			insert: testutils.NumbersRandom,
			delete: testutils.NumbersRandom,
		},
		{
			name:   "insert_random_delete_random_1000_cycles_1000",
			size:   1000,
			cycles: 1000,
			insert: testutils.NumbersRandom,
			delete: testutils.NumbersRandom,
		},
		{
			name:   "insert_random_delete_random_10000_cycles_1000",
			size:   10000,
			cycles: 1000,
			insert: testutils.NumbersRandom,
			delete: testutils.NumbersRandom,
		},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			b.ReportAllocs()

			b.ResetTimer()
			for range b.N {
				b.StopTimer()
				tr := tree(nil)
				insert := tc.insert(tc.size)
				delete := tc.delete(tc.size)
				b.StartTimer()

				for j := 0; j < tc.cycles; j++ {
					for _, num := range insert {
						tr.Insert(num)
					}

					for _, num := range delete {
						tr.Delete(num)
					}

					if j%100 == 0 {
						b.StopTimer()
						runtime.GC() // GC is very expensive, avoid counting it
						b.StartTimer()
					}
				}
			}
		})
	}
}
