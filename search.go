package rtree

import (
	"github.com/intdxdt/mbr"
)

//Search item
func (tree *RTree) Search(bbox *mbr.MBR) []*Node {

	var result []*Node
	var node = tree.Data

	if !intersects(bbox, node.bbox) {
		return result
	}

	var nodesToSearch []*Node
	var child *Node
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

	return result
}

//All items from  root node
func (tree *RTree) All() []*Node {
	return all(tree.Data, make([]*Node, 0))
}

//all - fetch all items from node
func all(node *Node, result []*Node) []*Node {
	var nodesToSearch = make([]*Node, 0)
	for true {
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
