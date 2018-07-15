package rtree

import (
	"math"
	"sort"
)

// _split overflowed rNode into two
func (tree *RTree) split(insertPath []*rNode, level int) {
	var node = insertPath[level]
	var newNode = newNode(emptyObject(), node.height, node.leaf, []*rNode{})
	var M = len(node.children)
	var m = tree.minEntries

	tree.chooseSplitAxis(node, m, M)
	var at = tree.chooseSplitIndex(node, m, M)
	//perform split at index
	node.children, newNode.children = splitAtIndex(node.children, at)

	calcBBox(node)
	calcBBox(newNode)

	if level > 0 {
		insertPath[level-1].addChild(newNode)
	} else {
		tree.splitRoot(node, newNode)
	}
}

//_splitRoot splits the root of tree.
func (tree *RTree) splitRoot(node, other *rNode) {
	// split root rNode
	tree.Data = newNode(emptyObject(), node.height+1, false, []*rNode{node, other})
	calcBBox(tree.Data)
}

//_chooseSplitIndex selects split index.
func (tree *RTree) chooseSplitIndex(node *rNode, m, M int) int {
	var i, index int
	var overlap, area, minOverlap, minArea float64

	minOverlap, minArea = math.Inf(1), math.Inf(1)

	for i = m; i <= M-m; i++ {
		var bbox1 = distBBox(node, 0, i)
		var bbox2 = distBBox(node, i, M)

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

//_chooseSplitAxis selects split axis : sorts rNode children
//by the best axis for split.
func (tree *RTree) chooseSplitAxis(node *rNode, m, M int) {
	var xMargin = tree.allDistMargin(node, m, M, ByX)
	var yMargin = tree.allDistMargin(node, m, M, ByY)

	// if total distributions margin value is minimal for x, sort by minX,
	// otherwise it's already sorted by minY
	if xMargin < yMargin {
		sort.Sort(&XNodePath{node.children})
	}
}
