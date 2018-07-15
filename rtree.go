package rtree
/*
 (c) 2015, Titus Tienaah
 A library for 2D spatial indexing of points and rectangles.
 https://github.com/mourner/rbush
 @after  (c) 2015, Vladimir Agafonkin
*/

//RTree type
type RTree struct {
    Data       *Node
    maxEntries int
    minEntries int
}

func NewRTree(nodeCap ...int) *RTree {
    var bucketSize = 8
    var self  = (&RTree{}).Clear()
    if len(nodeCap) > 0 {
        bucketSize = nodeCap[0]
    }
    // max entries in a Node is 9 by default min Node fill is 40% for best performance
    self.maxEntries = maxEntries(bucketSize)
    self.minEntries = minEntries(self.maxEntries)
    return self
}