package rbtree

import (
	"fmt"
)

// Valid returns an error if the tree is not valid, giving a detailed
// explanation of the reason why it is invalid. This function should only be used
// in tests, or to debug errors if they arise. The Red-Black tree implementation
// should never reach an invalid state.
func (t *Tree) Valid() error {
	if t == nil {
		return fmt.Errorf("invalid tree: tree is a null pointer")
	}

	if t.Root == nil {
		return nil
	}

	if t.Root.Red {
		return fmt.Errorf("invalid tree: root must be black")
	}

	visited := map[*Node]bool{}
	stack := []*Node{t.Root}
	countStack := []int{0}
	blackPathLength := -1

	for {
		if len(stack) == 0 {
			break
		}

		i := len(stack) - 1
		node := stack[i]
		stack = stack[:i]

		blackCount := countStack[i]
		countStack = countStack[:i]

		if _, ok := visited[node]; !ok {
			visited[node] = true

			if node.Red && node.Parent.Red {
				side := "L"
				if node.Parent.Right == node {
					side = "R"
				}
				return fmt.Errorf("invalid tree: two directly related red nodes: %d -%s-> %d", node.Parent.Price, side, node.Price)
			}

			if node.Parent != nil {
				if node == node.Parent.Left && node.Price >= node.Parent.Price {
					return fmt.Errorf("invalid tree: left child is not less than: %d -> %d", node.Parent.Price, node.Price)
				}

				if node == node.Parent.Right && node.Price <= node.Parent.Price {
					return fmt.Errorf("invalid tree: right child is not bigger than: %d -> %d", node.Parent.Price, node.Price)
				}
			}

			if node.Right != nil { // right first, so that it's processed last
				if node.Right.Parent != node {
					return fmt.Errorf("invalid tree: right child not referencing parent: %d <- %d", node.Price, node.Right.Price)
				}

				stack = append(stack, node.Right)

				if node.Right.Red {
					countStack = append(countStack, blackCount)
				} else {
					countStack = append(countStack, blackCount+1)
				}
			}

			if node.Left != nil {
				if node.Left.Parent != node {
					return fmt.Errorf("invalid tree: left child not referencing parent: %d <- %d", node.Price, node.Right.Price)
				}

				stack = append(stack, node.Left)

				if node.Left.Red {
					countStack = append(countStack, blackCount)
				} else {
					countStack = append(countStack, blackCount+1)
				}
			}

			if node.Left == nil || node.Right == nil {
				if blackPathLength == -1 {
					blackPathLength = blackCount
				} else if blackCount != blackPathLength {
					return fmt.Errorf("invalid tree: unequal black counts, want %d, got %d, from root to %d", blackCount, blackPathLength, node.Price)
				}
			}
		}
	}

	return nil
}
