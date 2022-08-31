package rtree

import (
	"math"
	"sort"
)

// _split overflowed node into two
func (tree *RTree) split(insertPath []*node, level int) {
	var nd = insertPath[level]
	var newNode = createNode(nil, nd.height, nd.leaf, []node{})
	var M = len(nd.children)
	var m = tree.minEntries

	tree.chooseSplitAxis(nd, m, M)
	var at = tree.chooseSplitIndex(nd, m, M)
	//perform split at index
	nd.children, newNode.children = splitAtIndex(nd.children, at)

	calcBBox(nd)
	calcBBox(&newNode)

	if level > 0 {
		insertPath[level-1].addChild(newNode)
	} else {
		tree.splitRoot(*nd, newNode)
	}
}

// _splitRoot splits the root of tree.
func (tree *RTree) splitRoot(nd, other node) {
	// split root node
	tree.Data = createNode(nil, nd.height+1, false, []node{nd, other})
	calcBBox(&tree.Data)
}

// chooseSplitIndex selects split index.
func (tree *RTree) chooseSplitIndex(nd *node, m, M int) int {
	var i int
	var index = -1
	var overlap, area, minOverlap, minArea float64

	minOverlap, minArea = math.Inf(1), math.Inf(1)

	for i = m; i <= M-m; i++ {
		var bbox1 = distBBox(nd, 0, i)
		var bbox2 = distBBox(nd, i, M)

		overlap = intersectionArea(&bbox1, &bbox2)
		area = bbox1.Area() + bbox2.Area()

		// choose distribution with minimum overlap
		if overlap < minOverlap {
			minOverlap = overlap
			index = i
			if area < minArea {
				minArea = area
			}
		} else if overlap == minOverlap {
			// otherwise choose distribution with minimum area
			if area < minArea {
				minArea = area
				index = i
			}
		}
	}

	if index < 0 {
		return M - m
	} //else
	return index
}

// _chooseSplitAxis selects split axis : sorts node children
// by the best axis for split.
func (tree *RTree) chooseSplitAxis(nd *node, m, M int) {
	var xMargin = tree.allDistMargin(nd, m, M, byX)
	var yMargin = tree.allDistMargin(nd, m, M, byY)

	// if total distributions margin value is minimal for x, sort by minX,
	// otherwise it's already sorted by minY
	if xMargin < yMargin {
		sort.Sort(&xNodePath{nd.children})
	}
}
