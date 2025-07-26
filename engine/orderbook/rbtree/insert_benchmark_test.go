package rbtree_test

import (
	"testing"

	"exchange/engine/testutils"
)

func Benchmark_Insert(b *testing.B) {
	testCases := []struct {
		name string
		seed func(uint64) []uint64
	}{
		{
			name: "ascending_10k",
			seed: testutils.NumbersAscending,
		},
		{
			name: "descending_10k",
			seed: testutils.NumbersDescending,
		},
		{
			name: "random_10k",
			seed: testutils.NumbersRandom,
		},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			b.ReportAllocs()

			seed := tc.seed(10000)

			b.ResetTimer()
			for range b.N {
				tr := tree(nil)
				for _, num := range seed {
					tr.Insert(num)
				}
			}
		})
	}
}
