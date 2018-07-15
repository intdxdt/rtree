package rtree

import (
	"sort"
	"testing"
	"github.com/intdxdt/mbr"
	"github.com/franela/goblin"
)

type Boxes []*mbr.MBR

//Len for sort interface
func (o Boxes) Len() int {
	return len(o)
}

//Swap for sort interface
func (o Boxes) Swap(i, j int) {
	o[i], o[j] = o[j], o[i]
}

//Less sorts boxes lexicographically
func (o Boxes) Less(i, j int) bool {
	var x, y = 0, 1
	var a, b = o[i].BBox(), o[j].BBox()
	var d = a[x] - b[x]
	//x's are close enough to each other
	if feq(d, 0.0) {
		d = a[y] - b[y]
	}
	//check if close enough ot zero
	return d < 0
}

func someData(n int) []*mbr.MBR {
	var data = make([]*mbr.MBR, n)
	for i := 0; i < n; i++ {
		data[i] = mbr.New(float64(i), float64(i), float64(i), float64(i))
	}
	return data
}

func testResults(g *goblin.G, nodes []*Node, boxes Boxes) {
	sort.Sort(&XYNodePath{nodes})
	sort.Sort(boxes)
	g.Assert(len(nodes)).Equal(len(boxes))
	for i, n := range nodes {
		g.Assert(n.bbox.Equals(boxes[i])).IsTrue()
	}
}

//data from rbush 1.4.2
var data = []*mbr.MBR{{0, 0, 0, 0}, {10, 10, 10, 10}, {20, 20, 20, 20}, {25, 0, 25, 0}, {35, 10, 35, 10}, {45, 20, 45, 20}, {0, 25, 0, 25}, {10, 35, 10, 35},
	{20, 45, 20, 45}, {25, 25, 25, 25}, {35, 35, 35, 35}, {45, 45, 45, 45}, {50, 0, 50, 0}, {60, 10, 60, 10}, {70, 20, 70, 20}, {75, 0, 75, 0},
	{85, 10, 85, 10}, {95, 20, 95, 20}, {50, 25, 50, 25}, {60, 35, 60, 35}, {70, 45, 70, 45}, {75, 25, 75, 25}, {85, 35, 85, 35}, {95, 45, 95, 45},
	{0, 50, 0, 50}, {10, 60, 10, 60}, {20, 70, 20, 70}, {25, 50, 25, 50}, {35, 60, 35, 60}, {45, 70, 45, 70}, {0, 75, 0, 75}, {10, 85, 10, 85},
	{20, 95, 20, 95}, {25, 75, 25, 75}, {35, 85, 35, 85}, {45, 95, 45, 95}, {50, 50, 50, 50}, {60, 60, 60, 60}, {70, 70, 70, 70}, {75, 50, 75, 50},
	{85, 60, 85, 60}, {95, 70, 95, 70}, {50, 75, 50, 75}, {60, 85, 60, 85}, {70, 95, 70, 95}, {75, 75, 75, 75}, {85, 85, 85, 85}, {95, 95, 95, 95}}

func TestRtreeRbush(t *testing.T) {
	g := goblin.Goblin(t)

	g.Describe("Rtree Tests - From Rbush", func() {
		g.It("should test load 9 & 10", func() {
			var tree0 = NewRTree().LoadBoxes(someData(0))
			g.Assert(tree0.Data.height).Equal(1)

			var tree1 = NewRTree(8).LoadBoxes(someData(8))
			g.Assert(tree1.Data.height).Equal(1)

			var tree2 = NewRTree(8).LoadBoxes(someData(10))
			g.Assert(tree2.Data.height).Equal(2)
		})

		g.It("tests search with some other", func() {
			var data = []*mbr.MBR{
				{-115, 45, -105, 55}, {105, 45, 115, 55}, {105, -55, 115, -45}, {-115, -55, -105, -45},
			}
			var tree = NewRTree(4)
			tree.LoadBoxes(data)
			testResults(g, tree.Search(mbr.New(-180, -90, 180, 90)), []*mbr.MBR{
				{-115, 45, -105, 55}, {105, 45, 115, 55}, {105, -55, 115, -45}, {-115, -55, -105, -45},
			})

			testResults(g, tree.Search(mbr.New(-180, -90, 0, 90)), []*mbr.MBR{
				{-115, 45, -105, 55}, {-115, -55, -105, -45},
			})
			testResults(g, tree.Search(mbr.New(0, -90, 180, 90)), []*mbr.MBR{
				{105, 45, 115, 55}, {105, -55, 115, -45},
			})
			testResults(g, tree.Search(mbr.New(-180, 0, 180, 90)), []*mbr.MBR{
				{-115, 45, -105, 55}, {105, 45, 115, 55},
			})
			testResults(g, tree.Search(mbr.New(-180, -90, 180, 0)), []*mbr.MBR{
				{105, -55, 115, -45}, {-115, -55, -105, -45},
			})
		})

		g.It("#load uses standard insertion when given a low number of items", func() {
			var tree = NewRTree(8).LoadBoxes(data)
			tree.LoadBoxes(data[0:3])
			var tree2 = NewRTree(8).LoadBoxes(data).Insert(
				&Obj{Id: 0, MBR: data[0]},
			).Insert(&Obj{Id: 1, MBR: data[1]}).Insert(&Obj{Id: 2, MBR: data[2]})
			g.Assert(tree.Data).Eql(tree2.Data)
		})

		g.It("#load does nothing if loading empty data", func() {
			var tree = NewRTree(0).Load(make([]*Obj, 0))
			g.Assert(tree.IsEmpty()).IsTrue()
		})

		g.It("#load properly splits tree root when merging trees of the same height", func() {
			var cloneData = make([]*mbr.MBR, len(data))
			for i := 0; i < len(data); i++ {
				cloneData[i] = data[i].Clone()
			}
			for i := 0; i < len(data); i++ {
				cloneData = append(cloneData, data[i].Clone())
			}
			var tree = NewRTree(4).LoadBoxes(data).LoadBoxes(data)
			testResults(g, tree.All(), cloneData)
		})

		g.It("#load properly merges data of smaller or bigger tree heights", func() {
			var smaller = someData(10)
			var cloneData = make([]*mbr.MBR, len(data))
			for i := 0; i < len(data); i++ {
				cloneData[i] = data[i].Clone()
			}
			for i := 0; i < len(smaller); i++ {
				cloneData = append(cloneData, smaller[i].Clone())
			}

			var tree1 = NewRTree(4).LoadBoxes(data).LoadBoxes(smaller)
			var tree2 = NewRTree(4).LoadBoxes(smaller).LoadBoxes(data)
			g.Assert(tree1.Data.height).Equal(tree2.Data.height)
			testResults(g, tree1.All(), cloneData)
			testResults(g, tree2.All(), cloneData)
		})

		g.It("#load properly merges data of smaller or bigger tree heights 2", func() {
			N = 8020
			var smaller = GenDataItems(N, 1)
			var larger = GenDataItems(2*N, 1)
			var cloneData = make([]*mbr.MBR, len(larger))

			for i := 0; i < len(larger); i++ {
				cloneData[i] = larger[i].Clone()
			}
			for i := 0; i < len(smaller); i++ {
				cloneData = append(cloneData, smaller[i].Clone())
			}

			var tree1 = NewRTree(64).LoadBoxes(larger).LoadBoxes(smaller)
			var tree2 = NewRTree(64).LoadBoxes(smaller).LoadBoxes(larger)
			g.Assert(tree1.Data.height).Equal(tree2.Data.height)
			testResults(g, tree1.All(), cloneData)
			testResults(g, tree2.All(), cloneData)
		})

		g.It("#search finds matching points in the tree given a bbox", func() {
			var tree = NewRTree(4).LoadBoxes(data)
			var result = tree.Search(mbr.New(40, 20, 80, 70))
			testResults(g, result, []*mbr.MBR{
				{70, 20, 70, 20}, {75, 25, 75, 25}, {45, 45, 45, 45}, {50, 50, 50, 50}, {60, 60, 60, 60}, {70, 70, 70, 70},
				{45, 20, 45, 20}, {45, 70, 45, 70}, {75, 50, 75, 50}, {50, 25, 50, 25}, {60, 35, 60, 35}, {70, 45, 70, 45},
			})
		})

		g.It("#collides returns true when search finds matching points", func() {
			var tree = NewRTree(4).LoadBoxes(data)
			g.Assert(tree.Collides(mbr.New(40, 20, 80, 70))).IsTrue()
			g.Assert(tree.Collides(mbr.New(200, 200, 210, 210))).IsFalse()
		})

		g.It("#search returns an empty array if nothing found", func() {
			var result = NewRTree(4).LoadBoxes(data).Search(mbr.New(200, 200, 210, 210))
			g.Assert(len(result)).Equal(0)
		})

		g.It("#all <==>.Data returns all points in the tree", func() {
			var cloneData = make([]*mbr.MBR, len(data))
			for i := 0; i < len(data); i++ {
				cloneData[i] = data[i]
			}

			var tree = NewRTree(4).LoadBoxes(data)
			var result = tree.Search(mbr.New(0, 0, 100, 100))
			testResults(g, result, cloneData)
		})

		g.It("#insert adds an item to an existing tree correctly", func() {
			var data = []*mbr.MBR{{0, 0, 0, 0}, {2, 2, 2, 2}, {1, 1, 1, 1},}
			var tree = NewRTree(4)
			tree.LoadBoxes(data)
			tree.Insert(Object(0, mbr.New(3, 3, 3, 3)))
			g.Assert(tree.Data.leaf).IsTrue()
			g.Assert(tree.Data.height).Equal(1)
			var box = mbr.New(0, 0, 3, 3)
			g.Assert(tree.Data.bbox.Equals(box)).IsTrue()
			testResults(g, tree.Data.children, []*mbr.MBR{
				{0, 0, 0, 0}, {1, 1, 1, 1}, {2, 2, 2, 2}, {3, 3, 3, 3},
			})
		})

		g.It("#insert does nothing if given nil", func() {
			var o *Obj
			var tree = NewRTree(4).LoadBoxes(data)
			g.Assert(tree.Data).Eql(NewRTree(4).LoadBoxes(data).Insert(o).Data)
		})

		g.It("#insert forms a valid tree if items are inserted one by one", func() {
			var tree = NewRTree(4)
			for i := 0; i < len(data); i++ {
				tree.Insert(Object( i,  data[i]))
			}

			var tree2 = NewRTree(4).LoadBoxes(data)
			g.Assert(tree.Data.height-tree2.Data.height <= 1).IsTrue()

			var boxes2 = make([]*mbr.MBR, 0)
			all2 := tree2.All()
			for i := 0; i < len(all2); i++ {
				boxes2 = append(boxes2, all2[i].bbox)
			}
			testResults(g, tree.All(), boxes2)
		})

		g.It("#remove removes items correctly", func() {
			var tree = NewRTree(4).LoadBoxes(data)
			var N = len(data)
			tree.RemoveMBR(data[0])
			tree.RemoveMBR(data[1])
			tree.RemoveMBR(data[2])

			tree.RemoveMBR(data[N-1])
			tree.RemoveMBR(data[N-2])
			tree.RemoveMBR(data[N-3])
			var cloneData []*mbr.MBR
			for i := 3; i < len(data)-3; i++ {
				cloneData = append(cloneData, data[i].Clone())
			}

			testResults(g, tree.All(), cloneData)

		})

		g.It("#remove does nothing if nothing found", func() {
			var item *Obj
			var tree = NewRTree(0).LoadBoxes(data)
			var tree2 = NewRTree(0).LoadBoxes(data)
			var query = mbr.New(13, 13, 13, 13)
			var querybox = Object(0, mbr.New(13, 13, 13, 13))
			g.Assert(tree.Data).Eql(tree2.RemoveMBR(query).Data)
			g.Assert(tree.Data).Eql(tree2.RemoveBoxObj(querybox).Data)
			g.Assert(tree.Data).Eql(tree2.RemoveBoxObj(item).Data)
			g.Assert(tree.Data).Eql(tree2.RemoveMBR(nil).Data)

		})

		g.It("#remove brings the tree to a clear state when removing everything one by one", func() {
			var tree = NewRTree(4).LoadBoxes(data)
			var result = tree.Search(mbr.New(0, 0, 100, 100))
			for i := 0; i < len(result); i++ {
				tree.RemoveNode(result[i])
			}
			g.Assert(tree.RemoveNode(&Node{}).IsEmpty()).IsTrue()
		})

		g.It("#clear should clear all the data in the tree", func() {
			var tree = NewRTree(4).LoadBoxes(data).Clear()
			g.Assert(tree.IsEmpty()).IsTrue()
		})

		g.It("should have chainable API", func() {
			g.Assert(NewRTree(4).LoadBoxes(
				data,
			).Insert(
				Object(0, data[0]),
			).RemoveMBR(data[0]).Clear().IsEmpty()).IsTrue()
		})
	})

}

/*
	g := goblin.Goblin(t)
	g.Describe("Rtree Tests - From Rbush", func() {
 */
func TestRtreeUtil(t *testing.T) {
	g := goblin.Goblin(t)
	g.Describe("Rtree Util", func() {
		g.It("tests pop nodes", func() {
			a := NewNode(Object( 0,  emptyMBR()), 0, true, nil)
			b := NewNode(Object( 0,  emptyMBR()), 1, true, nil)
			c := NewNode(Object( 0,  emptyMBR()), 1, true, nil)
			var nodes = make([]*Node, 0)
			var n *Node

			n, nodes = popNode(nodes)
			g.Assert(n == nil).IsTrue()

			nodes = []*Node{a, b, c}
			g.Assert(len(nodes)).Equal(3)

			n, nodes = popNode(nodes)
			g.Assert(len(nodes)).Equal(2)
			g.Assert(n).Eql(c)

			n, nodes = popNode(nodes)
			g.Assert(len(nodes)).Equal(1)
			g.Assert(n).Eql(b)

			n, nodes = popNode(nodes)
			g.Assert(len(nodes)).Equal(0)
			g.Assert(n).Eql(a)

			n, nodes = popNode(nodes)
			g.Assert(len(nodes)).Equal(0)
			g.Assert(n == nil).IsTrue()

			nodes = []*Node{a, b, c}
			g.Assert(len(nodes)).Equal(3)
			nodes = removeNode(nodes, 1)
			g.Assert(len(nodes)).Equal(2)
			nodes = removeNode(nodes, 4)
			g.Assert(len(nodes)).Equal(2)

		})
		g.It("tests pop index", func() {
			a := 0
			b := 1
			c := 2
			var indexes = make([]int, 0)
			var n int

			n = popIndex(&indexes)
			g.Assert(n == 0).IsTrue()

			indexes = []int{a, b, c}
			g.Assert(len(indexes)).Equal(3)

			n = popIndex(&indexes)
			g.Assert(len(indexes)).Equal(2)
			g.Assert(n).Eql(c)

			n = popIndex(&indexes)
			g.Assert(len(indexes)).Equal(1)
			g.Assert(n).Eql(b)

			n = popIndex(&indexes)
			g.Assert(len(indexes)).Equal(0)
			g.Assert(n).Eql(a)

			n = popIndex(&indexes)
			g.Assert(len(indexes)).Equal(0)
			g.Assert(n == 0).IsTrue()
		})
	})

}
