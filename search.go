package rtree

import (
	"github.com/intdxdt/mbr"
)

//Search item
func (tree *RTree) Search(bbox *mbr.MBR) []*Obj {

	var result []*rNode
	var node = tree.Data

	if !intersects(bbox, node.bbox) {
		return []*Obj{}
	}

	var nodesToSearch []*rNode
	var child *rNode
	var childBBox *mbr.MBR

	for {
		for i, length := 0, len(node.children); i < length; i++ {
			child = node.children[i]
			childBBox = child.bbox

			if intersects(bbox, childBBox) {
				if node.leaf {
					result = append(result, child)
				} else if contains(bbox, childBBox) {
					result = all(child, result)
				} else {
					nodesToSearch = append(nodesToSearch, child)
				}
			}
		}

		node, nodesToSearch = popNode(nodesToSearch)
		if node == nil {
			break
		}
	}

	var objs = make([]*Obj, 0, len(result))
	for i := range result {
		objs = append(objs, result[i].item)
	}
	return objs
}

//All items from  root rNode
func (tree *RTree) All() []*rNode {
	return all(tree.Data, make([]*rNode, 0))
}

//all - fetch all items from rNode
func all(node *rNode, result []*rNode) []*rNode {
	var nodesToSearch = make([]*rNode, 0)
	for {
		if node.leaf {
			for i := range node.children {
				result = append(result, node.children[i])
			}
		} else {
			for i := range node.children {
				nodesToSearch = append(nodesToSearch, node.children[i])
			}
		}

		node, nodesToSearch = popNode(nodesToSearch)
		if node == nil {
			break
		}
	}

	return result
}
