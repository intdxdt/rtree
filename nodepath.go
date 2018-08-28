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

//xNodePath is  type  for  x sorting of boxes
type xNodePath struct {
	NodePath
}

//Less sorts boxes by ll[x]
func (path xNodePath) Less(i, j int) bool {
	return path.NodePath[i].bbox.MinX < path.NodePath[j].bbox.MinX
}

//yNodePath is type  for  y sorting of boxes
type yNodePath struct {
	NodePath
}

//Less sorts boxes by ll[y]
func (path yNodePath) Less(i, j int) bool {
	return path.NodePath[i].bbox.MinY < path.NodePath[j].bbox.MinY
}

