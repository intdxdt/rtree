package rtree

import "github.com/intdxdt/mbr"

func (tree *RTree) Collides(query mbr.MBR) bool {
    var bbox = &query
    var nd = tree.Data
    if !intersects(bbox, nd.bbox) {
        return false
    }

    var searchList = make([]*node, 0)
    var child *node
    var bln  = false

    for !bln && nd != nil {
        for i, length := 0, len(nd.children); !bln && i < length; i++ {
            child = nd.children[i]
            if intersects(bbox, child.bbox) {
                bln =  nd.leaf || contains(bbox, child.bbox)
                searchList = append(searchList, child)
            }
        }
        nd, searchList = popNode(searchList)
    }

    return bln
}
