package rtree

import (
	"github.com/intdxdt/mbr"
)

//node type for internal node
type node struct {
	children []node
	item     BoxObj
	height   int
	leaf     bool
	bbox     *mbr.MBR
}

//newNode creates a new node
func newNode(item BoxObj, height int, leaf bool, children []node) node {
	return node{
		children: children,
		item:     item,
		height:   height,
		leaf:     leaf,
		bbox:     item.BBox(),
	}
}

//node type for internal node
func newLeafNode(item BoxObj) node {
	return node{
		children: []node{},
		item:     item,
		height:   1,
		leaf:     true,
		bbox:     item.BBox(),
	}
}


//MBR returns bbox property
func (nd *node) BBox() *mbr.MBR {
	return nd.bbox
}

//GetItem from node
func (nd *node) GetItem() BoxObj {
	return nd.item
}

//add child
func (nd *node) addChild(child node) {
	nd.children = append(nd.children, child)
}

//Constructs children of node
func makeChildren(items []BoxObj) []node {
	var chs = make([]node, 0, len(items))
	for i := range items {
		chs = append(chs, newLeafNode(items[i]))
	}
	return chs
}
