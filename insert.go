package rtree

import (
	"github.com/intdxdt/mbr"
	"github.com/intdxdt/math"
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
	var node *Node
	var bbox = item.BBox()
	var insertPath = make([]*Node, 0)

	// find the best node for accommodating the item, saving all nodes along the path too
	node, insertPath = chooseSubtree(bbox, tree.Data, level, insertPath)

	if n, ok := item.(*Node); ok {
		node.children = append(node.children, n)
	} else {
		// put the item into the node item_bbox
		node.addChild(newLeafNode(item))
	}
	extend(node.bbox, bbox)

	// split on node overflow propagate upwards if necessary
	for (level >= 0) && (len(insertPath[level].children) > tree.maxEntries) {
		tree.split(insertPath, level)
		level--
	}

	// adjust bboxes along the insertion path
	tree.adjustParentBBoxes(bbox, insertPath, level)
}

//_chooseSubtree select child of node and updates path to selected node.
func chooseSubtree(bbox *mbr.MBR, node *Node, level int, path []*Node) (*Node, []*Node) {

	var child *Node
	var targetNode *Node
	var minArea, minEnlargement, area, enlargement float64

	for {
		path = append(path, node)
		if node.leaf || len(path)-1 == level {
			break
		}
		minEnlargement = math.Inf(1)
		minArea = math.Inf(1)

		for i, length := 0, len(node.children); i < length; i++ {
			child = node.children[i]
			area = bboxArea(child.bbox)
			enlargement = enlargedArea(bbox, child.bbox) - area

			// choose entry with the least area enlargement
			if enlargement < minEnlargement {
				minEnlargement = enlargement
				if area < minArea {
					minArea = area
				}
				targetNode = child
			} else if math.FloatEqual(enlargement, minEnlargement) {
				// otherwise choose one with the smallest area
				if area < minArea {
					minArea = area
					targetNode = child
				}
			}
		}

		node = targetNode
	}

	return node, path
}

//extend bounding box
func extend(a, b *mbr.MBR) *mbr.MBR {
	a[x1] = min(a[x1], b[x1])
	a[y1] = min(a[y1], b[y1])
	a[x2] = max(a[x2], b[x2])
	a[y2] = max(a[y2], b[y2])
	return a
}

//computes area of bounding box
func bboxArea(a *mbr.MBR) float64 {
	return (a[x2] - a[x1]) * (a[y2] - a[y1])
}

//computes box margin
func bboxMargin(a *mbr.MBR) float64 {
	return (a[x2] - a[x1]) + (a[y2] - a[y1])
}

//computes enlarged area given two mbrs
func enlargedArea(a, b *mbr.MBR) float64 {
	return (max(a[x2], b[x2]) - min(a[x1], b[x1])) *
		   (max(a[y2], b[y2]) - min(a[y1], b[y1]))
}

//computes the intersection area of two mbrs
func intersectionArea(a, b *mbr.MBR) float64 {
	var minx, miny, maxx, maxy = a[x1], a[y1], a[x2], a[y2]

	if !intersects(a, b) {
		return 0.0
	}

	if b[x1] > minx {
		minx = b[x1]
	}

	if b[y1] > miny {
		miny = b[y1]
	}

	if b[x2] < maxx {
		maxx = b[x2]
	}

	if b[y2] < maxy {
		maxy = b[y2]
	}

	return (maxx - minx) * (maxy - miny)
}

//contains tests whether a contains b
func contains(a, b *mbr.MBR) bool {
	return b[x1] >= a[x1] &&
		b[x2] <= a[x2] &&
		b[y1] >= a[y1] &&
		b[y2] <= a[y2]
}

//intersects tests a intersect b (mbr)
func intersects(a, b *mbr.MBR) bool {
	return !(
		b[x1] > a[x2] ||
			b[x2] < a[x1] ||
			b[y1] > a[y2] ||
			b[y2] < a[y1])
}
