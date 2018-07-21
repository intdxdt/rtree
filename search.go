package rtree

import (
	"github.com/intdxdt/mbr"
)

//Search item
func (tree *RTree) Search(query mbr.MBR) []*Obj {
	var bbox = &query
	var result []*node
	var nd = tree.Data

	if !intersects(bbox, nd.bbox) {
		return []*Obj{}
	}

	var nodesToSearch []*node
	var child *node
	var childBBox *mbr.MBR

	for {
		for i, length := 0, len(nd.children); i < length; i++ {
			child = nd.children[i]
			childBBox = child.bbox

			if intersects(bbox, childBBox) {
				if nd.leaf {
					result = append(result, child)
				} else if contains(bbox, childBBox) {
					result = all(child, result)
				} else {
					nodesToSearch = append(nodesToSearch, child)
				}
			}
		}

		nd, nodesToSearch = popNode(nodesToSearch)
		if nd == nil {
			break
		}
	}

	var objs = make([]*Obj, 0, len(result))
	for i := range result {
		objs = append(objs, result[i].item)
	}
	return objs
}

//All items from  root node
func (tree *RTree) All() []*node {
	return all(tree.Data, make([]*node, 0))
}

//all - fetch all items from node
func all(nd *node, result []*node) []*node {
	var nodesToSearch = make([]*node, 0)
	for {
		if nd.leaf {
			for i := range nd.children {
				result = append(result, nd.children[i])
			}
		} else {
			for i := range nd.children {
				nodesToSearch = append(nodesToSearch, nd.children[i])
			}
		}

		nd, nodesToSearch = popNode(nodesToSearch)
		if nd == nil {
			break
		}
	}

	return result
}
