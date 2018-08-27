package rtree

import "github.com/intdxdt/mbr"

// adjust bboxes along the given tree path
func (tree *RTree) adjustParentBBoxes(bbox *mbr.MBR, path []*node, level int) {
	for i := level; i >= 0; i-- {
		extend(path[i].bbox, bbox)
	}
}
