package rbtree

// Transplant disconnects the current node and connects the new node
// in its place. The left and right childs of the new node are not
// updated here, that is the responsibility of the calling function.
func (t *Tree) Transplant(current *Node, new *Node) {
	if current.Parent == nil {
		t.Root = new
	} else if current == current.Parent.Left {
		current.Parent.Left = new
	} else {
		current.Parent.Right = new
	}

	if new != nil {
		new.Parent = current.Parent
	}
}
