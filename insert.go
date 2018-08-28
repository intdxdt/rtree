package rtree

import (
	"github.com/intdxdt/mbr"
)

//Insert item
func (tree *RTree) Insert(item BoxObj) *RTree {
	if item == nil {
		return tree
	}
	tree.insert(item, tree.Data.height-1)
	return tree
}

//insert - private
func (tree *RTree) insert(item BoxObj, level int) {
	var nd *node
	var insertPath = make([]*node, 0, tree.maxEntries)

	// find the best node for accommodating the item, saving all nodes along the path too
	nd, insertPath = chooseSubtree(item.BBox(), &tree.Data, level, insertPath)

	// put the item into the node item_bbox
	nd.addChild(newLeafNode(item))
	extend(&nd.bbox, item.BBox())

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
	extend(&nd.bbox, &item.bbox)

	// split on node overflow propagate upwards if necessary
	level, insertPath = tree.splitOnOverflow(level, insertPath)

	// adjust bboxes along the insertion path
	tree.adjustParentBBoxes(&item.bbox, insertPath, level)
}

// split on node overflow propagate upwards if necessary
func (tree *RTree) splitOnOverflow(level int, insertPath []*node) (int, []*node) {
	for (level >= 0) && (len(insertPath[level].children)  > tree.maxEntries) {
		tree.split(insertPath, level)
		level--
	}
	return level, insertPath
}

//_chooseSubtree select child of node and updates path to selected node.
func chooseSubtree(bbox *mbr.MBR, nd *node, level int, path []*node) (*node, []*node) {
	var child *node
	var targetNode *node
	var minArea, minEnlargement, area, enlargement float64

	for {
		path = append(path, nd)
		if nd.leaf || (len(path)-1 == level) {
			break
		}
		minArea, minEnlargement = inf, inf

		for i, length := 0, len(nd.children); i < length; i++ {
			child = &nd.children[i]
			area = bboxArea(&child.bbox)
			enlargement = enlargedArea(bbox, &child.bbox) - area

			// choose entry with the least area enlargement
			if enlargement < minEnlargement {
				minEnlargement = enlargement
				if area < minArea {
					minArea = area
				}
				targetNode = child
			} else if feq(enlargement, minEnlargement) {
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

//extend bounding box
func extend(a, b *mbr.MBR) *mbr.MBR {
	return a.ExpandIncludeMBR(b)
}

//computes area of bounding box
func bboxArea(a *mbr.MBR) float64 {
	return a.Area()
}

//computes box margin
func bboxMargin(a *mbr.MBR) float64 {
	return (a.MaxX - a.MinX) + (a.MaxY - a.MinY)
}

//computes enlarged area given two mbrs
func enlargedArea(a, b *mbr.MBR) float64 {
	return (max(a.MaxX, b.MaxX) - min(a.MinX, b.MinX)) * (max(a.MaxY, b.MaxY) - min(a.MinY, b.MinY))
}

//computes the intersection area of two mbrs
func intersectionArea(a, b *mbr.MBR) float64 {
	var minx, miny, maxx, maxy = a.MinX, a.MinY, a.MaxX, a.MaxY

	if !intersects(a, b) {
		return 0.0
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
	return (
		b.MinX >= a.MinX &&
		b.MaxX <= a.MaxX &&
		b.MinY >= a.MinY &&
		b.MaxY <= a.MaxY)
}

//intersects tests a intersect b (MBR)
func intersects(a, b *mbr.MBR) bool {
	return !(
		b.MinX > a.MaxX ||
		b.MaxX < a.MinX ||
		b.MinY > a.MaxY ||
		b.MaxY < a.MinY)
}
