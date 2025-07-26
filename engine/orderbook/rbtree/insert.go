package rbtree

// Insert may or may not actually insert the requested price. If the price
// already exists, its node is returned, if it doesn't, the new node is
// inserted, rules are fixed, and the node containing the price is returned.
//
// Do not search if the price exists before calling Insert. This function will
// have to find its place in the tree anyways, so call this function directly
// instead for efficient programming.
//
// O(log n)
func (t *Tree) Insert(price uint64) *Node {
	holder, side, parent := t.find(price)
	if holder != nil {
		return holder
	}

	holder = t.getNode()
	holder.Price = price

	if t.orientation == MinFirst {
		if t.head == nil || price < t.head.Price {
			t.head = holder
		}
	} else {
		if t.head == nil || price > t.head.Price {
			t.head = holder
		}
	}

	if parent == nil {
		t.Root = holder
		return t.Root
	}

	holder.Red = true
	if side == LeftChild {
		parent.Left = holder
	} else {
		parent.Right = holder
	}
	holder.Parent = parent

	if parent.IsBlack() {
		return holder
	}

	var y *Node
	z := holder
	for {
		if z.Parent.IsBlack() {
			break
		}

		if z.Parent == z.Parent.Parent.Left {
			y = z.Parent.Parent.Right
			if y.IsRed() {
				z.Parent.Red = false
				y.Red = false
				z.Parent.Parent.Red = true
				z = z.Parent.Parent
			} else {
				if z == z.Parent.Right {
					z = z.Parent
					t.RotateLeft(z)
				}

				z.Parent.Red = false
				z.Parent.Parent.Red = true
				t.RotateRight(z.Parent.Parent)
			}
		} else {
			y = z.Parent.Parent.Left
			if y.IsRed() {
				z.Parent.Red = false
				y.Red = false
				z.Parent.Parent.Red = true
				z = z.Parent.Parent
			} else {
				if z == z.Parent.Left {
					z = z.Parent
					t.RotateRight(z)
				}

				z.Parent.Red = false
				z.Parent.Parent.Red = true
				t.RotateLeft(z.Parent.Parent)
			}
		}
	}
	t.Root.Red = false
	return holder
}
