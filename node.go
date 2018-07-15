package rtree

import (
	"github.com/intdxdt/mbr"
)

//rNode type for internal rNode
type rNode struct {
	children []*rNode
	item     *Obj
	height   int
	leaf     bool
	bbox     *mbr.MBR
}

//newNode creates a new rNode
func newNode(item *Obj, height int, leaf bool, children []*rNode) *rNode {
	return &rNode{
		children: children,
		item:     item,
		height:   height,
		leaf:     leaf,
		bbox:     item.MBR,
	}
}

//rNode type for internal rNode
func newLeafNode(item *Obj) *rNode {
	return &rNode{
		children: []*rNode{},
		item:     item,
		height:   1,
		leaf:     true,
		bbox:     item.MBR,
	}
}

//MBR returns bbox property
func (n *rNode) BBox() *mbr.MBR {
	return n.bbox
}

//GetItem from rNode
func (n *rNode) GetItem() *Obj {
	return n.item
}

//add child
func (n *rNode) addChild(child *rNode) {
	n.children = append(n.children, child)
}



//Constructs children of rNode
func makeChildren(items []*Obj) []*rNode {
	var chs = make([]*rNode, len(items))
	for i := range items {
		chs[i] = newLeafNode(items[i])
	}
	return chs
}
