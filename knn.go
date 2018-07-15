package rtree

import (
	"github.com/intdxdt/heap"
	"github.com/intdxdt/mbr"
)

func predicate(_ *KObj) (bool, bool) {
	return true, false
}

func (tree *RTree) Knn(
	query *mbr.MBR, limit int, score func(*mbr.MBR, *KObj) float64,
	predicates ...func(*KObj) (bool, bool)) []*Obj {

	var predFn = predicate
	if len(predicates) > 0 {
		predFn = predicates[0]
	}

	var node = tree.Data
	var result []*Obj
	var child *rNode
	var queue = heap.NewHeap(kObjCmp, heap.NewHeapType().AsMin())
	var stop, pred bool
	var dist float64

	for !stop && (node != nil) {
		for i := 0; i < len(node.children); i++ {
			child = node.children[i]

			if len(child.children) == 0 {
				dist = score(query, &KObj{child, child.bbox, true, -1})
			} else {
				dist = score(query, &KObj{nil, child.bbox, false, -1})
			}

			queue.Push(&KObj{
				Node:   child,
				MBR:    child.bbox,
				IsItem: len(child.children) == 0,
				Dist:   dist,
			})
		}

		for !queue.IsEmpty() && queue.Peek().(*KObj).IsItem {
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
				node = q.(*KObj).Node
			}
		}
	}
	return result
}
