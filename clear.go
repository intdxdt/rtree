package rtree

import (
	"github.com/intdxdt/mbr"
	"math"
)

func emptyMBR() mbr.MBR {
	return mbr.MBR{
		math.Inf(1), math.Inf(1),
		math.Inf(-1), math.Inf(-1),
	}
}

func (tree *RTree) Clear() *RTree {
	tree.Data = createNode(nil, 1, true, []node{})
	return tree
}

//IsEmpty checks for empty tree
func (tree *RTree) IsEmpty() bool {
	return len(tree.Data.children) == 0
}
