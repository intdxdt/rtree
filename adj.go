package rtree

import "github.com/intdxdt/mbr"

// adjust bboxes along the given tree path
func (tree *RTree) adjustParentBBoxes(bbox *mbr.MBR[float64], path []*node, level int) {
	for i := level; i >= 0; i-- {
		path[i].bbox.ExpandIncludeMBR(bbox)
	}
}
