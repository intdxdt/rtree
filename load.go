package rtree

import (
	"github.com/intdxdt/mbr"
)

//LoadBoxes loads bounding boxes
func (tree *RTree) LoadBoxes(data []*mbr.MBR) *RTree {
	var items = make([]*Obj, 0, len(data))
	for i := range data {
		items = append(items, &Obj{Id: i, MBR: data[i]})
	}
	return tree.Load(items)
}

//Load implements bulk loading
func (tree *RTree) Load(items []*Obj) *RTree {
	var n  = len(items)
	var node *rNode
	if n < tree.minEntries {
		for i := range items {
			tree.Insert(items[i])
		}
		return tree
	}

	var data = make([]*Obj, 0, n)
	for i := range items {
		data = append(data, items[i])
	}


	// recursively build the tree with the given data from stratch using OMT algorithm
	node = tree.buildTree(data, 0, n-1, 0)

	if len(tree.Data.children) == 0 {
		// save as is if tree is empty
		tree.Data = node
	} else if tree.Data.height == node.height {
		// split root if trees have the same height
		tree.splitRoot(tree.Data, node)
	} else {
		if tree.Data.height < node.height {
			// swap trees if inserted one is bigger
			tree.Data, node = node, tree.Data
		}

		// insert the small tree into the large tree at appropriate level
		tree.insertNode(node, tree.Data.height-node.height-1)
	}

	return tree
}
