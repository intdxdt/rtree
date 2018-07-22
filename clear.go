package rtree

import (
	"math"
	"github.com/intdxdt/mbr"
)

func emptyMBR() *mbr.MBR {
	return &mbr.MBR{
		math.Inf(1), math.Inf(1),
		math.Inf(-1), math.Inf(-1),
	}
}

type Obj struct {
	Id     int
	Meta   int
	MBR    *mbr.MBR
	Object interface{}
}

func emptyObject() *Obj {
	return &Obj{
		Id:     -1,
		MBR:    emptyMBR(),
		Object: nil,
		Meta:   -1,
	}
}

func Object(id int, box *mbr.MBR, object ...interface{}) *Obj {
	var obj interface{}
	if len(object) > 0 {
		obj = object[0]
	}
	return &Obj{
		Id:     id,
		MBR:    box,
		Object: obj,
		Meta:   -1,
	}
}

func (tree *RTree) Clear() *RTree {
	tree.Data = newNode(
		emptyObject(),
		1, true, []node{},
	)
	return tree
}

//IsEmpty checks for empty tree
func (tree *RTree) IsEmpty() bool {
	return len(tree.Data.children) == 0
}
