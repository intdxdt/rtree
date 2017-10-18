package rtree

import (
	"github.com/intdxdt/mbr"
	"github.com/intdxdt/math"
)

//Node type for internal node
type Node struct {
	item     BoxObj
	bbox     *mbr.MBR
	height   int
	leaf     bool
	children []*Node
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

//NodePath slice of Node
type NodePath []*Node

//Len for sort interface
func (path NodePath) Len() int {
	return len(path)
}

//Swap for sort interface
func (path NodePath) Swap(i, j int) {
	path[i], path[j] = path[j], path[i]
}

//XNodePath is  type  for  x sorting of boxes
type XNodePath struct{
	NodePath
}

//Less sorts boxes by ll[x]
func (path XNodePath) Less(i, j int) bool {
	return path.NodePath[i].BBox()[0] < path.NodePath[j].BBox()[0]
}

//YNodePath is type  for  y sorting of boxes
type YNodePath struct{
	NodePath
}

//Less sorts boxes by ll[y]
func (path YNodePath) Less(i, j int) bool {
	return path.NodePath[i].BBox()[1] < path.NodePath[j].BBox()[1]
}

//XYNodePath is type  for  xy sorting of boxes
type XYNodePath struct{
	NodePath
}

//Less sorts boxes lexicographically
func (path XYNodePath) Less(i, j int) bool {
	var x, y = 0, 1
	a, b := path.NodePath[i].BBox(), path.NodePath[j].BBox()
	d := a[x] - b[x]
	//x's are close enougth to each other
	if math.FloatEqual(d, 0.0) {
		d = a[y] - b[y]
	}
	//check if close enougth ot zero
	if d < 0 {
		return true
	}
	return false
}

//NewNode creates a new node
func NewNode(item BoxObj, height int, leaf bool, children []*Node) *Node {
	bbox := item.BBox().Clone()
	node := Node{
		item:     item,
		bbox:     bbox,
		height:   height,
		leaf:     leaf,
		children: children,
	}
	return &node
}

//Node type for internal node
func newLeafNode(item BoxObj) *Node {
	bbox := item.BBox()
	chs := make([]*Node, 0)
	return &Node{
		item    :   item,
		bbox    :   bbox,
		height  :   1,
		leaf    :   true,
		children:   chs,
	}
}

//Constructs children of node
func makeChildren(items []BoxObj) []*Node {
	chs := make([]*Node, len(items))
	for i := range items {
		chs[i] = newLeafNode(items[i])
	}
	return chs
}






