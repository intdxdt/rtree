package rtree

import (
	"github.com/intdxdt/mbr"
)

//Node type for internal node
type Node struct {
	item     BoxObj
	bbox     *mbr.MBR
	height   int
	leaf     bool
	children []*Node
}

//NewNode creates a new node
func NewNode(item BoxObj, height int, leaf bool, children []*Node) *Node {
	return &Node{
		item:     item,
		bbox:     item.BBox().Clone(),
		height:   height,
		leaf:     leaf,
		children: children,
	}
}

//Node type for internal node
func newLeafNode(item BoxObj) *Node {
	return &Node{
		item:     item,
		bbox:     item.BBox(),
		height:   1,
		leaf:     true,
		children: []*Node{},
	}
}

//BBox returns bbox property
func (n *Node) BBox() *mbr.MBR {
	return n.bbox
}

//add child
func (n *Node) addChild(child *Node) {
	n.children = append(n.children, child)
}

//GetItem from node
func (n *Node) GetItem() BoxObj {
	return n.item
}

//Constructs children of node
func makeChildren(items []BoxObj) []*Node {
	var chs = make([]*Node, len(items))
	for i := range items {
		chs[i] = newLeafNode(items[i])
	}
	return chs
}
