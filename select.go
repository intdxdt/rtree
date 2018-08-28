package rtree

import "math"

// sort an array so that items come in groups of n unsorted items,
// with groups sorted between each other and
// combines selection algorithm with binary divide & conquer approach.
func multiSelect(arr []BoxObject, left, right, n int, compare compareObject) {
	var mid int
	var stack =[]int{left, right}

	for len(stack) > 0 {
		right, stack = popInt(stack)
		left, stack = popInt(stack)

		if right-left <= n {
			continue
		}

		mid = left + int(math.Ceil(float64(right-left)/float64(n)/2.0))*n
		selectBox(arr, mid, left, right,  compare)
		stack = appendInts(stack, left, mid, mid, right)
	}
}

// sort array between left and right (inclusive)
// so that the smallest k elements come first (unordered)
func selectBox(arr []BoxObject, k , left, right int, compare compareObject) {
	var i, j int
	var newLeft, newRight int
	var fn, fi,  fsn, fz, fs, fsd float64
	var fleft, fright, fk = float64(left), float64(right), float64(k)
	var t BoxObject

	for right > left {
		if right-left > 600 {
			fn = fright - fleft + 1
			fi = fk - fleft + 1
			fz = math.Log(fn)

			fs = 0.5 * math.Exp(2*fz/3.0)
			fsn = 1
			if (fi - fn/2) < 0 {
				fsn = -1
			}
			fsd = 0.5 * math.Sqrt(fz*fs*(fn-fs)/fn) * fsn
			newLeft  = int(max(fleft, math.Floor(fk-fi*fs/fn+fsd)))
			newRight = int(min(fright, math.Floor(fk+(fn-fi)*fs/fn+fsd)))
			selectBox(arr, k,  newLeft, newRight, compare)
		}

		t = arr[k]
		i = left
		j = right

		swapItem(arr, left, k)
		if compare(arr[right], t) > 0 {
			swapItem(arr, left, right)
		}

		for i < j {
			swapItem(arr, i, j)
			i++
			j--
			for compare(arr[i], t) < 0 {
				i++
			}
			for compare(arr[j], t) > 0 {
				j--
			}
		}

		if compare(arr[left], t) == 0 {
			swapItem(arr, left, j)
		} else {
			j++
			swapItem(arr, j, right)
		}

		if j <= k {
			left = j + 1
		}
		if k <= j {
			right = j - 1
		}
	}
}
