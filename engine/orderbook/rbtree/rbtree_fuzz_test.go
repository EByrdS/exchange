package rbtree_test

import (
	"testing"
)

func Fuzz_Insert_Delete(f *testing.F) {
	f.Add([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15})
	f.Add([]byte{15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1})
	f.Add([]byte{15, 1, 14, 2, 13, 3, 12, 4, 11, 5, 10, 6, 9, 7, 8})
	f.Add([]byte{8, 9, 7, 10, 6, 11, 5, 12, 4, 13, 3, 14, 2, 15, 1})

	f.Fuzz(func(t *testing.T, data []byte) {
		tr := tree(nil)
		prices := make([]uint64, 0, len(data))
		for i := range len(data) {
			prices = append(prices, uint64(data[i]))
		}

		for i, v := range prices {
			tr.Insert(v)

			if err := tr.Valid(); err != nil {
				t.Fatalf("Insert(%v). Invalid tree after sequence: %v. Error: %v", v, prices[:i], err)
			}
		}

		for i, v := range prices {
			tr.Delete(v)

			if err := tr.Valid(); err != nil {
				t.Fatalf("Delete(%v). Invalid tree after insert sequence: %v, and delete sequence: %v. Error: %v", v, prices, prices[:i], err)
			}
		}
	})
}
