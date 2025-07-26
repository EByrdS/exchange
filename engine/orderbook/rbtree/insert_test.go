package rbtree_test

import (
	"sync"
	"testing"

	"exchange/engine/orderbook/rbtree"
	"exchange/engine/testutils"
)

func Test_Insert(t *testing.T) {
	testCases := []struct {
		name     string
		tree     *rbtree.Tree
		price    uint64
		wantRed  bool
		wantTree *rbtree.Tree
	}{
		{
			name:     "no_root",
			tree:     tree(nil),
			price:    10,
			wantTree: tree(&rbtree.Node{Price: 10}),
		},
		{
			name:    "left",
			tree:    tree(&rbtree.Node{Price: 5}),
			price:   2,
			wantRed: true,
			wantTree: tree(&rbtree.Node{Price: 5,
				Left: &rbtree.Node{Price: 2, Red: true}}),
		},
		{
			name:    "right",
			tree:    tree(&rbtree.Node{Price: 5}),
			price:   7,
			wantRed: true,
			wantTree: tree(&rbtree.Node{Price: 5,
				Right: &rbtree.Node{Price: 7, Red: true}}),
		},
		{
			name: "left_line",
			tree: tree(&rbtree.Node{Price: 5,
				Left: &rbtree.Node{Price: 3, Red: true}}),
			price:   1,
			wantRed: true,
			wantTree: tree(&rbtree.Node{Price: 3,
				Left:  &rbtree.Node{Price: 1, Red: true},
				Right: &rbtree.Node{Price: 5, Red: true}}),
		},
		{
			name: "left_arrow",
			tree: tree(&rbtree.Node{Price: 5,
				Left: &rbtree.Node{Price: 3, Red: true}}),
			price: 4,
			wantTree: tree(&rbtree.Node{Price: 4,
				Left:  &rbtree.Node{Price: 3, Red: true},
				Right: &rbtree.Node{Price: 5, Red: true}}),
		},
		{
			name: "right_arrow",
			tree: tree(&rbtree.Node{Price: 5,
				Right: &rbtree.Node{Price: 7, Red: true}}),
			price: 6,
			wantTree: tree(&rbtree.Node{Price: 6,
				Left:  &rbtree.Node{Price: 5, Red: true},
				Right: &rbtree.Node{Price: 7, Red: true}}),
		},
		{
			name: "right_line",
			tree: tree(&rbtree.Node{Price: 5,
				Right: &rbtree.Node{Price: 7, Red: true}}),
			price:   8,
			wantRed: true,
			wantTree: tree(&rbtree.Node{Price: 7,
				Left:  &rbtree.Node{Price: 5, Red: true},
				Right: &rbtree.Node{Price: 8, Red: true}}),
		},
		{
			name: "left_left",
			tree: tree(&rbtree.Node{Price: 5,
				Left:  &rbtree.Node{Price: 3},
				Right: &rbtree.Node{Price: 7}}),
			price:   1,
			wantRed: true,
			wantTree: tree(&rbtree.Node{Price: 5,
				Left: &rbtree.Node{Price: 3,
					Left: &rbtree.Node{Price: 1, Red: true}},
				Right: &rbtree.Node{Price: 7}}),
		},
		{
			name: "left_right",
			tree: tree(&rbtree.Node{Price: 5,
				Left:  &rbtree.Node{Price: 3},
				Right: &rbtree.Node{Price: 7}}),
			price:   4,
			wantRed: true,
			wantTree: tree(&rbtree.Node{Price: 5,
				Left: &rbtree.Node{Price: 3,
					Right: &rbtree.Node{Price: 4, Red: true}},
				Right: &rbtree.Node{Price: 7}}),
		},
		{
			name: "right_left",
			tree: tree(&rbtree.Node{Price: 5,
				Left:  &rbtree.Node{Price: 3},
				Right: &rbtree.Node{Price: 7}}),
			price:   6,
			wantRed: true,
			wantTree: tree(&rbtree.Node{Price: 5,
				Left: &rbtree.Node{Price: 3},
				Right: &rbtree.Node{Price: 7,
					Left: &rbtree.Node{Price: 6, Red: true}}}),
		},
		{
			name: "right_right",
			tree: tree(&rbtree.Node{Price: 5,
				Left:  &rbtree.Node{Price: 3},
				Right: &rbtree.Node{Price: 7}}),
			price:   8,
			wantRed: true,
			wantTree: tree(&rbtree.Node{Price: 5,
				Left: &rbtree.Node{Price: 3},
				Right: &rbtree.Node{Price: 7,
					Right: &rbtree.Node{Price: 8, Red: true}}}),
		},
		{
			name: "left_line_extra_node",
			tree: tree(&rbtree.Node{Price: 5,
				Left:  &rbtree.Node{Price: 3, Red: true},
				Right: &rbtree.Node{Price: 7, Red: true}}),
			price:   1,
			wantRed: true,
			wantTree: tree(&rbtree.Node{Price: 5,
				Left: &rbtree.Node{Price: 3,
					Left: &rbtree.Node{Price: 1, Red: true}},
				Right: &rbtree.Node{Price: 7}}),
		},
		{
			name: "left_arrow_extra_node",
			tree: tree(&rbtree.Node{Price: 5,
				Left:  &rbtree.Node{Price: 3, Red: true},
				Right: &rbtree.Node{Price: 7, Red: true}}),
			price:   4,
			wantRed: true,
			wantTree: tree(&rbtree.Node{Price: 5,
				Left: &rbtree.Node{Price: 3,
					Right: &rbtree.Node{Price: 4, Red: true}},
				Right: &rbtree.Node{Price: 7}}),
		},
		{
			name: "right_line_extra_node",
			tree: tree(&rbtree.Node{Price: 5,
				Left:  &rbtree.Node{Price: 3, Red: true},
				Right: &rbtree.Node{Price: 7, Red: true}}),
			price:   8,
			wantRed: true,
			wantTree: tree(&rbtree.Node{Price: 5,
				Left: &rbtree.Node{Price: 3},
				Right: &rbtree.Node{Price: 7,
					Right: &rbtree.Node{Price: 8, Red: true}}}),
		},
		{
			name: "right_arrow_extra_node",
			tree: tree(&rbtree.Node{Price: 5,
				Left:  &rbtree.Node{Price: 3, Red: true},
				Right: &rbtree.Node{Price: 7, Red: true}}),
			price:   6,
			wantRed: true,
			wantTree: tree(&rbtree.Node{Price: 5,
				Left: &rbtree.Node{Price: 3},
				Right: &rbtree.Node{Price: 7,
					Left: &rbtree.Node{Price: 6, Red: true}}}),
		},
		{
			name: "nested_right_line",
			tree: tree(&rbtree.Node{Price: 4,
				Right: &rbtree.Node{Price: 6, Red: true,
					Right: &rbtree.Node{Price: 8,
						Right: &rbtree.Node{Price: 9, Red: true},
						Left:  &rbtree.Node{Price: 7, Red: true}},
					Left: &rbtree.Node{Price: 5}},
				Left: &rbtree.Node{Price: 2, Red: true,
					Right: &rbtree.Node{Price: 3},
					Left:  &rbtree.Node{Price: 1}}}),
			price:   10,
			wantRed: true,
			wantTree: tree(&rbtree.Node{Price: 4,
				Right: &rbtree.Node{Price: 6,
					Right: &rbtree.Node{Price: 8, Red: true,
						Right: &rbtree.Node{Price: 9,
							Right: &rbtree.Node{Price: 10, Red: true}},
						Left: &rbtree.Node{Price: 7}},
					Left: &rbtree.Node{Price: 5}},
				Left: &rbtree.Node{Price: 2,
					Right: &rbtree.Node{Price: 3},
					Left:  &rbtree.Node{Price: 1}}}),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			connectParents(tc.tree)
			if err := tc.tree.Valid(); err != nil {
				t.Fatalf("Invalid starting tree: %v", err)
			}
			connectParents(tc.wantTree)
			if err := tc.wantTree.Valid(); err != nil {
				t.Fatalf("Invalid want tree: %v", err)
			}

			gotNode := tc.tree.Insert(tc.price)

			if err := tc.tree.Valid(); err != nil {
				t.Fatalf("Invalid got tree: %v\ngot:\n%s", err, tc.tree.PrintTree())
			}

			if gotNode.Price != tc.price {
				t.Errorf("Insert(%v) expected node with price: %v, got price: %v", tc.price, tc.price, gotNode.Price)
			}

			if gotNode.Red != tc.wantRed {
				t.Errorf("Insert(%v) expected red: %v, got red: %v", tc.price, tc.wantRed, gotNode.Red)
			}

			gotTree := tc.tree.PrintTree()
			wantTree := tc.wantTree.PrintTree()
			if gotTree != wantTree {
				t.Errorf("Insert(%v)\nwant tree:\n%s\ngot tree:\n%s\n", tc.price, wantTree, gotTree)
			}
		})
	}
}

func Test_Insert_Head(t *testing.T) {
	testCases := []struct {
		name          string
		orientation   rbtree.TreeOrientation
		insertions    []uint64
		wantHead      bool
		wantHeadPrice uint64
	}{
		{
			name:        "min_empty",
			orientation: rbtree.MinFirst,
		},
		{
			name:        "max_empty",
			orientation: rbtree.MaxFirst,
		},
		{
			name:          "min_single",
			orientation:   rbtree.MinFirst,
			insertions:    []uint64{1},
			wantHead:      true,
			wantHeadPrice: 1,
		},
		{
			name:          "max_single",
			orientation:   rbtree.MaxFirst,
			insertions:    []uint64{1},
			wantHead:      true,
			wantHeadPrice: 1,
		},
		{
			name:          "min_keep_previous",
			orientation:   rbtree.MinFirst,
			insertions:    []uint64{1, 2},
			wantHead:      true,
			wantHeadPrice: 1,
		},
		{
			name:          "max_keep_previous",
			orientation:   rbtree.MaxFirst,
			insertions:    []uint64{2, 1},
			wantHead:      true,
			wantHeadPrice: 2,
		},
		{
			name:          "min_override",
			orientation:   rbtree.MinFirst,
			insertions:    []uint64{2, 1},
			wantHead:      true,
			wantHeadPrice: 1,
		},
		{
			name:          "max_override",
			orientation:   rbtree.MaxFirst,
			insertions:    []uint64{1, 2},
			wantHead:      true,
			wantHeadPrice: 2,
		},
		{
			name:          "min_multiple",
			orientation:   rbtree.MinFirst,
			insertions:    []uint64{7, 5, 6, 4, 1, 3, 2, 8},
			wantHead:      true,
			wantHeadPrice: 1,
		},
		{
			name:          "max_multiple",
			orientation:   rbtree.MaxFirst,
			insertions:    []uint64{1, 3, 2, 4, 8, 6, 7, 5},
			wantHead:      true,
			wantHeadPrice: 8,
		},
		{
			name:          "min_random",
			orientation:   rbtree.MinFirst,
			insertions:    testutils.NumbersRandom(100),
			wantHead:      true,
			wantHeadPrice: 1,
		},
		{
			name:          "max_random",
			orientation:   rbtree.MaxFirst,
			insertions:    testutils.NumbersRandom(100),
			wantHead:      true,
			wantHeadPrice: 100,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			pool := &sync.Pool{
				New: func() any {
					return rbtree.NewNode()
				},
			}
			tree := rbtree.NewTree(tc.orientation, pool)
			for _, num := range tc.insertions {
				tree.Insert(num)
			}

			head := tree.Head()
			if !tc.wantHead {
				if head != nil {
					t.Errorf("Head() wanted nil head, got: %v", head)
				}
			} else {
				if head == nil {
					t.Fatalf("Head() wanted head with price %d, got nil head", tc.wantHeadPrice)
				}

				if tc.wantHeadPrice != head.Price {
					t.Errorf("Head() wanted price %d, got price %d", tc.wantHeadPrice, head.Price)
				}
			}
		})
	}
}

func Test_Insert_Sequence(t *testing.T) {
	testCases := []struct {
		name string
		seed func(uint64) []uint64
	}{
		{
			name: "ascending",
			seed: testutils.NumbersAscending,
		},
		{
			name: "descending",
			seed: testutils.NumbersDescending,
		},
		{
			name: "random",
			seed: testutils.NumbersRandom,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			seed := tc.seed(100)
			tr := tree(nil)
			for _, num := range seed {
				start := tr.PrintTree()

				tr.Insert(num)

				if err := tr.Valid(); err != nil {
					t.Fatalf("Insert(%v) invalid tree: %v\nStart:\n%s\nGot:\n%s", num, err, start, tr.PrintTree())
				}
			}
		})
	}
}
