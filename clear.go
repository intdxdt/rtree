package rtree

import (
	"math"
	"github.com/intdxdt/mbr"
)

//universe type with bounds [+inf +inf -inf -inf]
type universe struct{}

func (u *universe) BBox() *mbr.MBR {
	return emptyMbr()
}

func emptyMbr() *mbr.MBR {
	return &mbr.MBR{
		math.Inf(1), math.Inf(1),
		math.Inf(-1), math.Inf(-1),
	}
}


func (tree *RTree) Clear() *RTree {
	node := NewNode(&universe{}, 1, true, make([]*Node, 0))
	tree.Data = node
	return tree
}

//IsEmpty checks for empty tree
func (tree *RTree) IsEmpty() bool {
	return len(tree.Data.children) == 0
}
