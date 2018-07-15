package rtree

import (
	"github.com/intdxdt/mbr"
)

//Node type for internal Node
type Node struct {
	children []*Node
	item     *Obj
	height   int
	leaf     bool
	bbox     *mbr.MBR
}

//NewNode creates a new Node
func NewNode(item *Obj, height int, leaf bool, children []*Node) *Node {
	return &Node{
		children: children,
		item:     item,
		height:   height,
		leaf:     leaf,
		bbox:     item.MBR,
	}
}

//Node type for internal Node
func newLeafNode(item *Obj) *Node {
	return &Node{
		children: []*Node{},
		item:     item,
		height:   1,
		leaf:     true,
		bbox:     item.MBR,
	}
}

//MBR returns bbox property
func (n *Node) BBox() *mbr.MBR {
	return n.bbox
}

//add child
func (n *Node) addChild(child *Node) {
	n.children = append(n.children, child)
}

//GetItem from Node
func (n *Node) GetItem() *Obj {
	return n.item
}

//Constructs children of Node
func makeChildren(items []*Obj) []*Node {
	var chs = make([]*Node, len(items))
	for i := range items {
		chs[i] = newLeafNode(items[i])
	}
	return chs
}
