package rtree

import (
	"github.com/intdxdt/math"
	"github.com/intdxdt/mbr"
)

//Insert item
func (tree *RTree) Insert(item BoxObject) *RTree {
	if item == nil {
		return tree
	}
	tree.insert(item, tree.Data.height-1)
	return tree
}

//insert - private
func (tree *RTree) insert(item BoxObject, level int) {
	var nd *node
	var insertPath = make([]*node, 0, tree.maxEntries)

	// find the best node for accommodating the item, saving all nodes along the path too
	nd, insertPath = chooseSubtree(item.BBox(), &tree.Data, level, insertPath)

	// put the item into the node item_bbox
	nd.addChild(newLeafNode(item))
	nd.bbox.ExpandIncludeMBR(item.BBox())

	// split on node overflow propagate upwards if necessary
	level, insertPath = tree.splitOnOverflow(level, insertPath)

	// adjust bboxes along the insertion path
	tree.adjustParentBBoxes(item.BBox(), insertPath, level)
}

//insert - private
func (tree *RTree) insertNode(item node, level int) {
	var nd *node
	var insertPath []*node

	// find the best node for accommodating the item, saving all nodes along the path too
	nd, insertPath = chooseSubtree(&item.bbox, &tree.Data, level, insertPath)

	nd.children = append(nd.children, item)
	nd.bbox.ExpandIncludeMBR(&item.bbox)

	// split on node overflow propagate upwards if necessary
	level, insertPath = tree.splitOnOverflow(level, insertPath)

	// adjust bboxes along the insertion path
	tree.adjustParentBBoxes(&item.bbox, insertPath, level)
}

// split on node overflow propagate upwards if necessary
func (tree *RTree) splitOnOverflow(level int, insertPath []*node) (int, []*node) {
	for (level >= 0) && (len(insertPath[level].children) > tree.maxEntries) {
		tree.split(insertPath, level)
		level--
	}
	return level, insertPath
}

//_chooseSubtree select child of node and updates path to selected node.
func chooseSubtree(bbox *mbr.MBR, nd *node, level int, path []*node) (*node, []*node) {
	var child, targetNode *node
	var minArea, minEnlargement float64
	var area, enlargement, d float64
	var minx, miny float64
	var maxx, maxy float64
	var ch_minx, ch_miny float64
	var ch_maxx, ch_maxy float64
	var b_minx, b_miny = bbox.MinX, bbox.MinY
	var b_maxx, b_maxy = bbox.MaxX, bbox.MaxY

	var chbox *mbr.MBR

	for {
		path = append(path, nd)
		if nd.leaf || (len(path)-1 == level) {
			break
		}
		minArea, minEnlargement = inf, inf

		for i, length := 0, len(nd.children); i < length; i++ {
			child = &nd.children[i]
			chbox = &child.bbox

			minx, miny = b_minx, b_miny
			maxx, maxy = b_maxx, b_maxy

			ch_minx, ch_miny = chbox.MinX, chbox.MinY
			ch_maxx, ch_maxy = chbox.MaxX, chbox.MaxY

			if ch_minx < minx {
				minx = ch_minx
			}
			if ch_miny < miny {
				miny = ch_miny
			}
			if ch_maxx > maxx {
				maxx = ch_maxx
			}
			if ch_maxy > maxy {
				maxy = ch_maxy
			}

			area = (ch_maxx - ch_minx) * (ch_maxy - ch_miny)
			enlargement = (maxx - minx) * (maxy - miny)

			d = enlargement - minEnlargement
			// choose entry with the least area enlargement
			if d < 0 {
				minEnlargement = enlargement
				if area < minArea {
					minArea = area
				}
				targetNode = child
			} else if d == 0 || math.Abs(d) < math.EPSILON {
				// otherwise choose one with the smallest area
				if area < minArea {
					minArea = area
					targetNode = child
				}
			}
		}

		nd = targetNode
	}

	return nd, path
}

//computes box margin
func bboxMargin(a *mbr.MBR) float64 {
	return (a.MaxX - a.MinX) + (a.MaxY - a.MinY)
}

//computes the intersection area of two mbrs
func intersectionArea(a, b *mbr.MBR) float64 {
	var minx, miny, maxx, maxy = a.MinX, a.MinY, a.MaxX, a.MaxY

	if !intersects(a, b) {
		return 0
	}

	if b.MinX > minx {
		minx = b.MinX
	}

	if b.MinY > miny {
		miny = b.MinY
	}

	if b.MaxX < maxx {
		maxx = b.MaxX
	}

	if b.MaxY < maxy {
		maxy = b.MaxY
	}

	return (maxx - minx) * (maxy - miny)
}

//contains tests whether a contains b
func contains(a, b *mbr.MBR) bool {
	return b.MinX >= a.MinX && b.MaxX <= a.MaxX && b.MinY >= a.MinY && b.MaxY <= a.MaxY
}

//intersects tests a intersect b (MBR)
func intersects(a, b *mbr.MBR) bool {
	return !(b.MinX > a.MaxX || b.MaxX < a.MinX || b.MinY > a.MaxY || b.MaxY < a.MinY)
}
