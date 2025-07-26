package rbtree

// https://stackoverflow.com/questions/4597650/code-with-explanation-for-binary-tree-rotation-Left-or-Right

// RotateLeft is an elementary operation of binary trees that keeps
// the binary order valid.
//
// It takes the node x, and makes its right child y the
// new root of the subtree. x becomes the left child of y and y's left subtree
// becomes the right subtree of x.
func (t *Tree) RotateLeft(x *Node) {
	if x == nil {
		return
	}

	if x.Right == nil {
		return
	}

	y := x.Right
	x.Right = y.Left
	if y.Left != nil { // otherwise .Parent would fail
		y.Left.Parent = x // set subtree on the new node
	}

	y.Parent = x.Parent
	if x.Parent == nil { // or, if t.root == x
		t.Root = y
	} else if x == x.Parent.Left { // if x is the Left child
		x.Parent.Left = y // the new Left is y
	} else { // x is the Right child
		x.Parent.Right = y // the new Right is y
	}

	y.Left = x
	x.Parent = y
}

// RotateRight is an elementary operation of binary trees that keeps
// the binary ordering valid.
//
// It takes the node x, and makes its left child y the
// new root of the subtree. x becomes the right child of y and y's right subtree
// becomes the left subtree of x.
func (t *Tree) RotateRight(x *Node) {
	if x == nil || x.Left == nil {
		return
	}

	y := x.Left
	x.Left = y.Right
	if y.Right != nil {
		y.Right.Parent = x
	}

	y.Parent = x.Parent
	if x.Parent == nil {
		t.Root = y
	} else if x == x.Parent.Left {
		x.Parent.Left = y
	} else {
		x.Parent.Right = y
	}

	y.Right = x
	x.Parent = y
}
