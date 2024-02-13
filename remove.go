package rtree

import (
	"github.com/intdxdt/math"
	"github.com/intdxdt/mbr"
)

// Node at index from slice of pointers , nil if not found
func nodeAtIndex(a []*node, i int) *node {
	var n = len(a)
	if i > n-1 || i < 0 || n == 0 {
		return nil
	}
	return a[i]
}

// Node at index from slice of nodes , nil if not found
func nodeAtIndex_(a []node, i int) *node {
	var n = len(a)
	if i > n-1 || i < 0 || n == 0 {
		return nil
	}
	return &a[i]
}

// Pop node
func popNode(a []*node) (*node, []*node) {
	var v *node
	var n int
	if len(a) == 0 {
		return nil, a
	}
	n = len(a) - 1
	v, a[n] = a[n], nil
	return v, a[:n]
}

// Remove node at given index from node slice.
func removeNode(a []node, i int) []node {
	var n = len(a) - 1
	if i > n {
		return a
	}
	//a, a[n] = append(a[:i], a[i+1:]...), nil //panics in gopherjs
	var b = make([]node, 0, math.Max(n-i, 0))
	if i < n {
		b = append(b, a[i+1:]...)
	}
	a[n] = node{}
	a = append(a[:i], b...)
	return a
}

func popIndex(indxs *[]int) int {
	var n, v int
	var a = *indxs
	n = len(a) - 1
	if n < 0 {
		return 0
	}
	v = a[n]
	a[n] = 0
	*indxs = a[:n]
	return v
}

// condense node and its path from the root
func (tree *RTree) condense(path []*node) {
	var parent *node
	var i = len(path) - 1
	// go through the path, removing empty nodes and updating bboxes
	for i >= 0 {
		if len(path[i].children) == 0 {
			//go through the path, removing empty nodes and updating bboxes
			if i > 0 {
				parent = path[i-1]
				var index = sliceIndex(len(parent.children), func(s int) bool {
					return path[i] == &parent.children[s]
				})
				if index != -1 {
					parent.children = removeNode(parent.children, index)
				}
			} else {
				tree.Clear() //root is empty, rest root
			}
		} else {
			calcBBox(path[i])
		}
		i--
	}
}

// Remove Item from RTree
func (tree *RTree) Remove(item BoxObject) *RTree {
	if item == nil {
		return tree
	}
	tree.removeItem(item.BBox(),
		func(nd *node, i int) bool {
			return nd.children[i].item == item
		})
	return tree
}

//Remove Item from RTree
//NOTE:if item is a bbox , then first found bbox is removed
//func (tree *RTree) Remove(item *mbr.MBR) *RTree {
//	tree.removeItem(item,
//		func(nd *node, i int) bool {
//			return nd.children[i].bbox.Equals(item)
//		})
//	return tree
//}

func (tree *RTree) removeItem(item *mbr.MBR[float64], predicate func(*node, int) bool) *RTree {
	var nd = &tree.Data
	var parent *node
	var bbox = item.BBox()
	var path []*node
	var indexes []int
	var i, index int
	var goingUp bool

	// depth-first iterative self traversal
	for (nd != nil) || (len(path) > 0) {
		if nd == nil {
			// go up
			nd, path = popNode(path)
			parent = nodeAtIndex(path, len(path)-1)
			i = popIndex(&indexes)
			goingUp = true
		}

		if nd.leaf {
			// check current node
			//index = node.children.indexOf(item)
			index = sliceIndex(len(nd.children), func(i int) bool {
				return predicate(nd, i)
			})

			//if found
			if index != -1 {
				//item found, remove the item and condense self upwards
				//node.children.splice(index, 1)
				nd.children = removeNode(nd.children, index)
				path = append(path, nd)
				tree.condense(path)
				return tree
			}
		}

		if !goingUp && !nd.leaf && contains(&nd.bbox, bbox) {
			// go down
			path = append(path, nd)
			indexes = append(indexes, i)
			i = 0
			parent = nd
			nd = &nd.children[0]
		} else if parent != nil {
			// go right
			i++
			nd = nodeAtIndex_(parent.children, i)
			goingUp = false
		} else {
			nd = nil
		} // nothing found
	}
	return tree
}
