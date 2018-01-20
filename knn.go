package rtree

import (
	"github.com/intdxdt/heap"
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
	var queue = heap.NewHeap(kobjCmp, heap.NewHeapType().AsMin())
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
