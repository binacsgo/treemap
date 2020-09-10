package treemap

const (
	red   = 0 // red
	black = 1 // black
)

// Keytype the key type
type Keytype interface {
	LessThan(interface{}) bool
	Equal(interface{}) bool
}

type valuetype interface{}

type Node struct {
	left, right, parent *Node
	color               int
	Key                 Keytype
	Value               valuetype
}

// Tree tree
type Tree struct {
	root *Node
	size int
}

// NewTree return a pointer to Tree
func NewTree() *Tree {
	return &Tree{}
}

// Find ...
func (t *Tree) Find(key Keytype) interface{} {
	n := t.findnode(key)
	if n != nil {
		return n.Value
	}
	return nil
}

// FindIter ...
func (t *Tree) FindIter(key Keytype) *Node {
	return t.findnode(key)
}

// Empty ...
func (t *Tree) Empty() bool {
	if t.root == nil {
		return true
	}
	return false
}

// Iterator ...
func (t *Tree) Iterator() *Node {
	return minimum(t.root)
}

// Size ...
func (t *Tree) Size() int {
	return t.size
}

// Clear ...
func (t *Tree) Clear() {
	t.root = nil
	t.size = 0
}

// Insert ...
func (t *Tree) Insert(key Keytype, value valuetype) {
	x := t.root
	var y *Node

	for x != nil {
		y = x
		if key.LessThan(x.Key) {
			x = x.left
		} else {
			x = x.right
		}
	}

	z := &Node{parent: y, color: red, Key: key, Value: value}
	t.size++

	if y == nil {
		z.color = black
		t.root = z
		return
	} else if z.Key.LessThan(y.Key) {
		y.left = z
	} else {
		y.right = z
	}
	t.rbInsertFixup(z)
}

// Delete ...
func (t *Tree) Delete(key Keytype) {
	z := t.findnode(key)
	if z == nil {
		return
	}

	var x, y *Node
	if z.left != nil && z.right != nil {
		y = successor(z)
	} else {
		y = z
	}

	if y.left != nil {
		x = y.left
	} else {
		x = y.right
	}

	xparent := y.parent
	if x != nil {
		x.parent = xparent
	}
	if y.parent == nil {
		t.root = x
	} else if y == y.parent.left {
		y.parent.left = x
	} else {
		y.parent.right = x
	}

	if y != z {
		z.Key = y.Key
		z.Value = y.Value
	}

	if y.color == black {
		t.rbDeleteFixup(x, xparent)
	}
	t.size--
}

func (t *Tree) rbInsertFixup(z *Node) {
	var y *Node
	for z.parent != nil && z.parent.color == red {
		if z.parent == z.parent.parent.left {
			y = z.parent.parent.right
			if y != nil && y.color == red {
				z.parent.color = black
				y.color = black
				z.parent.parent.color = red
				z = z.parent.parent
			} else {
				if z == z.parent.right {
					z = z.parent
					t.leftRotate(z)
				}
				z.parent.color = black
				z.parent.parent.color = red
				t.rightRotate(z.parent.parent)
			}
		} else {
			y = z.parent.parent.left
			if y != nil && y.color == red {
				z.parent.color = black
				y.color = black
				z.parent.parent.color = red
				z = z.parent.parent
			} else {
				if z == z.parent.left {
					z = z.parent
					t.rightRotate(z)
				}
				z.parent.color = black
				z.parent.parent.color = red
				t.leftRotate(z.parent.parent)
			}
		}
	}
	t.root.color = black
}

func (t *Tree) rbDeleteFixup(x, parent *Node) {
	var w *Node
	for x != t.root && getColor(x) == black {
		if x != nil {
			parent = x.parent
		}
		if x == parent.left {
			w = parent.right
			if w.color == red {
				w.color = black
				parent.color = red
				t.leftRotate(parent)
				w = parent.right
			}
			if getColor(w.left) == black && getColor(w.right) == black {
				w.color = red
				x = parent
			} else {
				if getColor(w.right) == black {
					if w.left != nil {
						w.left.color = black
					}
					w.color = red
					t.rightRotate(w)
					w = parent.right
				}
				w.color = parent.color
				parent.color = black
				if w.right != nil {
					w.right.color = black
				}
				t.leftRotate(parent)
				x = t.root
			}
		} else {
			w = parent.left
			if w.color == red {
				w.color = black
				parent.color = red
				t.rightRotate(parent)
				w = parent.left
			}
			if getColor(w.left) == black && getColor(w.right) == black {
				w.color = red
				x = parent
			} else {
				if getColor(w.left) == black {
					if w.right != nil {
						w.right.color = black
					}
					w.color = red
					t.leftRotate(w)
					w = parent.left
				}
				w.color = parent.color
				parent.color = black
				if w.left != nil {
					w.left.color = black
				}
				t.rightRotate(parent)
				x = t.root
			}
		}
	}
	if x != nil {
		x.color = black
	}
}

func (t *Tree) leftRotate(x *Node) {
	y := x.right
	x.right = y.left
	if y.left != nil {
		y.left.parent = x
	}
	y.parent = x.parent
	if x.parent == nil {
		t.root = y
	} else if x == x.parent.left {
		x.parent.left = y
	} else {
		x.parent.right = y
	}
	y.left = x
	x.parent = y
}

func (t *Tree) rightRotate(x *Node) {
	y := x.left
	x.left = y.right
	if y.right != nil {
		y.right.parent = x
	}
	y.parent = x.parent
	if x.parent == nil {
		t.root = y
	} else if x == x.parent.right {
		x.parent.right = y
	} else {
		x.parent.left = y
	}
	y.right = x
	x.parent = y
}

// findnode finds the node by key and return it, if not exists return nil.
func (t *Tree) findnode(key Keytype) *Node {
	x := t.root
	for x != nil {
		if key.LessThan(x.Key) {
			x = x.left
		} else {
			//if key == x.Key {
			if key.Equal(x.Key) {
				return x
			}
			x = x.right
		}
	}
	return nil
}

// Next returns the node's successor as an iterator.
func (n *Node) Next() *Node {
	return successor(n)
}

// successor returns the successor of the node
func successor(x *Node) *Node {
	if x.right != nil {
		return minimum(x.right)
	}
	y := x.parent
	for y != nil && x == y.right {
		x = y
		y = x.parent
	}
	return y
}

// getColor gets color of the node.
func getColor(n *Node) int {
	if n == nil {
		return black
	}
	return n.color
}

// minimum finds the minimum node of subtree n.
func minimum(n *Node) *Node {
	for n.left != nil {
		n = n.left
	}
	return n
}
