package rbtree_test

import (
	"exchange/engine/orderbook/rbtree"
	"sync"
	"testing"
)

func connectParents(t *rbtree.Tree) {
	if t == nil || t.Root == nil {
		return
	}

	visited := []*rbtree.Node{t.Root}
	for {
		if len(visited) == 0 {
			return
		}

		node := visited[0]
		visited = visited[1:]

		if node.Left != nil {
			node.Left.Parent = node
			visited = append(visited, node.Left)
		}

		if node.Right != nil {
			node.Right.Parent = node
			visited = append(visited, node.Right)
		}
	}
}

// tree is an easy way to create trees in tests
func tree(root *rbtree.Node) *rbtree.Tree {
	pool := &sync.Pool{
		New: func() any {
			return rbtree.NewNode()
		},
	}
	t := rbtree.NewTree(rbtree.MinFirst, pool)
	t.Root = root
	return t
}

func Test_connectParents(t *testing.T) {
	testCases := []struct {
		name string
		tree *rbtree.Tree
		want func() *rbtree.Tree
	}{
		{
			name: "empty",
			tree: &rbtree.Tree{},
			want: func() *rbtree.Tree {
				return &rbtree.Tree{}
			},
		},
		{
			name: "few_nodes",
			tree: &rbtree.Tree{
				Root: &rbtree.Node{
					Price: 2,
					Left:  &rbtree.Node{Price: 1},
					Right: &rbtree.Node{Price: 3},
				},
			},
			want: func() *rbtree.Tree {
				root := &rbtree.Node{Price: 2}
				root.Left = &rbtree.Node{Price: 1}
				root.Left.Parent = root

				root.Right = &rbtree.Node{Price: 3}
				root.Right.Parent = root
				return &rbtree.Tree{Root: root}
			},
		},
		{
			name: "multiple_nodes",
			tree: &rbtree.Tree{
				Root: &rbtree.Node{
					Price: 10,
					Left: &rbtree.Node{
						Price: 5,
						Left: &rbtree.Node{
							Price: 3,
							Left:  &rbtree.Node{Price: 1},
							Right: &rbtree.Node{Price: 4},
						},
						Right: &rbtree.Node{
							Price: 7,
							Right: &rbtree.Node{Price: 8},
						},
					},
					Right: &rbtree.Node{
						Price: 15,
						Left: &rbtree.Node{
							Price: 13,
							Left:  &rbtree.Node{Price: 11},
						},
						Right: &rbtree.Node{
							Price: 16,
							Right: &rbtree.Node{Price: 17},
						},
					},
				},
			},
			want: func() *rbtree.Tree {
				root := &rbtree.Node{Price: 10}

				rl := &rbtree.Node{Price: 5}
				root.Left = rl
				rl.Parent = root

				rll := &rbtree.Node{Price: 3}
				rl.Left = rll
				rll.Parent = rl
				rll.Left = &rbtree.Node{Price: 1}
				rll.Left.Parent = rll
				rll.Right = &rbtree.Node{Price: 4}
				rll.Right.Parent = rll

				rlr := &rbtree.Node{Price: 7}
				rl.Right = rlr
				rlr.Parent = rl
				rlr.Right = &rbtree.Node{Price: 8}
				rlr.Right.Parent = rlr

				rr := &rbtree.Node{Price: 15}
				root.Right = rr
				rr.Parent = root

				rrl := &rbtree.Node{Price: 13}
				rr.Left = rrl
				rrl.Parent = rr
				rrl.Left = &rbtree.Node{Price: 11}
				rrl.Left.Parent = rrl
				rrr := &rbtree.Node{Price: 16}
				rr.Right = rrr
				rrr.Parent = rr
				rrr.Right = &rbtree.Node{Price: 17}
				rrr.Right.Parent = rrr

				return &rbtree.Tree{Root: root}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			connectParents(tc.tree)

			got := tc.tree.PrintTree()
			want := tc.want().PrintTree()
			if got != want {
				t.Errorf("connectParents() error.\nwant:\n%s\ngot:\n%s", want, got)
			}
		})
	}
}
