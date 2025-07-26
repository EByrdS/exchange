package rbtree

// Delete may or may not actually delete the requested price. If the price
// exists, the node is deleted, if it doesn't then the function returns.
//
// Do not search if the price exists before calling Delete. This function will
// have to find its place in the tree anyways, so call this function directly
// instead for efficient programming.
//
// O(log n)
func (t *Tree) Delete(price uint64) {
	holder, holderSide, holderParent := t.find(price)

	t.delete(holder, holderSide, holderParent)
}

// DeleteNode skips finding the node in the tree and goes straight to deleting
// it, making it much faster than Delete.
//
// The node contents are destroyed after using this function; do not re-use it.
//
// It is responsibility of the caller to send a node that exists in this tree.
//
// O(log n)
func (t *Tree) DeleteNode(node *Node) {
	nodeSide := LeftChild
	if node.Parent != nil && node.Parent.Right == node {
		nodeSide = RightChild
	}

	t.delete(node, nodeSide, node.Parent)
}

// O(log n)
func (t *Tree) delete(holder *Node, holderSide ChildSide, holderParent *Node) {
	if holder == nil {
		return
	}

	if t.head == holder {
		if t.orientation == MinFirst {
			if holder.Right == nil {
				t.head = holder.Parent
			} else {
				t.head = t.min(holder.Right)
			}
		} else {
			if holder.Left == nil {
				t.head = holder.Parent
			} else {
				t.head = t.max(holder.Left)
			}
		}
	}

	originalRed := holder.Red

	var x, xParent *Node
	var xSide ChildSide
	if holder.Left == nil {
		// case 1: left is nil
		x = holder.Right // x can be nil

		// because we will transplant x in place of the holder
		xParent = holderParent
		xSide = holderSide

		t.Transplant(holder, x)
	} else if holder.Right == nil {
		// case 2: right is nil
		x = holder.Left // x can't be nil

		// because we will transplant x in place of the holder
		xParent = holderParent
		xSide = holderSide

		t.Transplant(holder, x)
	} else {
		// case 3: neither child is nil
		y := t.min(holder.Right)
		originalRed = y.Red

		// x can be nil but not always
		// y is the leftmost node, but it can have a right child
		x = y.Right
		xParent = y
		xSide = RightChild

		if y.Parent == holder {
			if x != nil { // x is nil if y.Right == nil
				x.Parent = y
			}
		} else {

			// Because we will transplant x in place of y, copy its references
			xParent = y.Parent // y.Parent is holder.Right or something deeper
			xSide = LeftChild  // y is always the left child

			t.Transplant(y, x)
			y.Right = holder.Right
			y.Right.Parent = y
		}

		t.Transplant(holder, y)
		y.Left = holder.Left
		y.Left.Parent = y
		y.Red = holder.Red
	}

	if !originalRed {
		t.deleteFixup(x, xParent, xSide)
	}

	t.putNode(holder)
}

func (t *Tree) safeChild(parent *Node, side ChildSide) *Node {
	if parent == nil {
		return nil
	}
	if side == LeftChild {
		return parent.Left
	}
	return parent.Right
}

func (t *Tree) deleteFixup(x *Node, xParent *Node, xSide ChildSide) {
	if x == nil && xParent == nil {
		return
	}

	// w is x sibling
	// Four cases:
	// case 1: when w is red
	// case 2: when w is black and w.left and w.right are black
	// case 3: when w is black and w.right is red and w.left is black
	// case 4: when w is black and w.right is red

	var w *Node
	for t.Root != x && x.IsBlack() {
		if xSide == LeftChild {
			w = t.safeChild(xParent, RightChild)

			// case 1
			if w.IsRed() {
				w.Red = false
				w.Parent.Red = true // or xParent.Red = true
				t.RotateLeft(xParent)
				w = xParent.Right
			}

			// case 2
			if w != nil && w.Left.IsBlack() && w.Right.IsBlack() {
				w.Red = true
				x = xParent

				xParent = x.Parent
				xSide = LeftChild
				if xParent != nil && x == xParent.Right {
					xSide = RightChild
				}

			} else {
				// case 3
				if w != nil && w.Right.IsBlack() {
					w.Left.Red = false
					w.Red = true
					t.RotateRight(w)
					w = xParent.Right
				}

				// case 4
				if w != nil {
					w.Red = xParent.Red
					if w.Right != nil {
						w.Right.Red = false
					}
				}
				xParent.Red = false
				t.RotateLeft(xParent)
				x = t.Root
			}
		} else {
			w = t.safeChild(xParent, LeftChild)
			// case 1
			if w.IsRed() {
				w.Red = false
				w.Parent.Red = true // or xParent.Red = true
				t.RotateRight(xParent)
				w = xParent.Left
			}

			// case 2
			if w != nil && w.Right.IsBlack() && w.Left.IsBlack() {
				w.Red = true
				x = xParent
				xParent = x.Parent
				xSide = LeftChild
				if xParent != nil && x == xParent.Right {
					xSide = RightChild
				}
			} else {
				// case 3
				if w != nil && w.Left.IsBlack() {
					w.Red = true
					if w.Right != nil {
						w.Right.Red = false
					}
					t.RotateLeft(w)
					w = xParent.Left
				}

				// case 4
				if w != nil {
					w.Red = xParent.Red
					if w.Left != nil {
						w.Left.Red = false
					}
				}
				xParent.Red = false
				t.RotateRight(xParent)
				x = t.Root
			}
		}
	}

	x.Red = false
}
