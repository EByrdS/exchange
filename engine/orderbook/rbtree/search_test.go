package rbtree_test

import (
	"exchange/engine/orderbook/rbtree"
	"testing"
)

func Test_Min(t *testing.T) {
	testCases := []struct {
		name       string
		tree       *rbtree.Tree
		want       uint64
		wantExists bool
	}{
		{
			name: "empty",
			tree: &rbtree.Tree{},
		},
		{
			name:       "root",
			tree:       &rbtree.Tree{Root: &rbtree.Node{Price: 41}},
			want:       41,
			wantExists: true,
		},
		{
			name: "root_with_right",
			tree: &rbtree.Tree{
				Root: &rbtree.Node{
					Price: 2,
					Left:  &rbtree.Node{Price: 1, Red: true},
				},
			},
			want:       1,
			wantExists: true,
		},
		{
			name: "root_with_left",
			tree: &rbtree.Tree{
				Root: &rbtree.Node{
					Price: 1,
					Right: &rbtree.Node{Price: 2, Red: true},
				},
			},
			want:       1,
			wantExists: true,
		},
		{
			name: "one_level",
			tree: &rbtree.Tree{
				Root: &rbtree.Node{
					Price: 4,
					Left: &rbtree.Node{
						Price: 2,
						Left:  &rbtree.Node{Price: 1},
						Right: &rbtree.Node{Price: 3},
					},
					Right: &rbtree.Node{
						Price: 6,
						Left:  &rbtree.Node{Price: 5},
						Right: &rbtree.Node{Price: 7},
					},
				},
			},
			want:       1,
			wantExists: true,
		},
		{
			name: "two_levels",
			tree: &rbtree.Tree{
				Root: &rbtree.Node{
					Price: 8,
					Left: &rbtree.Node{
						Price: 4,
						Left: &rbtree.Node{
							Price: 2,
							Left:  &rbtree.Node{Price: 1},
							Right: &rbtree.Node{Price: 3},
						},
						Right: &rbtree.Node{
							Price: 6,
							Left:  &rbtree.Node{Price: 5},
							Right: &rbtree.Node{Price: 7},
						},
					},
					Right: &rbtree.Node{
						Price: 12,
						Left: &rbtree.Node{
							Price: 10,
							Left:  &rbtree.Node{Price: 9},
							Right: &rbtree.Node{Price: 11},
						},
						Right: &rbtree.Node{
							Price: 14,
							Left:  &rbtree.Node{Price: 13},
							Right: &rbtree.Node{Price: 15},
						},
					},
				},
			},
			want:       1,
			wantExists: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			connectParents(tc.tree)

			if err := tc.tree.Valid(); err != nil {
				t.Fatalf("The tree must be valid: %v", err)
			}

			got := tc.tree.Min()
			if got != nil {
				if !tc.wantExists {
					t.Fatalf("Min() expected the node NOT to exist, but it does: %v", got)
				}

				if tc.want != got.Price {
					t.Errorf("Min(), want price: %v, got node: %+v", tc.want, got)
				}
			} else if tc.wantExists {
				t.Errorf("Min() expected the node to exist, but it doesn't")
			}
		})
	}
}

func Test_Max(t *testing.T) {
	testCases := []struct {
		name       string
		tree       *rbtree.Tree
		want       uint64
		wantExists bool
	}{
		{
			name: "empty",
			tree: &rbtree.Tree{},
		},
		{
			name:       "root",
			tree:       &rbtree.Tree{Root: &rbtree.Node{Price: 41}},
			want:       41,
			wantExists: true,
		},
		{
			name: "root_with_right",
			tree: &rbtree.Tree{
				Root: &rbtree.Node{
					Price: 2,
					Left:  &rbtree.Node{Price: 1, Red: true},
				},
			},
			want:       2,
			wantExists: true,
		},
		{
			name: "root_with_left",
			tree: &rbtree.Tree{
				Root: &rbtree.Node{
					Price: 1,
					Right: &rbtree.Node{Price: 2, Red: true},
				},
			},
			want:       2,
			wantExists: true,
		},
		{
			name: "one_level",
			tree: &rbtree.Tree{
				Root: &rbtree.Node{
					Price: 4,
					Left: &rbtree.Node{
						Price: 2,
						Left:  &rbtree.Node{Price: 1},
						Right: &rbtree.Node{Price: 3},
					},
					Right: &rbtree.Node{
						Price: 6,
						Left:  &rbtree.Node{Price: 5},
						Right: &rbtree.Node{Price: 7},
					},
				},
			},
			want:       7,
			wantExists: true,
		},
		{
			name: "two_levels",
			tree: &rbtree.Tree{
				Root: &rbtree.Node{
					Price: 8,
					Left: &rbtree.Node{
						Price: 4,
						Left: &rbtree.Node{
							Price: 2,
							Left:  &rbtree.Node{Price: 1},
							Right: &rbtree.Node{Price: 3},
						},
						Right: &rbtree.Node{
							Price: 6,
							Left:  &rbtree.Node{Price: 5},
							Right: &rbtree.Node{Price: 7},
						},
					},
					Right: &rbtree.Node{
						Price: 12,
						Left: &rbtree.Node{
							Price: 10,
							Left:  &rbtree.Node{Price: 9},
							Right: &rbtree.Node{Price: 11},
						},
						Right: &rbtree.Node{
							Price: 14,
							Left:  &rbtree.Node{Price: 13},
							Right: &rbtree.Node{Price: 15},
						},
					},
				},
			},
			want:       15,
			wantExists: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			connectParents(tc.tree)

			if err := tc.tree.Valid(); err != nil {
				t.Fatalf("The tree must be valid: %v", err)
			}

			got := tc.tree.Max()
			if got != nil {
				if !tc.wantExists {
					t.Fatalf("Max() expected the node NOT to exist, but it does: %v", got)
				}

				if tc.want != got.Price {
					t.Errorf("Max(), want price: %v, got node: %+v", tc.want, got)
				}
			} else if tc.wantExists {
				t.Errorf("Max() expected the node to exist, but it doesn't")
			}
		})
	}
}

func Test_Get(t *testing.T) {
	testCases := []struct {
		name    string
		tree    *rbtree.Tree
		price   uint64
		wantErr bool
	}{
		{
			name:    "empty",
			tree:    &rbtree.Tree{},
			price:   1,
			wantErr: true,
		},
		{
			name:  "root_exists",
			tree:  &rbtree.Tree{Root: &rbtree.Node{Price: 41}},
			price: 41,
		},
		{
			name:    "root_unknown",
			tree:    &rbtree.Tree{Root: &rbtree.Node{Price: 41}},
			price:   2,
			wantErr: true,
		},
		{
			name: "root_with_right_exists",
			tree: &rbtree.Tree{
				Root: &rbtree.Node{
					Price: 2,
					Left:  &rbtree.Node{Price: 1, Red: true},
				},
			},
			price: 1,
		},
		{
			name: "root_with_right_unknown",
			tree: &rbtree.Tree{
				Root: &rbtree.Node{
					Price: 2,
					Left:  &rbtree.Node{Price: 1, Red: true},
				},
			},
			price:   41,
			wantErr: true,
		},
		{
			name: "root_with_left_exists",
			tree: &rbtree.Tree{
				Root: &rbtree.Node{
					Price: 1,
					Right: &rbtree.Node{Price: 2, Red: true},
				},
			},
			price: 2,
		},
		{
			name: "root_with_left_unknown",
			tree: &rbtree.Tree{
				Root: &rbtree.Node{
					Price: 1,
					Right: &rbtree.Node{Price: 2, Red: true},
				},
			},
			price:   41,
			wantErr: true,
		},
		{
			name: "one_level_exists",
			tree: &rbtree.Tree{
				Root: &rbtree.Node{
					Price: 4,
					Left: &rbtree.Node{
						Price: 2,
						Left:  &rbtree.Node{Price: 1},
						Right: &rbtree.Node{Price: 3},
					},
					Right: &rbtree.Node{
						Price: 6,
						Left:  &rbtree.Node{Price: 5},
						Right: &rbtree.Node{Price: 7},
					},
				},
			},
			price: 6,
		},
		{
			name: "one_level_unknown",
			tree: &rbtree.Tree{
				Root: &rbtree.Node{
					Price: 4,
					Left: &rbtree.Node{
						Price: 2,
						Left:  &rbtree.Node{Price: 1},
						Right: &rbtree.Node{Price: 3},
					},
					Right: &rbtree.Node{
						Price: 6,
						Left:  &rbtree.Node{Price: 5},
						Right: &rbtree.Node{Price: 7},
					},
				},
			},
			price:   41,
			wantErr: true,
		},
		{
			name: "two_levels_exists",
			tree: &rbtree.Tree{
				Root: &rbtree.Node{
					Price: 8,
					Left: &rbtree.Node{
						Price: 4,
						Left: &rbtree.Node{
							Price: 2,
							Left:  &rbtree.Node{Price: 1},
							Right: &rbtree.Node{Price: 3},
						},
						Right: &rbtree.Node{
							Price: 6,
							Left:  &rbtree.Node{Price: 5},
							Right: &rbtree.Node{Price: 7},
						},
					},
					Right: &rbtree.Node{
						Price: 12,
						Left: &rbtree.Node{
							Price: 10,
							Left:  &rbtree.Node{Price: 9},
							Right: &rbtree.Node{Price: 11},
						},
						Right: &rbtree.Node{
							Price: 14,
							Left:  &rbtree.Node{Price: 13},
							Right: &rbtree.Node{Price: 15},
						},
					},
				},
			},
			price: 10,
		},
		{
			name: "two_levels_unknown",
			tree: &rbtree.Tree{
				Root: &rbtree.Node{
					Price: 8,
					Left: &rbtree.Node{
						Price: 4,
						Left: &rbtree.Node{
							Price: 2,
							Left:  &rbtree.Node{Price: 1},
							Right: &rbtree.Node{Price: 3},
						},
						Right: &rbtree.Node{
							Price: 6,
							Left:  &rbtree.Node{Price: 5},
							Right: &rbtree.Node{Price: 7},
						},
					},
					Right: &rbtree.Node{
						Price: 12,
						Left: &rbtree.Node{
							Price: 10,
							Left:  &rbtree.Node{Price: 9},
							Right: &rbtree.Node{Price: 11},
						},
						Right: &rbtree.Node{
							Price: 14,
							Left:  &rbtree.Node{Price: 13},
							Right: &rbtree.Node{Price: 15},
						},
					},
				},
			},
			price:   41,
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			connectParents(tc.tree)

			if err := tc.tree.Valid(); err != nil {
				t.Fatalf("The tree must be valid: %v", err)
			}

			got, err := tc.tree.Get(tc.price)
			if err == nil {
				if tc.wantErr {
					t.Fatalf("Get(%v) expected error, got OK: %v", tc.price, got)
				}

				if got.Price != tc.price {
					t.Fatalf("Get(%v) got %v: %+v", tc.price, got.Price, got)
				}
			} else {
				if !tc.wantErr {
					t.Fatalf("Get(%v) expected no error, got: %v", tc.price, err)
				}
			}
		})
	}
}
