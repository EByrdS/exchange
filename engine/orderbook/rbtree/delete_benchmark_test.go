package rbtree_test

import (
	"testing"

	"exchange/engine/testutils"
)

func Benchmark_Delete(b *testing.B) {
	testCases := []struct {
		name   string
		insert func(uint64) []uint64
		delete func(uint64) []uint64
	}{
		{
			name:   "insert_asc_delete_asc_10k",
			insert: testutils.NumbersAscending,
			delete: testutils.NumbersAscending,
		},
		{
			name:   "insert_asc_delete_desc_10k",
			insert: testutils.NumbersAscending,
			delete: testutils.NumbersDescending,
		},
		{
			name:   "insert_asc_delete_random_10k",
			insert: testutils.NumbersAscending,
			delete: testutils.NumbersRandom,
		},
		{
			name:   "insert_desc_delete_asc_10k",
			insert: testutils.NumbersDescending,
			delete: testutils.NumbersAscending,
		},
		{
			name:   "insert_desc_delete_desc_10k",
			insert: testutils.NumbersDescending,
			delete: testutils.NumbersDescending,
		},
		{
			name:   "insert_desc_delete_random_10k",
			insert: testutils.NumbersDescending,
			delete: testutils.NumbersRandom,
		},
		{
			name:   "insert_random_delete_asc_10k",
			insert: testutils.NumbersRandom,
			delete: testutils.NumbersAscending,
		},
		{
			name:   "insert_random_delete_desc_10k",
			insert: testutils.NumbersRandom,
			delete: testutils.NumbersDescending,
		},
		{
			name:   "insert_random_delete_random_10k",
			insert: testutils.NumbersRandom,
			delete: testutils.NumbersRandom,
		},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			b.ReportAllocs()

			insert := tc.insert(10000)
			delete := tc.delete(10000)

			b.ResetTimer()
			for range b.N {
				b.StopTimer()
				tr := tree(nil)
				for _, num := range insert {
					tr.Insert(num)
				}

				b.StartTimer()
				for _, num := range delete {
					tr.Delete(num)
				}
			}
		})
	}
}
