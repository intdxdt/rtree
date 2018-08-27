package rtree

import (
	"github.com/intdxdt/mbr"
)

//Search item
func (tree *RTree) Search(query mbr.MBR) []BoxObj {
	var bbox = &query
	var result []BoxObj
	var nd = &tree.Data

	if !intersects(bbox, nd.bbox) {
		return []BoxObj{}
	}

	var nodesToSearch []*node
	var child *node
	var childBBox *mbr.MBR

	for {
		for i, length := 0, len(nd.children); i < length; i++ {
			child = &nd.children[i]
			childBBox = child.bbox

			if intersects(bbox, childBBox) {
				if nd.leaf {
					result = append(result, child.item)
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

	//var objs = make([]BoxObj, 0, len(result))
	//for i := range result {
	//	objs = append(objs, result[i].item)
	//}
	return result
}

//All items from  root node
func (tree *RTree) All() []BoxObj {
	return all(&tree.Data, []BoxObj{})
}

//all - fetch all items from node
func all(nd *node, result []BoxObj) []BoxObj {
	var nodesToSearch []*node
	for {
		if nd.leaf {
			for i := range nd.children {
				result = append(result, nd.children[i].item)
			}
		} else {
			for i := range nd.children {
				nodesToSearch = append(nodesToSearch, &nd.children[i])
			}
		}

		nd, nodesToSearch = popNode(nodesToSearch)
		if nd == nil {
			break
		}
	}

	return result
}
