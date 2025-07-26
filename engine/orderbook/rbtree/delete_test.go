package rbtree_test

import (
	"sync"
	"testing"

	"exchange/engine/orderbook/rbtree"
	"exchange/engine/testutils"
)

func Test_Delete(t *testing.T) {
	testCases := []struct {
		name     string
		tree     *rbtree.Tree
		price    uint64
		wantTree *rbtree.Tree
	}{
		{
			name:     "no_root",
			tree:     tree(nil),
			price:    10,
			wantTree: tree(nil),
		},
		{
			name: "not_found",
			tree: tree(&rbtree.Node{Price: 5,
				Right: &rbtree.Node{Price: 7},
				Left:  &rbtree.Node{Price: 3}}),
			price: 10,
			wantTree: tree(&rbtree.Node{Price: 5,
				Right: &rbtree.Node{Price: 7},
				Left:  &rbtree.Node{Price: 3}}),
		},
		{
			name:     "root",
			tree:     tree(&rbtree.Node{Price: 5}),
			price:    5,
			wantTree: tree(nil),
		},
		{
			name: "root_with_left_child",
			tree: tree(&rbtree.Node{Price: 5,
				Left: &rbtree.Node{Price: 3, Red: true}}),
			price:    5,
			wantTree: tree(&rbtree.Node{Price: 3}),
		},
		{
			name: "root_with_right_child",
			tree: tree(&rbtree.Node{Price: 5,
				Right: &rbtree.Node{Price: 7, Red: true}}),
			price:    5,
			wantTree: tree(&rbtree.Node{Price: 7}),
		},
		{
			name: "root_with_children",
			tree: tree(&rbtree.Node{Price: 5,
				Right: &rbtree.Node{Price: 6},
				Left:  &rbtree.Node{Price: 4}}),
			price: 5,
			wantTree: tree(&rbtree.Node{Price: 6,
				Left: &rbtree.Node{Price: 4, Red: true}}),
		},
		{
			name: "root_with_grandchildren",
			tree: tree(&rbtree.Node{Price: 5,
				Right: &rbtree.Node{Price: 7,
					Right: &rbtree.Node{Price: 8},
					Left:  &rbtree.Node{Price: 6}},
				Left: &rbtree.Node{Price: 3,
					Right: &rbtree.Node{Price: 4},
					Left:  &rbtree.Node{Price: 2}}}),
			price: 5,
			wantTree: tree(&rbtree.Node{Price: 6,
				Right: &rbtree.Node{Price: 7,
					Right: &rbtree.Node{Price: 8, Red: true}},
				Left: &rbtree.Node{Price: 3, Red: true,
					Right: &rbtree.Node{Price: 4},
					Left:  &rbtree.Node{Price: 2}}}),
		},
		{
			name: "root_whose_successor_has_right_child",
			tree: tree(&rbtree.Node{Price: 5,
				Right: &rbtree.Node{Price: 8,
					Right: &rbtree.Node{Price: 9},
					Left: &rbtree.Node{Price: 6,
						Right: &rbtree.Node{Price: 7, Red: true}}},
				Left: &rbtree.Node{Price: 3,
					Right: &rbtree.Node{Price: 4},
					Left:  &rbtree.Node{Price: 2}}}),
			price: 5,
			wantTree: tree(&rbtree.Node{Price: 6,
				Right: &rbtree.Node{Price: 8,
					Right: &rbtree.Node{Price: 9},
					Left:  &rbtree.Node{Price: 7}},
				Left: &rbtree.Node{Price: 3,
					Right: &rbtree.Node{Price: 4},
					Left:  &rbtree.Node{Price: 2}}}),
		},
		{
			name: "left_leaf",
			tree: tree(&rbtree.Node{Price: 5,
				Left: &rbtree.Node{Price: 3, Red: true}}),
			price:    3,
			wantTree: tree(&rbtree.Node{Price: 5}),
		},
		{
			name: "right_leaf",
			tree: tree(&rbtree.Node{Price: 5,
				Right: &rbtree.Node{Price: 7, Red: true}}),
			price:    7,
			wantTree: tree(&rbtree.Node{Price: 5}),
		},
		{
			name: "left_leaf_red_sibling",
			tree: tree(&rbtree.Node{Price: 7,
				Right: &rbtree.Node{Price: 8, Red: true},
				Left:  &rbtree.Node{Price: 5, Red: true}}),
			price: 5,
			wantTree: tree(&rbtree.Node{Price: 7,
				Right: &rbtree.Node{Price: 8, Red: true}}),
		},
		{
			name: "right_lead_red_sibling",
			tree: tree(&rbtree.Node{Price: 7,
				Right: &rbtree.Node{Price: 8, Red: true},
				Left:  &rbtree.Node{Price: 5, Red: true}}),
			price: 8,
			wantTree: tree(&rbtree.Node{Price: 7,
				Left: &rbtree.Node{Price: 5, Red: true}}),
		},
		{
			name: "left_left_leaf",
			tree: tree(&rbtree.Node{Price: 5,
				Right: &rbtree.Node{Price: 7},
				Left: &rbtree.Node{Price: 3,
					Left: &rbtree.Node{Price: 1, Red: true}}}),
			price: 1,
			wantTree: tree(&rbtree.Node{Price: 5,
				Right: &rbtree.Node{Price: 7},
				Left:  &rbtree.Node{Price: 3}}),
		},
		{
			name: "left_right_leaf",
			tree: tree(&rbtree.Node{Price: 5,
				Right: &rbtree.Node{Price: 7},
				Left: &rbtree.Node{Price: 3,
					Right: &rbtree.Node{Price: 4, Red: true}}}),
			price: 4,
			wantTree: tree(&rbtree.Node{Price: 5,
				Right: &rbtree.Node{Price: 7},
				Left:  &rbtree.Node{Price: 3}}),
		},
		{
			name: "right_left_leaf",
			tree: tree(&rbtree.Node{Price: 5,
				Right: &rbtree.Node{Price: 7,
					Left: &rbtree.Node{Price: 6, Red: true}},
				Left: &rbtree.Node{Price: 3}}),
			price: 6,
			wantTree: tree(&rbtree.Node{Price: 5,
				Right: &rbtree.Node{Price: 7},
				Left:  &rbtree.Node{Price: 3}}),
		},
		{
			name: "right_right_leaf",
			tree: tree(&rbtree.Node{Price: 5,
				Right: &rbtree.Node{Price: 7,
					Right: &rbtree.Node{Price: 8, Red: true}},
				Left: &rbtree.Node{Price: 3}}),
			price: 8,
			wantTree: tree(&rbtree.Node{Price: 5,
				Right: &rbtree.Node{Price: 7},
				Left:  &rbtree.Node{Price: 3}}),
		},
		{
			name: "left_line_extra_node",
			tree: tree(&rbtree.Node{Price: 5,
				Right: &rbtree.Node{Price: 7},
				Left: &rbtree.Node{Price: 3,
					Left: &rbtree.Node{Price: 1, Red: true}}}),
			price: 1,
			wantTree: tree(&rbtree.Node{Price: 5,
				Right: &rbtree.Node{Price: 7},
				Left:  &rbtree.Node{Price: 3}}),
		},
		{
			name: "left_arrow_extra_node",
			tree: tree(&rbtree.Node{Price: 5,
				Right: &rbtree.Node{Price: 7},
				Left: &rbtree.Node{Price: 3,
					Right: &rbtree.Node{Price: 4, Red: true}}}),
			price: 4,
			wantTree: tree(&rbtree.Node{Price: 5,
				Right: &rbtree.Node{Price: 7},
				Left:  &rbtree.Node{Price: 3}}),
		},
		{
			name: "right_line_extra_node",
			tree: tree(&rbtree.Node{Price: 5,
				Right: &rbtree.Node{Price: 7,
					Right: &rbtree.Node{Price: 8, Red: true}},
				Left: &rbtree.Node{Price: 3}}),
			price: 8,
			wantTree: tree(&rbtree.Node{Price: 5,
				Right: &rbtree.Node{Price: 7},
				Left:  &rbtree.Node{Price: 3}}),
		},
		{
			name: "right_arrow_extra_node",
			tree: tree(&rbtree.Node{Price: 5,
				Right: &rbtree.Node{Price: 7,
					Left: &rbtree.Node{Price: 6, Red: true}},
				Left: &rbtree.Node{Price: 3}}),
			price: 6,
			wantTree: tree(&rbtree.Node{Price: 5,
				Right: &rbtree.Node{Price: 7},
				Left:  &rbtree.Node{Price: 3}}),
		},
		{
			name: "black_no_children_right_far_red_nephew",
			tree: tree(&rbtree.Node{Price: 2,
				Right: &rbtree.Node{Price: 3,
					Right: &rbtree.Node{Price: 4, Red: true}},
				Left: &rbtree.Node{Price: 1}}),
			price: 1,
			wantTree: tree(&rbtree.Node{Price: 3,
				Right: &rbtree.Node{Price: 4},
				Left:  &rbtree.Node{Price: 2}}),
		},
		{
			name: "black_no_children_right_close_red_nephew",
			tree: tree(&rbtree.Node{Price: 2,
				Right: &rbtree.Node{Price: 4,
					Left: &rbtree.Node{Price: 3, Red: true}},
				Left: &rbtree.Node{Price: 1}}),
			price: 1,
			wantTree: tree(&rbtree.Node{Price: 3,
				Right: &rbtree.Node{Price: 4},
				Left:  &rbtree.Node{Price: 2}}),
		},
		{
			name: "black_no_children_left_far_red_nephew",
			tree: tree(&rbtree.Node{Price: 3,
				Right: &rbtree.Node{Price: 4},
				Left: &rbtree.Node{Price: 2,
					Left: &rbtree.Node{Price: 1, Red: true}}}),
			price: 4,
			wantTree: tree(&rbtree.Node{Price: 2,
				Right: &rbtree.Node{Price: 3},
				Left:  &rbtree.Node{Price: 1}}),
		},
		{
			name: "black_no_children_left_close_red_nephew",
			tree: tree(&rbtree.Node{Price: 3,
				Right: &rbtree.Node{Price: 4},
				Left: &rbtree.Node{Price: 1,
					Right: &rbtree.Node{Price: 2, Red: true}}}),
			price: 4,
			wantTree: tree(&rbtree.Node{Price: 2,
				Right: &rbtree.Node{Price: 3},
				Left:  &rbtree.Node{Price: 1}}),
		},
		{
			name: "black_no_children_black_nephews",
			tree: tree(&rbtree.Node{Price: 4,
				Right: &rbtree.Node{Price: 6,
					Right: &rbtree.Node{Price: 8, Red: true,
						Right: &rbtree.Node{Price: 9,
							Right: &rbtree.Node{Price: 10, Red: true}},
						Left: &rbtree.Node{Price: 7}},
					Left: &rbtree.Node{Price: 5}},
				Left: &rbtree.Node{Price: 2,
					Right: &rbtree.Node{Price: 3},
					Left:  &rbtree.Node{Price: 1}}}),
			price: 1,
			wantTree: tree(&rbtree.Node{Price: 6,
				Right: &rbtree.Node{Price: 8,
					Right: &rbtree.Node{Price: 9,
						Right: &rbtree.Node{Price: 10, Red: true}},
					Left: &rbtree.Node{Price: 7}},
				Left: &rbtree.Node{Price: 4,
					Right: &rbtree.Node{Price: 5},
					Left: &rbtree.Node{Price: 2,
						Right: &rbtree.Node{Price: 3, Red: true}}}}),
		},
		{
			name: "two_level_black_tree_delete_rightmost",
			tree: tree(&rbtree.Node{Price: 4,
				Right: &rbtree.Node{Price: 6,
					Right: &rbtree.Node{Price: 7},
					Left:  &rbtree.Node{Price: 5}},
				Left: &rbtree.Node{Price: 2,
					Right: &rbtree.Node{Price: 3},
					Left:  &rbtree.Node{Price: 1}}}),
			price: 7,
			wantTree: tree(&rbtree.Node{Price: 4,
				Right: &rbtree.Node{Price: 6,
					Left: &rbtree.Node{Price: 5, Red: true}},
				Left: &rbtree.Node{Price: 2, Red: true,
					Right: &rbtree.Node{Price: 3},
					Left:  &rbtree.Node{Price: 1}}}),
		},
		{
			name: "random_case_1",
			tree: tree(&rbtree.Node{Price: 8,
				Right: &rbtree.Node{Price: 10,
					Right: &rbtree.Node{Price: 12},
					Left:  &rbtree.Node{Price: 9}},
				Left: &rbtree.Node{Price: 4,
					Right: &rbtree.Node{Price: 6,
						Right: &rbtree.Node{Price: 7, Red: true}},
					Left: &rbtree.Node{Price: 3,
						Left: &rbtree.Node{Price: 1, Red: true}}}}),
			price: 8,
			wantTree: tree(&rbtree.Node{Price: 9,
				Right: &rbtree.Node{Price: 10,
					Right: &rbtree.Node{Price: 12, Red: true}},
				Left: &rbtree.Node{Price: 4, Red: true,
					Right: &rbtree.Node{Price: 6,
						Right: &rbtree.Node{Price: 7, Red: true}},
					Left: &rbtree.Node{Price: 3,
						Left: &rbtree.Node{Price: 1, Red: true}}}}),
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

			tc.tree.Delete(tc.price)

			if err := tc.tree.Valid(); err != nil {
				t.Fatalf("Invalid got tree: %v\ngot:\n%s", err, tc.tree.PrintTree())
			}

			gotTree := tc.tree.PrintTree()
			wantTree := tc.wantTree.PrintTree()
			if gotTree != wantTree {
				t.Errorf("Delete(%v)\nwant tree:\n%s\ngot tree:\n%s\n", tc.price, wantTree, gotTree)
			}
		})
	}
}

func Test_DeleteNode(t *testing.T) {
	testCases := []struct {
		name     string
		tree     *rbtree.Tree
		getNode  func(t *rbtree.Tree) *rbtree.Node
		wantTree *rbtree.Tree
	}{
		{
			name: "root",
			tree: tree(&rbtree.Node{Price: 5}),
			getNode: func(t *rbtree.Tree) *rbtree.Node {
				return t.Root
			},
			wantTree: tree(nil),
		},
		{
			name: "root_with_left_child",
			tree: tree(&rbtree.Node{Price: 5,
				Left: &rbtree.Node{Price: 3, Red: true}}),
			getNode: func(t *rbtree.Tree) *rbtree.Node {
				return t.Root
			},
			wantTree: tree(&rbtree.Node{Price: 3}),
		},
		{
			name: "root_with_right_child",
			tree: tree(&rbtree.Node{Price: 5,
				Right: &rbtree.Node{Price: 7, Red: true}}),
			getNode: func(t *rbtree.Tree) *rbtree.Node {
				return t.Root
			},
			wantTree: tree(&rbtree.Node{Price: 7}),
		},
		{
			name: "root_with_children",
			tree: tree(&rbtree.Node{Price: 5,
				Right: &rbtree.Node{Price: 6},
				Left:  &rbtree.Node{Price: 4}}),
			getNode: func(t *rbtree.Tree) *rbtree.Node {
				return t.Root
			},
			wantTree: tree(&rbtree.Node{Price: 6,
				Left: &rbtree.Node{Price: 4, Red: true}}),
		},
		{
			name: "root_with_grandchildren",
			tree: tree(&rbtree.Node{Price: 5,
				Right: &rbtree.Node{Price: 7,
					Right: &rbtree.Node{Price: 8},
					Left:  &rbtree.Node{Price: 6}},
				Left: &rbtree.Node{Price: 3,
					Right: &rbtree.Node{Price: 4},
					Left:  &rbtree.Node{Price: 2}}}),
			getNode: func(t *rbtree.Tree) *rbtree.Node {
				return t.Root
			},
			wantTree: tree(&rbtree.Node{Price: 6,
				Right: &rbtree.Node{Price: 7,
					Right: &rbtree.Node{Price: 8, Red: true}},
				Left: &rbtree.Node{Price: 3, Red: true,
					Right: &rbtree.Node{Price: 4},
					Left:  &rbtree.Node{Price: 2}}}),
		},
		{
			name: "root_whose_successor_has_right_child",
			tree: tree(&rbtree.Node{Price: 5,
				Right: &rbtree.Node{Price: 8,
					Right: &rbtree.Node{Price: 9},
					Left: &rbtree.Node{Price: 6,
						Right: &rbtree.Node{Price: 7, Red: true}}},
				Left: &rbtree.Node{Price: 3,
					Right: &rbtree.Node{Price: 4},
					Left:  &rbtree.Node{Price: 2}}}),
			getNode: func(t *rbtree.Tree) *rbtree.Node {
				return t.Root
			},
			wantTree: tree(&rbtree.Node{Price: 6,
				Right: &rbtree.Node{Price: 8,
					Right: &rbtree.Node{Price: 9},
					Left:  &rbtree.Node{Price: 7}},
				Left: &rbtree.Node{Price: 3,
					Right: &rbtree.Node{Price: 4},
					Left:  &rbtree.Node{Price: 2}}}),
		},
		{
			name: "left_leaf",
			tree: tree(&rbtree.Node{Price: 5,
				Left: &rbtree.Node{Price: 3, Red: true}}),
			getNode: func(t *rbtree.Tree) *rbtree.Node {
				return t.Root.Left
			},
			wantTree: tree(&rbtree.Node{Price: 5}),
		},
		{
			name: "right_leaf",
			tree: tree(&rbtree.Node{Price: 5,
				Right: &rbtree.Node{Price: 7, Red: true}}),
			getNode: func(t *rbtree.Tree) *rbtree.Node {
				return t.Root.Right
			},
			wantTree: tree(&rbtree.Node{Price: 5}),
		},
		{
			name: "left_leaf_red_sibling",
			tree: tree(&rbtree.Node{Price: 7,
				Right: &rbtree.Node{Price: 8, Red: true},
				Left:  &rbtree.Node{Price: 5, Red: true}}),
			getNode: func(t *rbtree.Tree) *rbtree.Node {
				return t.Root.Left
			},
			wantTree: tree(&rbtree.Node{Price: 7,
				Right: &rbtree.Node{Price: 8, Red: true}}),
		},
		{
			name: "right_lead_red_sibling",
			tree: tree(&rbtree.Node{Price: 7,
				Right: &rbtree.Node{Price: 8, Red: true},
				Left:  &rbtree.Node{Price: 5, Red: true}}),
			getNode: func(t *rbtree.Tree) *rbtree.Node {
				return t.Root.Right
			},
			wantTree: tree(&rbtree.Node{Price: 7,
				Left: &rbtree.Node{Price: 5, Red: true}}),
		},
		{
			name: "left_left_leaf",
			tree: tree(&rbtree.Node{Price: 5,
				Right: &rbtree.Node{Price: 7},
				Left: &rbtree.Node{Price: 3,
					Left: &rbtree.Node{Price: 1, Red: true}}}),
			getNode: func(t *rbtree.Tree) *rbtree.Node {
				return t.Root.Left.Left
			},
			wantTree: tree(&rbtree.Node{Price: 5,
				Right: &rbtree.Node{Price: 7},
				Left:  &rbtree.Node{Price: 3}}),
		},
		{
			name: "left_right_leaf",
			tree: tree(&rbtree.Node{Price: 5,
				Right: &rbtree.Node{Price: 7},
				Left: &rbtree.Node{Price: 3,
					Right: &rbtree.Node{Price: 4, Red: true}}}),
			getNode: func(t *rbtree.Tree) *rbtree.Node {
				return t.Root.Left.Right
			},
			wantTree: tree(&rbtree.Node{Price: 5,
				Right: &rbtree.Node{Price: 7},
				Left:  &rbtree.Node{Price: 3}}),
		},
		{
			name: "right_left_leaf",
			tree: tree(&rbtree.Node{Price: 5,
				Right: &rbtree.Node{Price: 7,
					Left: &rbtree.Node{Price: 6, Red: true}},
				Left: &rbtree.Node{Price: 3}}),
			getNode: func(t *rbtree.Tree) *rbtree.Node {
				return t.Root.Right.Left
			},
			wantTree: tree(&rbtree.Node{Price: 5,
				Right: &rbtree.Node{Price: 7},
				Left:  &rbtree.Node{Price: 3}}),
		},
		{
			name: "right_right_leaf",
			tree: tree(&rbtree.Node{Price: 5,
				Right: &rbtree.Node{Price: 7,
					Right: &rbtree.Node{Price: 8, Red: true}},
				Left: &rbtree.Node{Price: 3}}),
			getNode: func(t *rbtree.Tree) *rbtree.Node {
				return t.Root.Right.Right
			},
			wantTree: tree(&rbtree.Node{Price: 5,
				Right: &rbtree.Node{Price: 7},
				Left:  &rbtree.Node{Price: 3}}),
		},
		{
			name: "left_line_extra_node",
			tree: tree(&rbtree.Node{Price: 5,
				Right: &rbtree.Node{Price: 7},
				Left: &rbtree.Node{Price: 3,
					Left: &rbtree.Node{Price: 1, Red: true}}}),
			getNode: func(t *rbtree.Tree) *rbtree.Node {
				return t.Root.Left.Left
			},
			wantTree: tree(&rbtree.Node{Price: 5,
				Right: &rbtree.Node{Price: 7},
				Left:  &rbtree.Node{Price: 3}}),
		},
		{
			name: "left_arrow_extra_node",
			tree: tree(&rbtree.Node{Price: 5,
				Right: &rbtree.Node{Price: 7},
				Left: &rbtree.Node{Price: 3,
					Right: &rbtree.Node{Price: 4, Red: true}}}),
			getNode: func(t *rbtree.Tree) *rbtree.Node {
				return t.Root.Left.Right
			},
			wantTree: tree(&rbtree.Node{Price: 5,
				Right: &rbtree.Node{Price: 7},
				Left:  &rbtree.Node{Price: 3}}),
		},
		{
			name: "right_line_extra_node",
			tree: tree(&rbtree.Node{Price: 5,
				Right: &rbtree.Node{Price: 7,
					Right: &rbtree.Node{Price: 8, Red: true}},
				Left: &rbtree.Node{Price: 3}}),
			getNode: func(t *rbtree.Tree) *rbtree.Node {
				return t.Root.Right.Right
			},
			wantTree: tree(&rbtree.Node{Price: 5,
				Right: &rbtree.Node{Price: 7},
				Left:  &rbtree.Node{Price: 3}}),
		},
		{
			name: "right_arrow_extra_node",
			tree: tree(&rbtree.Node{Price: 5,
				Right: &rbtree.Node{Price: 7,
					Left: &rbtree.Node{Price: 6, Red: true}},
				Left: &rbtree.Node{Price: 3}}),
			getNode: func(t *rbtree.Tree) *rbtree.Node {
				return t.Root.Right.Left
			},
			wantTree: tree(&rbtree.Node{Price: 5,
				Right: &rbtree.Node{Price: 7},
				Left:  &rbtree.Node{Price: 3}}),
		},
		{
			name: "black_no_children_right_far_red_nephew",
			tree: tree(&rbtree.Node{Price: 2,
				Right: &rbtree.Node{Price: 3,
					Right: &rbtree.Node{Price: 4, Red: true}},
				Left: &rbtree.Node{Price: 1}}),
			getNode: func(t *rbtree.Tree) *rbtree.Node {
				return t.Root.Left
			},
			wantTree: tree(&rbtree.Node{Price: 3,
				Right: &rbtree.Node{Price: 4},
				Left:  &rbtree.Node{Price: 2}}),
		},
		{
			name: "black_no_children_right_close_red_nephew",
			tree: tree(&rbtree.Node{Price: 2,
				Right: &rbtree.Node{Price: 4,
					Left: &rbtree.Node{Price: 3, Red: true}},
				Left: &rbtree.Node{Price: 1}}),
			getNode: func(t *rbtree.Tree) *rbtree.Node {
				return t.Root.Left
			},
			wantTree: tree(&rbtree.Node{Price: 3,
				Right: &rbtree.Node{Price: 4},
				Left:  &rbtree.Node{Price: 2}}),
		},
		{
			name: "black_no_children_left_far_red_nephew",
			tree: tree(&rbtree.Node{Price: 3,
				Right: &rbtree.Node{Price: 4},
				Left: &rbtree.Node{Price: 2,
					Left: &rbtree.Node{Price: 1, Red: true}}}),
			getNode: func(t *rbtree.Tree) *rbtree.Node {
				return t.Root.Right
			},
			wantTree: tree(&rbtree.Node{Price: 2,
				Right: &rbtree.Node{Price: 3},
				Left:  &rbtree.Node{Price: 1}}),
		},
		{
			name: "black_no_children_left_close_red_nephew",
			tree: tree(&rbtree.Node{Price: 3,
				Right: &rbtree.Node{Price: 4},
				Left: &rbtree.Node{Price: 1,
					Right: &rbtree.Node{Price: 2, Red: true}}}),
			getNode: func(t *rbtree.Tree) *rbtree.Node {
				return t.Root.Right
			},
			wantTree: tree(&rbtree.Node{Price: 2,
				Right: &rbtree.Node{Price: 3},
				Left:  &rbtree.Node{Price: 1}}),
		},
		{
			name: "black_no_children_black_nephews",
			tree: tree(&rbtree.Node{Price: 4,
				Right: &rbtree.Node{Price: 6,
					Right: &rbtree.Node{Price: 8, Red: true,
						Right: &rbtree.Node{Price: 9,
							Right: &rbtree.Node{Price: 10, Red: true}},
						Left: &rbtree.Node{Price: 7}},
					Left: &rbtree.Node{Price: 5}},
				Left: &rbtree.Node{Price: 2,
					Right: &rbtree.Node{Price: 3},
					Left:  &rbtree.Node{Price: 1}}}),
			getNode: func(t *rbtree.Tree) *rbtree.Node {
				return t.Root.Left.Left
			},
			wantTree: tree(&rbtree.Node{Price: 6,
				Right: &rbtree.Node{Price: 8,
					Right: &rbtree.Node{Price: 9,
						Right: &rbtree.Node{Price: 10, Red: true}},
					Left: &rbtree.Node{Price: 7}},
				Left: &rbtree.Node{Price: 4,
					Right: &rbtree.Node{Price: 5},
					Left: &rbtree.Node{Price: 2,
						Right: &rbtree.Node{Price: 3, Red: true}}}}),
		},
		{
			name: "two_level_black_tree_delete_rightmost",
			tree: tree(&rbtree.Node{Price: 4,
				Right: &rbtree.Node{Price: 6,
					Right: &rbtree.Node{Price: 7},
					Left:  &rbtree.Node{Price: 5}},
				Left: &rbtree.Node{Price: 2,
					Right: &rbtree.Node{Price: 3},
					Left:  &rbtree.Node{Price: 1}}}),
			getNode: func(t *rbtree.Tree) *rbtree.Node {
				return t.Root.Right.Right
			},
			wantTree: tree(&rbtree.Node{Price: 4,
				Right: &rbtree.Node{Price: 6,
					Left: &rbtree.Node{Price: 5, Red: true}},
				Left: &rbtree.Node{Price: 2, Red: true,
					Right: &rbtree.Node{Price: 3},
					Left:  &rbtree.Node{Price: 1}}}),
		},
		{
			name: "random_case_1",
			tree: tree(&rbtree.Node{Price: 8,
				Right: &rbtree.Node{Price: 10,
					Right: &rbtree.Node{Price: 12},
					Left:  &rbtree.Node{Price: 9}},
				Left: &rbtree.Node{Price: 4,
					Right: &rbtree.Node{Price: 6,
						Right: &rbtree.Node{Price: 7, Red: true}},
					Left: &rbtree.Node{Price: 3,
						Left: &rbtree.Node{Price: 1, Red: true}}}}),
			getNode: func(t *rbtree.Tree) *rbtree.Node {
				return t.Root
			},
			wantTree: tree(&rbtree.Node{Price: 9,
				Right: &rbtree.Node{Price: 10,
					Right: &rbtree.Node{Price: 12, Red: true}},
				Left: &rbtree.Node{Price: 4, Red: true,
					Right: &rbtree.Node{Price: 6,
						Right: &rbtree.Node{Price: 7, Red: true}},
					Left: &rbtree.Node{Price: 3,
						Left: &rbtree.Node{Price: 1, Red: true}}}}),
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

			node := tc.getNode(tc.tree)
			tc.tree.DeleteNode(node)

			if err := tc.tree.Valid(); err != nil {
				t.Fatalf("Invalid got tree: %v\ngot:\n%s", err, tc.tree.PrintTree())
			}

			gotTree := tc.tree.PrintTree()
			wantTree := tc.wantTree.PrintTree()
			if gotTree != wantTree {
				t.Errorf("DeleteNode(%v)\nwant tree:\n%s\ngot tree:\n%s\n", node, wantTree, gotTree)
			}
		})
	}
}

func Test_Delete_Head(t *testing.T) {
	testCases := []struct {
		name          string
		orientation   rbtree.TreeOrientation
		insertions    []uint64
		deletions     []uint64
		wantHead      bool
		wantHeadPrice uint64
	}{
		{
			name:        "min_empty",
			insertions:  []uint64{1},
			deletions:   []uint64{1},
			orientation: rbtree.MinFirst,
		},
		{
			name:        "max_empty",
			insertions:  []uint64{1},
			deletions:   []uint64{1},
			orientation: rbtree.MaxFirst,
		},
		{
			name:          "min_single",
			orientation:   rbtree.MinFirst,
			insertions:    []uint64{1, 2},
			deletions:     []uint64{2},
			wantHead:      true,
			wantHeadPrice: 1,
		},
		{
			name:          "max_single",
			orientation:   rbtree.MaxFirst,
			insertions:    []uint64{2, 1},
			deletions:     []uint64{1},
			wantHead:      true,
			wantHeadPrice: 2,
		},
		{
			name:          "min_keep_previous",
			orientation:   rbtree.MinFirst,
			insertions:    []uint64{1, 2},
			deletions:     []uint64{2},
			wantHead:      true,
			wantHeadPrice: 1,
		},
		{
			name:          "max_keep_previous",
			orientation:   rbtree.MaxFirst,
			insertions:    []uint64{2, 1},
			deletions:     []uint64{1},
			wantHead:      true,
			wantHeadPrice: 2,
		},
		{
			name:          "min_override",
			orientation:   rbtree.MinFirst,
			insertions:    []uint64{5, 4, 3, 2, 1},
			deletions:     []uint64{1, 2},
			wantHead:      true,
			wantHeadPrice: 3,
		},
		{
			name:          "max_override",
			orientation:   rbtree.MaxFirst,
			insertions:    []uint64{1, 2, 3, 4, 5},
			deletions:     []uint64{5, 4},
			wantHead:      true,
			wantHeadPrice: 3,
		},
		{
			name:          "min_multiple",
			orientation:   rbtree.MinFirst,
			insertions:    []uint64{1, 2, 3, 4, 5, 6, 7, 8},
			deletions:     []uint64{8, 1, 7, 2},
			wantHead:      true,
			wantHeadPrice: 3,
		},
		{
			name:          "max_multiple",
			orientation:   rbtree.MaxFirst,
			insertions:    []uint64{1, 2, 3, 4, 5, 6, 7, 8},
			deletions:     []uint64{8, 1, 7, 2},
			wantHead:      true,
			wantHeadPrice: 6,
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

			for _, num := range tc.deletions {
				tree.Delete(num)
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

func Test_Delete_Sequence(t *testing.T) {
	testCases := []struct {
		name   string
		insert func(uint64) []uint64
		delete func(uint64) []uint64
	}{
		{
			name:   "insert_asc_delete_asc",
			insert: testutils.NumbersAscending,
			delete: testutils.NumbersAscending,
		},
		{
			name:   "insert_asc_delete_desc",
			insert: testutils.NumbersAscending,
			delete: testutils.NumbersDescending,
		},
		{
			name:   "insert_asc_delete_random",
			insert: testutils.NumbersAscending,
			delete: testutils.NumbersRandom,
		},
		{
			name:   "insert_desc_delete_asc",
			insert: testutils.NumbersDescending,
			delete: testutils.NumbersAscending,
		},
		{
			name:   "insert_desc_delete_desc",
			insert: testutils.NumbersDescending,
			delete: testutils.NumbersDescending,
		},
		{
			name:   "insert_desc_delete_random",
			insert: testutils.NumbersDescending,
			delete: testutils.NumbersRandom,
		},
		{
			name:   "insert_random_delete_asc",
			insert: testutils.NumbersRandom,
			delete: testutils.NumbersAscending,
		},
		{
			name:   "insert_random_delete_desc",
			insert: testutils.NumbersRandom,
			delete: testutils.NumbersDescending,
		},
		{
			name:   "insert_random_delete_random",
			insert: testutils.NumbersRandom,
			delete: testutils.NumbersRandom,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tr := tree(nil)
			for _, num := range tc.insert(100) {
				start := tr.PrintTree()

				tr.Insert(num)

				if err := tr.Valid(); err != nil {
					t.Fatalf("Insert(%v) invalid tree: %v\nStart:\n%s\nGot:\n%s", num, err, start, tr.PrintTree())
				}
			}

			for _, num := range tc.delete(100) {
				start := tr.PrintTree()

				tr.Delete(num)

				if err := tr.Valid(); err != nil {
					t.Fatalf("Delete(%v) invalid tree: %v\nStart:\n%s\nGot:\n%s", num, err, start, tr.PrintTree())
				}
			}
		})
	}
}
