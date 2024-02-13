package rtree

import (
	"github.com/intdxdt/mbr"
)

// Node type for internal rtree node
type node struct {
	children []node
	item     BoxObject
	height   int
	leaf     bool
	bbox     mbr.MBR[float64]
}

// Creates a node
func createNode(item BoxObject, height int, leaf bool, children []node) node {
	var box mbr.MBR[float64]
	if item == nil {
		box = emptyMBR()
	} else {
		box = item.BBox().Clone()
	}
	return node{
		children: children,
		item:     item,
		height:   height,
		leaf:     leaf,
		bbox:     box,
	}
}

// node type for internal node
func newLeafNode(item BoxObject) node {
	return node{
		children: []node{},
		item:     item,
		height:   1,
		leaf:     true,
		bbox:     item.BBox().Clone(),
	}
}

// MBR returns bbox property
func (nd *node) BBox() *mbr.MBR[float64] {
	return &nd.bbox
}

// GetItem from node
func (nd *node) GetItem() BoxObject {
	return nd.item
}

// add child
func (nd *node) addChild(child node) {
	nd.children = append(nd.children, child)
}

// Constructs children of node
func makeChildren(items []BoxObject) []node {
	var chs = make([]node, 0, len(items))
	for i := range items {
		chs = append(chs, newLeafNode(items[i]))
	}
	return chs
}
