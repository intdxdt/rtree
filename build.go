package rtree

import "math"

//build
func (tree *RTree) buildTree(items []BoxObject, left, right, height int) node {
	var N = float64(right - left + 1)
	var M = float64(tree.maxEntries)
	//var n *node
	if N <= M {
		// reached leaf level return leaf
		var n = createNode(
			nil , 1, true,
			makeChildren(items[left:right+1:right+1]),
		)
		calcBBox(&n)
		return n
	}

	if height == 0 {
		// target height of the bulk-loaded tree
		height = int(
			math.Ceil(math.Log(N) / math.Log(M)))

		// target number of root entries to maximize storage utilization
		M = math.Ceil(N / math.Pow(M, float64(height-1)))
	}

	// TODO eliminate recursion?
	var n = createNode(nil, height, false, []node{})

	// split the items into M mostly square tiles

	var N2 = int(math.Ceil(N / M))
	var N1 = N2 * int(math.Ceil(math.Sqrt(M)))
	var i, j, right2, right3 int

	multiSelect(items, left, right, N1, cmpMinX)

	for i = left; i <= right; i += N1 {
		right2 = minInt(i+N1-1, right)
		multiSelect(items, i, right2, N2, cmpMinY)

		for j = i; j <= right2; j += N2 {
			right3 = minInt(j+N2-1, right2)
			// pack each entry recursively
			n.addChild(tree.buildTree(items, j, right3, height-1))
		}
	}

	calcBBox(&n)
	return n
}
