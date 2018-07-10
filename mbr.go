package rtree

import (
	"sort"
	"github.com/intdxdt/mbr"
)

//calcBBox calculates its bbox from bboxes of its children.
func calcBBox(node *Node) {
	node.bbox = distBBox(node, 0, len(node.children))
}

//distBBox computes min bounding rectangle of node children from k to p-1.
func distBBox(node *Node, k, p int) mbr.MBR {
	var bbox = emptyMbr()
	var child *Node
	for i := k; i < p; i++ {
		child = node.children[i]
		extend(&bbox, &child.bbox)
	}
	return bbox
}

//allDistMargin computes total margin of all possible split distributions.
//Each node is at least m full.
func (tree *RTree) allDistMargin(node *Node, m, M int, sortBy SortBy) float64 {
	if sortBy == ByX {
		sort.Sort(XNodePath{node.children})
		//bubbleAxis(*node.getChildren(), ByX, ByY)
	} else if sortBy == ByY {
		sort.Sort(YNodePath{node.children})
		//bubbleAxis(*node.getChildren(), ByY, ByX)
	}

	var i int
	var child *Node

	var leftBBox = distBBox(node, 0, m)
	var rightBBox = distBBox(node, M-m, M)
	var margin = bboxMargin(&leftBBox) + bboxMargin(&rightBBox)

	for i = m; i < M-m; i++ {
		child = node.children[i]
		extend(&leftBBox, &child.bbox)
		margin += bboxMargin(&leftBBox)
	}

	for i = M - m - 1; i >= m; i-- {
		child = node.children[i]
		extend(&rightBBox, &child.bbox)
		margin += bboxMargin(&rightBBox)
	}
	return margin
}
