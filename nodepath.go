package rtree

//nodePath slice of node
type nodePath []node

//Len for sort interface
func (path nodePath) Len() int {
	return len(path)
}

//Swap for sort interface
func (path nodePath) Swap(i, j int) {
	path[i], path[j] = path[j], path[i]
}

//xNodePath is  type  for  x sorting of boxes
type xNodePath struct {
	nodePath
}

//Less sorts boxes by ll[x]
func (path xNodePath) Less(i, j int) bool {
	return path.nodePath[i].bbox.MinX < path.nodePath[j].bbox.MinX
}

//yNodePath is type  for  y sorting of boxes
type yNodePath struct {
	nodePath
}

//Less sorts boxes by ll[y]
func (path yNodePath) Less(i, j int) bool {
	return path.nodePath[i].bbox.MinY < path.nodePath[j].bbox.MinY
}

