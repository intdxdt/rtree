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

	for !stop && (node != nil) {
		for i := 0; i < len(node.children); i++ {
			child = node.children[i]
			var o = &KObj{
				child,
				child.bbox,
				len(child.children) == 0,
				-1,
			}
			o.Dist = score(query, o)
			queue.Push(o)
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
			var q = queue.Pop()
			if q == nil {
				node = nil
			} else {
				node = q.(*KObj).node
			}
		}
	}
	return result
}
