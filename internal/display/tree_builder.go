package display

import (
	"github.com/Equationzhao/g/internal/display/tree"
	"github.com/Equationzhao/g/internal/item"
)

type TreeBuilder struct{}

func NewTreeBuilder() *TreeBuilder {
	return &TreeBuilder{}
}

func (b *TreeBuilder) Build(items []*item.FileInfo) *tree.Tree {
	if len(items) == 0 {
		return tree.NewTree()
	}
	buildTree := tree.NewTree(tree.WithCap(len(items) / 2))
	level := make(map[int][]*item.FileInfo)
	for _, v := range items {
		level[v.Level] = append(level[v.Level], v)
	}

	nodeMap := make(map[string]*tree.Node, len(level))
	root := level[0][0]
	buildTree.Root.Meta = root
	nodeMap[root.FullPath] = buildTree.Root

	for i := 1; i < len(level); i++ {
		for _, v := range level[i] {
			node := nodeMap[v.ParentPath]
			c := &tree.Node{
				Parent: node,
				Child:  make([]*tree.Node, 0, 10),
				Level:  i,
				Meta:   v,
			}
			nodeMap[v.FullPath] = c
			node.AddChild(c)
		}
	}

	return buildTree
}
