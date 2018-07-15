package rtree

import "github.com/intdxdt/mbr"

func nodeAtIndex(a []*rNode, i int) *rNode {
	var n = len(a)
	if i > n-1 || i < 0 || n == 0 {
		return nil
	}
	return a[i]
}

func popNode(a []*rNode) (*rNode, []*rNode) {
	var v *rNode
	var n int
	if len(a) == 0 {
		return nil, a
	}
	n = len(a) - 1
	v, a[n] = a[n], nil
	return v, a[:n]
}

func popIndex(indxs *[]int) int {
	var n, v int
	a := *indxs
	n = len(a) - 1
	if n < 0 {
		return 0
	}
	v, a[n] = a[n], 0
	*indxs = a[:n]
	return v
}

//remove rNode at given index from rNode slice.
func removeNode(a []*rNode, i int) []*rNode {
	var n = len(a) - 1
	if i > n {
		return a
	}
	a, a[n] = append(a[:i], a[i+1:]...), nil
	return a
}

//condense rNode and its path from the root
func (tree *RTree) condense(path []*rNode) {
	var parent *rNode
	var i = len(path) - 1
	// go through the path, removing empty nodes and updating bboxes
	for i >= 0 {
		if len(path[i].children) == 0 {
			//go through the path, removing empty nodes and updating bboxes
			if i > 0 {
				parent = path[i-1]
				index := sliceIndex(len(parent.children), func(s int) bool {
					return path[i] == parent.children[s]
				})
				if index != -1 {
					parent.children = removeNode(parent.children, index)
				}
			} else {
				//root is empty, rest root
				tree.Clear()
			}
		} else {
			calcBBox(path[i])
		}
		i--
	}
}

//Remove Item from RTree
//NOTE:if item is a bbox , then first found bbox is removed
func (tree *RTree) RemoveObj(item *Obj) *RTree {
	if (item == nil) {
		return tree
	}
	tree.removeItem(item.MBR,
		func(node *rNode, i int) bool {
			return node.children[i].item == item
		})
	return tree
}

//Remove Item from RTree
//NOTE:if item is a bbox , then first found bbox is removed
func (tree *RTree) RemoveMBR(item *mbr.MBR) *RTree {
	tree.removeItem(item,
		func(node *rNode, i int) bool {
			return node.children[i].bbox.Equals(item)
		})
	return tree
}

func (tree *RTree) removeItem(item *mbr.MBR, predicate func(*rNode, int) bool) *RTree {
	if item == nil {
		return tree
	}

	var node = tree.Data
	var parent *rNode
	var bbox = item.BBox()
	var path = make([]*rNode, 0)
	var indexes = make([]int, 0)
	var i, index int
	var goingUp bool

	// depth-first iterative self traversal
	for (node != nil) || (len(path) > 0) {
		if node == nil {
			// go up
			node, path = popNode(path)
			parent = nodeAtIndex(path, len(path)-1)
			i = popIndex(&indexes)
			goingUp = true
		}

		if node.leaf {
			// check current rNode
			//index = rNode.children.indexOf(item)
			index = sliceIndex(len(node.children), func(i int) bool {
				return predicate(node, i)
			})

			//if found
			if index != -1 {
				//item found, remove the item and condense self upwards
				//rNode.children.splice(index, 1)
				node.children = removeNode(node.children, index)
				path = append(path, node)
				tree.condense(path)
				return tree
			}
		}

		if !goingUp && !node.leaf && contains(node.bbox, bbox) {
			// go down
			path = append(path, node)
			indexes = append(indexes, i)
			i = 0
			parent = node
			node = node.children[0]
		} else if parent != nil {
			// go right
			i++
			node = nodeAtIndex(parent.children, i)
			goingUp = false
		} else {
			node = nil
		} // nothing found
	}
	return tree
}
