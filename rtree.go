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

func NewRTree(nodeCapacity int) *RTree {
    self  := &RTree{}
    self.Clear()

    if nodeCapacity == 0 {
        nodeCapacity = 9
    }
    // max entries in a node is 9 by default min node fill is 40% for best performance
    self.maxEntries = maxEntries(nodeCapacity)
    self.minEntries = minEntries(self.maxEntries)
    return self
}