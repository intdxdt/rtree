package rtree

import (
	"fmt"
	"github.com/intdxdt/mbr"
	"github.com/intdxdt/heap"
	"github.com/intdxdt/math"
)

func predicate(_ *KObj) (bool, bool) {
	return true, false
}

func (tree *RTree) KNN(
	query BoxObj, limit int, score func(BoxObj, BoxObj) float64,
	predicates ...func(*KObj) (bool, bool)) []BoxObj {

	var predFn = predicate
	if len(predicates) > 0 {
		predFn = predicates[0]
	}

	var node = tree.Data
	var result = make([]BoxObj, 0)
	var child *Node
	var queue = heap.NewHeap(kobj_cmp, heap.NewHeapType().AsMin())
	var stop, pred bool
	var dist float64

	for !stop && (node != nil) {
		for i := 0; i < len(node.children); i++ {
			child = node.children[i]

			if len(child.children) == 0 {
				dist = score(query, child.GetItem())
			} else {
				dist = score(query, child.BBox())
			}

			queue.Push(&KObj{
				node:   child,
				isItem: len(child.children) == 0,
				dist:   dist,
			})
		}

		for !queue.IsEmpty() && queue.Peek().(*KObj).isItem {
			var candidate = queue.Pop().(*KObj)
			pred, stop = predFn(candidate)
			if pred {
				result = append(result, candidate.GetItem())
			}

			if stop {
				break
			}

			if limit != 0 && len(result) == limit {
				return result
			}
		}

		if !stop {
			q := queue.Pop()
			if q == nil {
				node = nil
			} else {
				node = q.(*KObj).node
			}
		}
	}
	return result
}

//KObj instance struct
type KObj struct {
	node   *Node
	isItem bool
	dist   float64
}

//Score of knn object
func (kobj *KObj) Score() float64 {
	return kobj.dist
}

//IsItem check is knn object is leaf node item
func (kobj *KObj) IsItem() bool {
	return kobj.isItem
}

func (kobj *KObj) GetItem() BoxObj {
	return kobj.node.GetItem()
}

//BBox - satisfies BoxObj interface
func (kobj *KObj) BBox() *mbr.MBR {
	return kobj.node.BBox()
}

//String representation of knn object
func (kobj *KObj) String() string {
	return fmt.Sprintf("%v -> %v", kobj.node.bbox.String(), kobj.dist)
}

//Compare - cmp interface
func kobj_cmp(a interface{}, b interface{}) int {
	self, other := a.(*KObj), b.(*KObj)
	dx := self.dist - other.dist
	if math.FloatEqual(dx, 0.0) {
		return 0
	} else if dx < 0 {
		return -1
	}
	return 1
}
