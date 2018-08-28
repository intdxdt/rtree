package rtree

import (
	"github.com/intdxdt/math"
	"github.com/intdxdt/mbr"
)
const (
	byX SortBy = iota
	byY
)

type SortBy int
type compareObject func(BoxObject, BoxObject) float64
type BoxObject interface {
    BBox() *mbr.MBR
}

var inf = math.Inf(1)


func maxEntries(x int) int {
	return maxInt(4, x)
}

func minEntries(x int) int {
	return maxInt(2, int(math.Ceil(float64(x)*0.4)))
}

//compareMinX computes change in minimum x
func compareMinX(a, b BoxObject) float64 {
	return a.BBox().MinX - b.BBox().MinX
}

//compareMinY computes change in minimum y
func compareMinY(a, b BoxObject) float64 {
	return a.BBox().MinY - b.BBox().MinY
}

func swapItem(arr []BoxObject, i, j int) {
	arr[i], arr[j] = arr[j], arr[i]
}

func popInt(a []int) (int, []int) {
	var n, v int
	n = len(a) - 1
	v, a[n] = a[n], 0
	a = a[:n]
	return v, a
}

func appendInts(a []int, v ...int) []int {
	for i := range v {
		a = append(a, v[i])
	}
	return a
}
