package rtree

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
	return path.NodePath[i].bbox[1] < path.NodePath[j].bbox[1]
}

//XYNodePath is type  for  xy sorting of boxes
type XYNodePath struct {
	NodePath
}

//Less sorts boxes lexicographically
func (path XYNodePath) Less(i, j int) bool {
	var d = path.NodePath[i].bbox[0] - path.NodePath[j].bbox[0]
	if feq(d, 0.0) {
		d = path.NodePath[i].bbox[1] - path.NodePath[j].bbox[1]
	}
	return d < 0
}
