package rbtree_test

import (
	"exchange/engine/orderbook/rbtree"
	"testing"
)

func Test_RotateLeft(t *testing.T) {
	testCases := []struct {
		name   string
		tree   func() *rbtree.Tree
		rotate func(tree *rbtree.Tree) *rbtree.Node
		want   *rbtree.Tree
	}{
		{
			name: "nil",
			tree: func() *rbtree.Tree {
				return &rbtree.Tree{}
			},
			rotate: func(tree *rbtree.Tree) *rbtree.Node {
				return tree.Root
			},
			want: &rbtree.Tree{},
		},
		{
			name: "no_right_child",
			tree: func() *rbtree.Tree {
				return &rbtree.Tree{Root: &rbtree.Node{Price: 1}}
			},
			rotate: func(tree *rbtree.Tree) *rbtree.Node {
				return tree.Root
			},
			want: &rbtree.Tree{Root: &rbtree.Node{Price: 1}},
		},
		{
			name: "root_no_subtree",
			tree: func() *rbtree.Tree {
				return &rbtree.Tree{Root: &rbtree.Node{Price: 1, Right: &rbtree.Node{Price: 2}}}
			},
			rotate: func(tree *rbtree.Tree) *rbtree.Node {
				return tree.Root
			},
			want: &rbtree.Tree{Root: &rbtree.Node{Price: 2, Left: &rbtree.Node{Price: 1}}},
		},
		{
			name: "root_subtree",
			tree: func() *rbtree.Tree {
				return &rbtree.Tree{Root: &rbtree.Node{Price: 1, Right: &rbtree.Node{Price: 2, Left: &rbtree.Node{Price: 3}}}}
			},
			rotate: func(tree *rbtree.Tree) *rbtree.Node {
				return tree.Root
			},
			want: &rbtree.Tree{Root: &rbtree.Node{Price: 2, Left: &rbtree.Node{Price: 1, Right: &rbtree.Node{Price: 3}}}},
		},
		{
			name: "root_multiple_subtrees",
			tree: func() *rbtree.Tree {
				return &rbtree.Tree{
					Root: &rbtree.Node{Price: 1, // x
						Left: &rbtree.Node{Price: 5}, // a
						Right: &rbtree.Node{Price: 2, // y
							Left:  &rbtree.Node{Price: 4},  // _b
							Right: &rbtree.Node{Price: 3}}, // _y
					},
				}
			},
			rotate: func(tree *rbtree.Tree) *rbtree.Node {
				return tree.Root
			},
			want: &rbtree.Tree{
				Root: &rbtree.Node{Price: 2, // y
					Left: &rbtree.Node{Price: 1, // x
						Left:  &rbtree.Node{Price: 5},  // a
						Right: &rbtree.Node{Price: 4}}, // _b
					Right: &rbtree.Node{Price: 3}, // _y
				},
			},
		},
		{
			name: "left_child_no_subtree",
			tree: func() *rbtree.Tree {
				tr := &rbtree.Tree{}
				tr.Root = &rbtree.Node{Price: 10}
				tr.Root.Left = &rbtree.Node{Price: 1, Right: &rbtree.Node{Price: 2}}
				tr.Root.Left.Parent = tr.Root
				return tr
			},
			rotate: func(tree *rbtree.Tree) *rbtree.Node {
				return tree.Root.Left
			},
			want: &rbtree.Tree{
				Root: &rbtree.Node{Price: 10,
					Left: &rbtree.Node{Price: 2, Left: &rbtree.Node{Price: 1}}}},
		},
		{
			name: "right_child_no_subtree",
			tree: func() *rbtree.Tree {
				tr := &rbtree.Tree{}
				tr.Root = &rbtree.Node{Price: 10}
				tr.Root.Right = &rbtree.Node{Price: 1, Right: &rbtree.Node{Price: 2}}
				tr.Root.Right.Parent = tr.Root
				return tr
			},
			rotate: func(tree *rbtree.Tree) *rbtree.Node {
				return tree.Root.Right
			},
			want: &rbtree.Tree{
				Root: &rbtree.Node{Price: 10,
					Right: &rbtree.Node{Price: 2, Left: &rbtree.Node{Price: 1}}}},
		},
		{
			name: "left_child_subtree",
			tree: func() *rbtree.Tree {
				tr := &rbtree.Tree{}
				tr.Root = &rbtree.Node{Price: 10}
				tr.Root.Left = &rbtree.Node{Price: 1, Right: &rbtree.Node{Price: 2, Left: &rbtree.Node{Price: 3}}}
				tr.Root.Left.Parent = tr.Root
				return tr
			},
			rotate: func(tree *rbtree.Tree) *rbtree.Node {
				return tree.Root.Left
			},
			want: &rbtree.Tree{Root: &rbtree.Node{Price: 10, Left: &rbtree.Node{Price: 2, Left: &rbtree.Node{Price: 1, Right: &rbtree.Node{Price: 3}}}}},
		},
		{
			name: "right_child_subtree",
			tree: func() *rbtree.Tree {
				tr := &rbtree.Tree{}
				tr.Root = &rbtree.Node{Price: 10}
				tr.Root.Right = &rbtree.Node{Price: 1, Right: &rbtree.Node{Price: 2, Left: &rbtree.Node{Price: 3}}}
				tr.Root.Right.Parent = tr.Root
				return tr
			},
			rotate: func(tree *rbtree.Tree) *rbtree.Node {
				return tree.Root.Right
			},
			want: &rbtree.Tree{Root: &rbtree.Node{Price: 10, Right: &rbtree.Node{Price: 2, Left: &rbtree.Node{Price: 1, Right: &rbtree.Node{Price: 3}}}}},
		},
		{
			name: "left_child_multiple_subtrees",
			tree: func() *rbtree.Tree {
				tr := &rbtree.Tree{}
				tr.Root = &rbtree.Node{Price: 10}
				tr.Root.Left = &rbtree.Node{Price: 1, // x
					Left: &rbtree.Node{Price: 5}, // a
					Right: &rbtree.Node{Price: 2, // y
						Left:  &rbtree.Node{Price: 4},  // _b
						Right: &rbtree.Node{Price: 3}}, // _y
				}
				tr.Root.Left.Parent = tr.Root
				return tr
			},
			rotate: func(tree *rbtree.Tree) *rbtree.Node {
				return tree.Root.Left
			},
			want: &rbtree.Tree{Root: &rbtree.Node{Price: 10, Left: &rbtree.Node{Price: 2, // y
				Left: &rbtree.Node{Price: 1, // x
					Left:  &rbtree.Node{Price: 5},  // a
					Right: &rbtree.Node{Price: 4}}, // _b
				Right: &rbtree.Node{Price: 3}, // _y
			}}},
		},
		{
			name: "right_child_multiple_subtrees",
			tree: func() *rbtree.Tree {
				tr := &rbtree.Tree{}
				tr.Root = &rbtree.Node{Price: 10}
				tr.Root.Right = &rbtree.Node{Price: 1, // x
					Left: &rbtree.Node{Price: 5}, // a
					Right: &rbtree.Node{Price: 2, // y
						Left:  &rbtree.Node{Price: 4},  // _b
						Right: &rbtree.Node{Price: 3}}, // _y
				}
				tr.Root.Right.Parent = tr.Root
				return tr
			},
			rotate: func(tree *rbtree.Tree) *rbtree.Node {
				return tree.Root.Right
			},
			want: &rbtree.Tree{Root: &rbtree.Node{Price: 10, Right: &rbtree.Node{Price: 2, // y
				Left: &rbtree.Node{Price: 1, // x
					Left:  &rbtree.Node{Price: 5},  // a
					Right: &rbtree.Node{Price: 4}}, // _b
				Right: &rbtree.Node{Price: 3}, // _y
			}}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tree := tc.tree()
			connectParents(tree)

			tree.RotateLeft(tc.rotate(tree))

			connectParents(tc.want)
			want := tc.want.PrintTree()
			got := tree.PrintTree()
			if want != got {
				t.Errorf("tree.RotateLeft() error.\nwant:\n%s,\ngot:\n%s", want, got)
			}
		})
	}
}

func Test_RotateRight(t *testing.T) {
	testCases := []struct {
		name   string
		tree   func() *rbtree.Tree
		rotate func(tree *rbtree.Tree) *rbtree.Node
		want   *rbtree.Tree
	}{
		{
			name: "nil",
			tree: func() *rbtree.Tree {
				return &rbtree.Tree{}
			},
			rotate: func(tree *rbtree.Tree) *rbtree.Node {
				return tree.Root
			},
			want: &rbtree.Tree{},
		},
		{
			name: "no_left_child",
			tree: func() *rbtree.Tree {
				return &rbtree.Tree{Root: &rbtree.Node{Price: 1}}
			},
			rotate: func(tree *rbtree.Tree) *rbtree.Node {
				return tree.Root
			},
			want: &rbtree.Tree{Root: &rbtree.Node{Price: 1}},
		},
		{
			name: "root_no_subtree",
			tree: func() *rbtree.Tree {
				return &rbtree.Tree{Root: &rbtree.Node{Price: 1, Left: &rbtree.Node{Price: 2}}}
			},
			rotate: func(tree *rbtree.Tree) *rbtree.Node {
				return tree.Root
			},
			want: &rbtree.Tree{Root: &rbtree.Node{Price: 2, Right: &rbtree.Node{Price: 1}}},
		},
		{
			name: "root_subtree",
			tree: func() *rbtree.Tree {
				return &rbtree.Tree{Root: &rbtree.Node{Price: 1, Left: &rbtree.Node{Price: 2, Right: &rbtree.Node{Price: 3}}}}
			},
			rotate: func(tree *rbtree.Tree) *rbtree.Node {
				return tree.Root
			},
			want: &rbtree.Tree{Root: &rbtree.Node{Price: 2, Right: &rbtree.Node{Price: 1, Left: &rbtree.Node{Price: 3}}}},
		},
		{
			name: "root_multiple_subtrees",
			tree: func() *rbtree.Tree {
				return &rbtree.Tree{
					Root: &rbtree.Node{Price: 2, // y
						Left: &rbtree.Node{Price: 1, // x
							Left:  &rbtree.Node{Price: 5},  // a
							Right: &rbtree.Node{Price: 4}}, // _b
						Right: &rbtree.Node{Price: 3}, // _y
					},
				}
			},
			rotate: func(tree *rbtree.Tree) *rbtree.Node {
				return tree.Root
			},
			want: &rbtree.Tree{
				Root: &rbtree.Node{Price: 1, // x
					Left: &rbtree.Node{Price: 5}, // a
					Right: &rbtree.Node{Price: 2, // y
						Left:  &rbtree.Node{Price: 4},  // _b
						Right: &rbtree.Node{Price: 3}}, // _y
				},
			},
		},
		{
			name: "left_child_no_subtree",
			tree: func() *rbtree.Tree {
				tr := &rbtree.Tree{}
				tr.Root = &rbtree.Node{Price: 10}
				tr.Root.Left = &rbtree.Node{Price: 1, Left: &rbtree.Node{Price: 2}}
				tr.Root.Left.Parent = tr.Root
				return tr
			},
			rotate: func(tree *rbtree.Tree) *rbtree.Node {
				return tree.Root.Left
			},
			want: &rbtree.Tree{
				Root: &rbtree.Node{Price: 10,
					Left: &rbtree.Node{Price: 2, Right: &rbtree.Node{Price: 1}}}},
		},
		{
			name: "right_child_no_subtree",
			tree: func() *rbtree.Tree {
				tr := &rbtree.Tree{}
				tr.Root = &rbtree.Node{Price: 10}
				tr.Root.Right = &rbtree.Node{Price: 1, Left: &rbtree.Node{Price: 2}}
				tr.Root.Right.Parent = tr.Root
				return tr
			},
			rotate: func(tree *rbtree.Tree) *rbtree.Node {
				return tree.Root.Right
			},
			want: &rbtree.Tree{
				Root: &rbtree.Node{Price: 10,
					Right: &rbtree.Node{Price: 2, Right: &rbtree.Node{Price: 1}}}},
		},
		{
			name: "left_child_subtree",
			tree: func() *rbtree.Tree {
				tr := &rbtree.Tree{}
				tr.Root = &rbtree.Node{Price: 10}
				tr.Root.Left = &rbtree.Node{Price: 1, Left: &rbtree.Node{Price: 2, Right: &rbtree.Node{Price: 3}}}
				tr.Root.Left.Parent = tr.Root
				return tr
			},
			rotate: func(tree *rbtree.Tree) *rbtree.Node {
				return tree.Root.Left
			},
			want: &rbtree.Tree{Root: &rbtree.Node{Price: 10, Left: &rbtree.Node{Price: 2, Right: &rbtree.Node{Price: 1, Left: &rbtree.Node{Price: 3}}}}},
		},
		{
			name: "right_child_subtree",
			tree: func() *rbtree.Tree {
				tr := &rbtree.Tree{}
				tr.Root = &rbtree.Node{Price: 10}
				tr.Root.Right = &rbtree.Node{Price: 1, Left: &rbtree.Node{Price: 2, Right: &rbtree.Node{Price: 3}}}
				tr.Root.Right.Parent = tr.Root
				return tr
			},
			rotate: func(tree *rbtree.Tree) *rbtree.Node {
				return tree.Root.Right
			},
			want: &rbtree.Tree{Root: &rbtree.Node{Price: 10, Right: &rbtree.Node{Price: 2, Right: &rbtree.Node{Price: 1, Left: &rbtree.Node{Price: 3}}}}},
		},
		{
			name: "left_child_multiple_subtrees",
			tree: func() *rbtree.Tree {
				tr := &rbtree.Tree{}
				tr.Root = &rbtree.Node{Price: 10}
				tr.Root.Left = &rbtree.Node{Price: 2, // y
					Left: &rbtree.Node{Price: 1, // x
						Left:  &rbtree.Node{Price: 5},  // a
						Right: &rbtree.Node{Price: 4}}, // _b
					Right: &rbtree.Node{Price: 3}, // _y
				}
				tr.Root.Left.Parent = tr.Root
				return tr
			},
			rotate: func(tree *rbtree.Tree) *rbtree.Node {
				return tree.Root.Left
			},
			want: &rbtree.Tree{Root: &rbtree.Node{Price: 10,
				Left: &rbtree.Node{Price: 1, // x
					Left: &rbtree.Node{Price: 5}, // a
					Right: &rbtree.Node{Price: 2, // y
						Left:  &rbtree.Node{Price: 4},  // _b
						Right: &rbtree.Node{Price: 3}}, // _y
				}}},
		},
		{
			name: "right_child_multiple_subtrees",
			tree: func() *rbtree.Tree {
				tr := &rbtree.Tree{}
				tr.Root = &rbtree.Node{Price: 10}
				tr.Root.Right = &rbtree.Node{Price: 2, // y
					Left: &rbtree.Node{Price: 1, // x
						Left:  &rbtree.Node{Price: 5},  // a
						Right: &rbtree.Node{Price: 4}}, // _b
					Right: &rbtree.Node{Price: 3}, // _y
				}
				tr.Root.Right.Parent = tr.Root
				return tr
			},
			rotate: func(tree *rbtree.Tree) *rbtree.Node {
				return tree.Root.Right
			},
			want: &rbtree.Tree{Root: &rbtree.Node{Price: 10, Right: &rbtree.Node{Price: 1, // x
				Left: &rbtree.Node{Price: 5}, // a
				Right: &rbtree.Node{Price: 2, // y
					Left:  &rbtree.Node{Price: 4},  // _b
					Right: &rbtree.Node{Price: 3}}, // _y
			}}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tree := tc.tree()
			connectParents(tree)

			tree.RotateRight(tc.rotate(tree))

			connectParents(tc.want)
			want := tc.want.PrintTree()
			got := tree.PrintTree()
			if want != got {
				t.Errorf("tree.RotateRight() error.\nwant:\n%s,\ngot:\n%s", want, got)
			}
		})
	}
}
