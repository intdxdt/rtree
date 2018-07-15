package rtree

import "github.com/intdxdt/mbr"

func (tree *RTree) Collides(bbox *mbr.MBR) bool {
    var node = tree.Data
    if !intersects(bbox, node.bbox) {
        return false
    }

    var searchList = make([]*rNode, 0)
    var child *rNode
    var bln  = false

    for !bln && node != nil {
        for i, length := 0, len(node.children); !bln && i < length; i++ {
            child = node.children[i]
            if intersects(bbox, child.bbox) {
                bln =  node.leaf || contains(bbox, child.bbox)
                searchList = append(searchList, child)
            }
        }
        node, searchList = popNode(searchList)
    }

    return bln
}
