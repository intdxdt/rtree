package rtree

import (
    "github.com/intdxdt/mbr"
)

//LoadBoxes loads bounding boxes
func (tree *RTree)  LoadBoxes(data []*mbr.MBR) *RTree {
    items := make([]BoxObj, len(data))
    for i:=0 ; i < len(data) ; i++{
        items[i] = data[i]
    }
    return tree.Load(items)
}

//Load implements bulk loading
func (tree *RTree) Load(items []BoxObj) *RTree {
    var node *Node
    //if len(items) == 0 {
    //    return tree
    //}

    if len(items) < tree.minEntries {
        for _, o := range items {
            tree.Insert(o)
        }
        return tree
    }

    var data = make([]BoxObj, len(items))
	copy(data, items)

    // recursively build the tree with the given data from stratch using OMT algorithm
    node = tree._build(data, 0, len(data) - 1, 0)

    if len(tree.Data.children) == 0 {
        // save as is if tree is empty
        tree.Data = node
    }  else if tree.Data.height == node.height {
        // split root if trees have the same height
        tree.splitRoot(tree.Data, node)
    }  else {
        if tree.Data.height < node.height {
            // swap trees if inserted one is bigger
            tree.Data, node = node, tree.Data
        }

        // insert the small tree into the large tree at appropriate level
        tree.insert(node, tree.Data.height - node.height - 1)
    }

    return tree
}


