package rtree

import "github.com/intdxdt/mbr"

func (tree *RTree) adjustParentBBoxes(bbox *mbr.MBR, path []*Node, level int) {
	// adjust bboxes along the given tree path
	for i := level; i >= 0; i-- {
		extend(path[i].bbox, bbox)
	}
}
