package tree

import "github.com/Equationzhao/g/internal/item"

/*
build tree like this:
drwxr-xr-x@    - mr.black 10 7 03:38 ├── filter
drwxr-xr-x@    - mr.black 10 7 03:38 │  ├── content
.rw-r--r--@ 1.7k mr.black 10 7 03:38 │  │  ├── charset.go
.rw-r--r--@ 4.3k mr.black  7 7 23:39 │  │  ├── duplicate.go
.rw-r--r--@ 2.5k mr.black  7 7 20:35 │  │  ├── git.go
.rw-r--r--@  949 mr.black  5 7 01:26 │  │  ├── group.go
.rw-r--r--@  323 mr.black  5 7 01:26 │  │  ├── index.go
.rw-r--r--@  479 mr.black 10 7 03:38 │  │  ├── inode.go
.rw-r--r--@  564 mr.black  5 7 01:26 │  │  ├── link.go
.rw-r--r--@ 1.5k mr.black 10 7 03:38 │  │  ├── mimetype.go
.rw-r--r--@ 5.8k mr.black  7 7 20:35 │  │  ├── name.go
.rw-r--r--@  972 mr.black  5 7 01:26 │  │  ├── owner.go
.rw-r--r--@  743 mr.black  5 7 01:26 │  │  ├── permission.go
.rw-r--r--@ 5.5k mr.black 10 7 03:38 │  │  ├── size.go
.rw-r--r--@ 2.7k mr.black  7 7 20:35 │  │  ├── sum.go
.rw-r--r--@ 1.7k mr.black  7 7 20:35 │  │  └── time.go
.rw-r--r--@ 2.8k mr.black 10 7 03:38 │  ├── contentfilter.go
.rw-r--r--@ 5.0k mr.black 10 7 03:38 │  └── itemfliter.go
.rw-r--r--@ 9.6k mr.black  7 7 23:39 ├── g.md
...
*/

type Node struct {
	Parent *Node
	Child  []*Node
	Level  int
	Meta   *item.FileInfo
}

func (n *Node) Apply2Child(f func(node *Node)) {
	if n.Child == nil {
		return
	}
	for _, child := range n.Child {
		f(child)
		child.Apply2Child(f)
	}
}

func (n *Node) AddChild(child *Node) *Node {
	n.Child = append(n.Child, child)
	child.Parent = n
	child.Level = n.Level + 1
	return n
}

func (n *Node) Apply2ChildSlice(connectors func(nodes []*Node)) {
	if n.Child == nil {
		return
	}
	connectors(n.Child)
	for _, child := range n.Child {
		child.Apply2ChildSlice(connectors)
	}
}

func (n *Node) ApplyThis(p func(node *Node)) {
	p(n)
}

type Tree struct {
	Root *Node
}

type Option = func(tree *Tree)

func WithCap(cap int) Option {
	return func(tree *Tree) {
		tree.Root.Child = make([]*Node, 0, cap)
	}
}

func NewTree(ops ...Option) *Tree {
	t := &Tree{
		Root: &Node{
			Parent: nil,
			Level:  0,
		},
	}
	for _, op := range ops {
		op(t)
	}
	if t.Root.Child == nil {
		t.Root.Child = make([]*Node, 0, 10)
	}
	return t
}
