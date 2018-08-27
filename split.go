package rtree

import (
	"math"
	"sort"
)

// _split overflowed node into two
func (tree *RTree) split(insertPath []*node, level int) {
	var nd = insertPath[level]
	var newNode = newNode(universe{}, nd.height, nd.leaf, []node{})
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

//_splitRoot splits the root of tree.
func (tree *RTree) splitRoot(nd, other node) {
	// split root node
	tree.Data = newNode(universe{}, nd.height+1, false, []node{nd, other})
	calcBBox(&tree.Data)
}

//_chooseSplitIndex selects split index.
func (tree *RTree) chooseSplitIndex(nd *node, m, M int) int {
	var i, index int
	var overlap, area, minOverlap, minArea float64

	minOverlap, minArea = math.Inf(1), math.Inf(1)

	for i = m; i <= M-m; i++ {
		var bbox1 = distBBox(nd, 0, i)
		var bbox2 = distBBox(nd, i, M)

		overlap = intersectionArea(bbox1, bbox2)
		area = bboxArea(bbox1) + bboxArea(bbox2)

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

	return index
}

//_chooseSplitAxis selects split axis : sorts node children
//by the best axis for split.
func (tree *RTree) chooseSplitAxis(nd *node, m, M int) {
	var xMargin = tree.allDistMargin(nd, m, M, ByX)
	var yMargin = tree.allDistMargin(nd, m, M, ByY)

	// if total distributions margin value is minimal for x, sort by minX,
	// otherwise it's already sorted by minY
	if xMargin < yMargin {
		sort.Sort(&XNodePath{nd.children})
	}
}
