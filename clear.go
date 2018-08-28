package rtree

import (
	"math"
	"github.com/intdxdt/mbr"
)

//universe type with bounds [+inf +inf -inf -inf]
type universe struct{ bounds mbr.MBR }

func (u universe) BBox() *mbr.MBR {
	return &u.bounds
}

func createUniverse() universe{
	return universe{emptyMBR()}
}

func emptyMBR() mbr.MBR {
	return mbr.MBR{
		math.Inf(1), math.Inf(1),
		math.Inf(-1), math.Inf(-1),
	}
}

func (tree *RTree) Clear() *RTree {
	tree.Data = createNode(createUniverse(), 1, true, []node{})
	return tree
}

//IsEmpty checks for empty tree
func (tree *RTree) IsEmpty() bool {
	return len(tree.Data.children) == 0
}
