package rtree

import (
	"github.com/intdxdt/math"
)

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
type XNodePath struct {
	NodePath
}

//Less sorts boxes by ll[x]
func (path XNodePath) Less(i, j int) bool {
	return path.NodePath[i].bbox[0] < path.NodePath[j].bbox[0]
}

//YNodePath is type  for  y sorting of boxes
type YNodePath struct {
	NodePath
}

//Less sorts boxes by ll[y]
func (path YNodePath) Less(i, j int) bool {
	return path.NodePath[i].BBox()[1] < path.NodePath[j].BBox()[1]
}

//XYNodePath is type  for  xy sorting of boxes
type XYNodePath struct {
	NodePath
}

//Less sorts boxes lexicographically
func (path XYNodePath) Less(i, j int) bool {
	var x, y = 0, 1
	a, b := path.NodePath[i].BBox(), path.NodePath[j].BBox()
	d := a[x] - b[x]
	//x's are close enough to each other
	if math.FloatEqual(d, 0.0) {
		d = a[y] - b[y]
	}
	//check if close enough ot zero
	if d < 0 {
		return true
	}
	return false
}
