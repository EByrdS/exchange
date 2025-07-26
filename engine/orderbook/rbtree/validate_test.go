package rbtree_test

import (
	"exchange/engine/orderbook/rbtree"
	"testing"
)

func Test_Valid(t *testing.T) {
	testCases := []struct {
		name   string
		tree   *rbtree.Tree
		wantOk bool
	}{
		{
			name: "nil",
		},
		{
			name:   "empty",
			tree:   &rbtree.Tree{},
			wantOk: true,
		},
		{
			name: "red_root",
			tree: &rbtree.Tree{Root: &rbtree.Node{Price: 1, Red: true}},
		},
		{
			name: "larger_left",
			tree: &rbtree.Tree{
				Root: &rbtree.Node{
					Price: 1,
					Left:  &rbtree.Node{Price: 2},
				},
			},
		},
		{
			name: "smaller_right",
			tree: &rbtree.Tree{
				Root: &rbtree.Node{
					Price: 2,
					Right: &rbtree.Node{Price: 1},
				},
			},
		},
		{
			name: "equal_left",
			tree: &rbtree.Tree{
				Root: &rbtree.Node{
					Price: 1,
					Left:  &rbtree.Node{Price: 1},
				},
			},
		},
		{
			name: "equal_right",
			tree: &rbtree.Tree{
				Root: &rbtree.Node{
					Price: 1,
					Right: &rbtree.Node{Price: 1},
				},
			},
		},
		{
			name: "non_root_larger_left",
			tree: &rbtree.Tree{
				Root: &rbtree.Node{
					Price: 2,
					Left:  &rbtree.Node{Price: 99}, // fail
					Right: &rbtree.Node{
						Price: 4,
						Red:   true,
						Left:  &rbtree.Node{Price: 3},
						Right: &rbtree.Node{Price: 5},
					},
				},
			},
		},
		{
			name: "non_root_smaller_right",
			tree: &rbtree.Tree{
				Root: &rbtree.Node{
					Price: 2,
					Left:  &rbtree.Node{Price: 1},
					Right: &rbtree.Node{
						Price: 4,
						Red:   true,
						Left:  &rbtree.Node{Price: 3},
						Right: &rbtree.Node{Price: 0}, // fail
					},
				},
			},
		},
		{
			name: "red_node_red_left",
			tree: &rbtree.Tree{
				Root: &rbtree.Node{
					Price: 2,
					Left:  &rbtree.Node{Price: 1},
					Right: &rbtree.Node{
						Price: 4,
						Red:   true,
						Left:  &rbtree.Node{Price: 3, Red: true}, // fail
						Right: &rbtree.Node{Price: 5},
					},
				},
			},
		},
		{
			name: "red_node_red_right",
			tree: &rbtree.Tree{
				Root: &rbtree.Node{
					Price: 2,
					Left:  &rbtree.Node{Price: 1},
					Right: &rbtree.Node{
						Price: 4,
						Red:   true,
						Left:  &rbtree.Node{Price: 3},
						Right: &rbtree.Node{Price: 5, Red: true}, // fail
					},
				},
			},
		},
		{
			name: "two_levels",
			tree: &rbtree.Tree{
				Root: &rbtree.Node{
					Price: 2,
					Left:  &rbtree.Node{Price: 1},
					Right: &rbtree.Node{Price: 3},
				},
			},
			wantOk: true,
		},
		{
			name: "two_levels_with_red",
			tree: &rbtree.Tree{
				Root: &rbtree.Node{
					Price: 2,
					Left:  &rbtree.Node{Price: 1, Red: true},
					Right: &rbtree.Node{Price: 3, Red: true},
				},
			},
			wantOk: true,
		},
		{
			name: "unequal_black_counts",
			tree: &rbtree.Tree{
				Root: &rbtree.Node{
					Price: 2,
					Left:  &rbtree.Node{Price: 1},
					Right: &rbtree.Node{
						Price: 4,
						Left:  &rbtree.Node{Price: 3}, // should be red
					},
				},
			},
		},
		{
			name: "three_levels_one_red",
			tree: &rbtree.Tree{
				Root: &rbtree.Node{
					Price: 2,
					Left:  &rbtree.Node{Price: 1},
					Right: &rbtree.Node{
						Price: 4,
						Left:  &rbtree.Node{Price: 3, Red: true},
					},
				},
			},
			wantOk: true,
		},
		{
			name: "three_levels_two_reds",
			tree: &rbtree.Tree{
				Root: &rbtree.Node{
					Price: 2,
					Left:  &rbtree.Node{Price: 1},
					Right: &rbtree.Node{
						Price: 4,
						Left:  &rbtree.Node{Price: 3, Red: true},
						Right: &rbtree.Node{Price: 5, Red: true},
					},
				},
			},
			wantOk: true,
		},
		{
			name: "three_levels_one_middle_red",
			tree: &rbtree.Tree{
				Root: &rbtree.Node{
					Price: 2,
					Left:  &rbtree.Node{Price: 1},
					Right: &rbtree.Node{
						Price: 4,
						Red:   true,
						Left:  &rbtree.Node{Price: 3},
						Right: &rbtree.Node{Price: 5},
					},
				},
			},
			wantOk: true,
		},
		{
			name: "four_levels_balanced",
			tree: &rbtree.Tree{
				Root: &rbtree.Node{
					Price: 16,
					Left: &rbtree.Node{
						Price: 8,
						Red:   true,
						Left:  &rbtree.Node{Price: 3},
						Right: &rbtree.Node{
							Price: 9,
							Right: &rbtree.Node{Price: 10, Red: true},
						},
					},
					Right: &rbtree.Node{
						Price: 21,
						Red:   true,
						Left:  &rbtree.Node{Price: 18},
						Right: &rbtree.Node{
							Price: 27,
							Right: &rbtree.Node{Price: 29, Red: true},
						},
					},
				},
			},
			wantOk: true,
		},
		{
			// https://en.wikipedia.org/wiki/Red%E2%80%93black_tree
			name: "four_levels_example",
			tree: &rbtree.Tree{
				Root: &rbtree.Node{
					Price: 13,
					Left: &rbtree.Node{
						Price: 8,
						Red:   true,
						Left: &rbtree.Node{
							Price: 1,
							Right: &rbtree.Node{
								Price: 6,
								Red:   true,
							},
						},
						Right: &rbtree.Node{
							Price: 11,
						},
					},
					Right: &rbtree.Node{
						Price: 17,
						Red:   true,
						Left: &rbtree.Node{
							Price: 15,
						},
						Right: &rbtree.Node{
							Price: 25,
							Left: &rbtree.Node{
								Price: 22,
								Red:   true,
							},
							Right: &rbtree.Node{
								Price: 27,
								Red:   true,
							},
						},
					},
				},
			},
			wantOk: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			connectParents(tc.tree)

			err := tc.tree.Valid()
			if err != nil {
				if tc.wantOk {
					t.Errorf("Valid() wanted ok, got error: %v", err)
				}
			} else {
				if !tc.wantOk {
					t.Errorf("Valid() wanted NOT ok, but did not raise")
				}
			}
		})
	}
}
