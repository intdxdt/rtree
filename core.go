package rtree

import (
	"github.com/intdxdt/math"
	"github.com/intdxdt/mbr"
)

var inf = math.Inf(1)
var feq =  math.FloatEqual
type SortBy int

const (
	byX SortBy = iota
	byY
)

const (
	x1 = iota
	y1
	x2
	y2
)

type compareNode func(BoxObj, BoxObj) float64
type BoxObj interface {
    BBox() *mbr.MBR
}

func maxEntries(x int) int {
	return maxInt(4, x)
}

func minEntries(x int) int {
	return maxInt(2, int(math.Ceil(float64(x)*0.4)))
}

//compareNodeMinX computes change in minimum x
func compareNodeMinX(a, b BoxObj) float64 {
	return a.BBox().MinX - b.BBox().MinX
}

//compareNodeMinY computes change in minimum y
func compareNodeMinY(a, b BoxObj) float64 {
	return a.BBox().MinY - b.BBox().MinY
}

func swapItem(arr []BoxObj, i, j int) {
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
