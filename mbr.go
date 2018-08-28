package rtree

import (
	"sort"
	"github.com/intdxdt/mbr"
)

//calcBBox calculates its bbox from bboxes of its children.
func calcBBox(nd *node) {
	nd.bbox = distBBox(nd, 0, len(nd.children))
}

//distBBox computes min bounding rectangle of node children from k to p-1.
func distBBox(nd *node, k, p int) mbr.MBR {
	var bbox = emptyMBR()
	for i := k; i < p; i++ {
		bbox.ExpandIncludeMBR(&nd.children[i].bbox)
	}
	return bbox
}

//allDistMargin computes total margin of all possible split distributions.
//Each node is at least m full.
func (tree *RTree) allDistMargin(nd *node, m, M int, sortBy SortBy) float64 {
	if sortBy == byX {
		sort.Sort(XNodePath{nd.children})
	} else if sortBy == byY {
		sort.Sort(YNodePath{nd.children})
	}

	var i int
	var leftBBox = distBBox(nd, 0, m)
	var rightBBox = distBBox(nd, M-m, M)
	var margin = bboxMargin(&leftBBox) + bboxMargin(&rightBBox)

	for i = m; i < M-m; i++ {
		leftBBox.ExpandIncludeMBR(&nd.children[i].bbox)
		margin += bboxMargin(&leftBBox)
	}

	for i = M - m - 1; i >= m; i-- {
		rightBBox.ExpandIncludeMBR(&nd.children[i].bbox)
		margin += bboxMargin(&rightBBox)
	}
	return margin
}
