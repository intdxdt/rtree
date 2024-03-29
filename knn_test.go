package rtree

import (
	"github.com/franela/goblin"
	"github.com/intdxdt/mbr"
	"testing"
	"time"
)

var knnData []BoxObject

func scoreFn(query *mbr.MBR[float64], boxer *KObj) float64 {
	return query.Distance(boxer.MBR)
}

func initKnn() {

	var dat = []mbr.MBR[float64]{
		{87, 55, 87, 56}, {38, 13, 39, 16}, {7, 47, 8, 47}, {89, 9, 91, 12}, {4, 58, 5, 60}, {0, 11, 1, 12}, {0, 5, 0, 6}, {69, 78, 73, 78},
		{56, 77, 57, 81}, {23, 7, 24, 9}, {68, 24, 70, 26}, {31, 47, 33, 50}, {11, 13, 14, 15}, {1, 80, 1, 80}, {72, 90, 72, 91}, {59, 79, 61, 83},
		{98, 77, 101, 77}, {11, 55, 14, 56}, {98, 4, 100, 6}, {21, 54, 23, 58}, {44, 74, 48, 74}, {70, 57, 70, 61}, {32, 9, 33, 12}, {43, 87, 44, 91},
		{38, 60, 38, 60}, {62, 48, 66, 50}, {16, 87, 19, 91}, {5, 98, 9, 99}, {9, 89, 10, 90}, {89, 2, 92, 6}, {41, 95, 45, 98}, {57, 36, 61, 40},
		{50, 1, 52, 1}, {93, 87, 96, 88}, {29, 42, 33, 42}, {34, 43, 36, 44}, {41, 64, 42, 65}, {87, 3, 88, 4}, {56, 50, 56, 52}, {32, 13, 35, 15},
		{3, 8, 5, 11}, {16, 33, 18, 33}, {35, 39, 38, 40}, {74, 54, 78, 56}, {92, 87, 95, 90}, {12, 97, 16, 98}, {76, 39, 78, 40}, {16, 93, 18, 95},
		{62, 40, 64, 42}, {71, 87, 71, 88}, {60, 85, 63, 86}, {39, 52, 39, 56}, {15, 18, 19, 18}, {91, 62, 94, 63}, {10, 16, 10, 18}, {5, 86, 8, 87},
		{85, 85, 88, 86}, {44, 84, 44, 88}, {3, 94, 3, 97}, {79, 74, 81, 78}, {21, 63, 24, 66}, {16, 22, 16, 22}, {68, 97, 72, 97}, {39, 65, 42, 65},
		{51, 68, 52, 69}, {61, 38, 61, 42}, {31, 65, 31, 65}, {16, 6, 19, 6}, {66, 39, 66, 41}, {57, 32, 59, 35}, {54, 80, 58, 84}, {5, 67, 7, 71},
		{49, 96, 51, 98}, {29, 45, 31, 47}, {31, 72, 33, 74}, {94, 25, 95, 26}, {14, 7, 18, 8}, {29, 0, 31, 1}, {48, 38, 48, 40}, {34, 29, 34, 32},
		{99, 21, 100, 25}, {79, 3, 79, 4}, {87, 1, 87, 5}, {9, 77, 9, 81}, {23, 25, 25, 29}, {83, 48, 86, 51}, {79, 94, 79, 95}, {33, 95, 33, 99},
		{1, 14, 1, 14}, {33, 77, 34, 77}, {94, 56, 98, 59}, {75, 25, 78, 26}, {17, 73, 20, 74}, {11, 3, 12, 4}, {45, 12, 47, 12}, {38, 39, 39, 39},
		{99, 3, 103, 5}, {41, 92, 44, 96}, {79, 40, 79, 41}, {29, 2, 29, 4},
	}
	knnData = make([]BoxObject, 0, len(dat))
	for i := range dat {
		knnData = append(knnData, &dat[i])
	}
}

func foundIn(needle *mbr.MBR[float64], haystack []mbr.MBR[float64]) bool {
	var found = false
	for i := range haystack {
		found = needle.Equals(&haystack[i])
		if found {
			break
		}
	}
	return found
}

func TestRtreeKNN(t *testing.T) {
	g := goblin.Goblin(t)
	g.Describe("Rtree Knn", func() {
		initKnn()
		g.It("finds n neighbours", func() {
			var rt = NewRTree(9)
			rt.Load(knnData)
			var nn = rt.Knn(mbr.CreateMBR[float64](40, 40, 40, 40), 10, scoreFn)
			var result = []mbr.MBR[float64]{
				{38, 39, 39, 39}, {35, 39, 38, 40}, {34, 43, 36, 44},
				{29, 42, 33, 42}, {48, 38, 48, 40}, {31, 47, 33, 50},
				{34, 29, 34, 32}, {29, 45, 31, 47}, {39, 52, 39, 56},
				{57, 36, 61, 40},
			}
			for _, n := range nn {
				g.Assert(foundIn(n.BBox(), result)).IsTrue()
			}

			nn = rt.Knn(mbr.CreateMBR[float64](40, 40, 40, 40), 1000, scoreFn)
			g.Assert(len(nn)).Equal(len(knnData))
		})
	})
}

func TestRtreeKNNPredScore(t *testing.T) {
	g := goblin.Goblin(t)
	g.Describe("Rtree Knn - Pred, Score", func() {

		initKnn()
		g.It("finds n neighbours with geoms", func() {

			var scoreFunc = func(query *mbr.MBR[float64], obj *KObj) float64 {
				return query.Distance(obj.MBR)
			}

			var predicateMbr []*mbr.MBR[float64]

			var createPredicate = func(dist float64) func(*KObj) (bool, bool) {
				return func(candidate *KObj) (bool, bool) {
					g.Assert(candidate.IsItem).IsTrue()
					if candidate.Dist <= dist {
						predicateMbr = append(predicateMbr, candidate.MBR)
						return true, false
					}
					return false, true
				}
			}
			rt := NewRTree(9)
			rt.Load(knnData)
			prefFn := createPredicate(6)
			query := mbr.CreateMBR(
				74.88825108886668, 82.678427498132,
				74.88825108886668, 82.678427498132,
			)

			var res = rt.Knn(query, 10, scoreFunc, prefFn)

			g.Assert(len(res)).Equal(2)
			for i, r := range res {
				g.Assert(r.BBox().Equals(predicateMbr[i])).IsTrue()
			}
		})
	})
}

type RichData struct {
	*mbr.MBR[float64]
	version int
}

func fnRichData() []BoxObject {
	var richData = make([]BoxObject, 0)
	var data = []mbr.MBR[float64]{
		{1, 2, 1, 2}, {3, 3, 3, 3}, {5, 5, 5, 5},
		{4, 2, 4, 2}, {2, 4, 2, 4}, {5, 3, 5, 3},
		{3, 4, 3, 4}, {2.5, 4, 2.5, 4},
	}
	for i := range data {
		richData = append(richData, &RichData{&data[i], i + 1})
	}
	return richData
}

func TestQobj_String(t *testing.T) {
	g := goblin.Goblin(t)
	g.Describe("", func() {
		g.It("test qobject", func() {
			var box = mbr.MBR[float64]{3, 3, 3, 3}
			var nd = newLeafNode(&box)
			var qo = &KObj{&nd, &box, true, 3.4}
			g.Assert(box.String() + " -> 3.4").Equal(qo.String())
		})
	})
}

func TestRtreeKNNPredicate(t *testing.T) {
	g := goblin.Goblin(t)
	g.Describe("Rtree Knn Predicate", func() {
		g.It("find n neighbours that do satisfy a given predicate", func() {
			g.Timeout(1 * time.Hour)
			var rt = NewRTree(9)
			rt.Load(fnRichData())

			var scoreFn = func(query *mbr.MBR[float64], boxer *KObj) float64 {
				return query.Distance(boxer.MBR)
			}
			var predicate = func(v *KObj) (bool, bool) {
				return v.GetItem().(*RichData).version < 5, false
			}
			var result = rt.Knn(mbr.CreateMBR[float64](2, 4, 2, 4), 1, scoreFn, predicate)

			g.Assert(len(result)).Equal(1)

			var v = result[0].(*RichData)
			var expectsVersion = 2

			g.Assert(v.MBR.Equals(&mbr.MBR[float64]{3, 3, 3, 3})).IsTrue()
			g.Assert(v.version).Equal(expectsVersion)
		})
	})
}
