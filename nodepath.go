package rtree

//NodePath slice of node
type NodePath []node

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

