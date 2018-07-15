package rtree

//split at index
func splitAtIndex(nodes []*rNode, index int) ([]*rNode, []*rNode) {
	var ln = len(nodes)
	var ext = make([]*rNode, 0, ln-index)
	for i := index; i < ln; i++ {
		ext = append(ext, nodes[i])
		nodes[i] = nil
	}
	return nodes[:index], ext
}

//slice index
func sliceIndex(limit int, predicate func(i int) bool) int {
	var index = -1
	for i := 0; i < limit; i++ {
		if predicate(i) {
			index = i
			break
		}
	}
	return index
}

//minimum float
func min(a, b float64) float64 {
	var m = a
	if b < a {
		m = b
	}
	return m
}

//maximum float
func max(a, b float64) float64 {
	var m = a
	if b > a {
		m = b
	}
	return m
}

//min integer
func minInt(a, b int) int {
	var m = a
	if b < a {
		m = b
	}
	return m
}

//maximum integer
func maxInt(a, b int) int {
	var m = a
	if b > a {
		m = b
	}
	return m
}
