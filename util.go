package rtree

//split at index
func splitAtIndex(nodes []*Node, index int) ([]*Node, []*Node) {
	ln := len(nodes)
	newNodes := make([]*Node, ln - index)
	copy(newNodes, nodes[index:])
	for i := index; i < ln; i++ {
		nodes[i] = nil
	}
	nodes = nodes[:index]
	return nodes, newNodes
}

//slice index
func sliceIndex(limit int, predicate func(i int) bool) int {
	for i := 0; i < limit; i++ {
		if predicate(i) {
			return i
		}
	}
	return -1
}

//minimum float
func min(a, b float64) float64 {
	if b < a {
		return b
	}
	return a
}

//maximum float
func max(a, b float64) float64 {
	if b > a {
		return b
	}
	return a
}

//min integer
func minInt(a, b int) int {
	if b < a {
		return b
	}
	return a
}

//maximum integer
func maxInt(a, b int) int {
	if b > a {
		return b
	}
	return a
}
