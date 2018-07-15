package rtree

import "math"

// sort an array so that items come in groups of n unsorted items,
// with groups sorted between each other and
// combines selection algorithm with binary divide & conquer approach.
func multiSelect(arr []*Obj, left, right, n int, compare compareNode) {
	var mid int
	var stack = make([]int, 2)
	stack[0], stack[1] = left, right

	for len(stack) > 0 {
		right, stack = popInt(stack)
		left, stack = popInt(stack)

		if right-left <= n {
			continue
		}

		mid = left + int(math.Ceil(float64(right-left)/float64(n)/2.0))*n
		selectBox(arr, left, right, mid, compare)
		stack = appendInts(stack, left, mid, mid, right)
	}
}

// sort array between left and right (inclusive) so that the smallest k elements come first (unordered)
func selectBox(arr []*Obj, left, right, k int, compare compareNode) {
	var i, j int
	var fn, fi, fNewLeft, fNewRight, fsn, fz, fs, fsd float64
	var fLeft, fRight, fk = float64(left), float64(right), float64(k)
	var t *Obj

	for right > left {
		//the arbitrary constants 600 and 0.5 are used in the original
		// version to minimize execution time
		if right-left > 600 {
			fn = fRight - fLeft + 1.0
			fi = fk - fLeft + 1.0
			fz = math.Log(fn)
			fs = 0.5 * math.Exp(2*fz/3.0)
			fsn = 1.0
			if (fi - fn/2) < 0 {
				fsn = -1.0
			}
			fsd = 0.5 * math.Sqrt(fz*fs*(fn-fs)/fn) * (fsn)
			fNewLeft = max(fLeft, math.Floor(fk-fi*fs/fn+fsd))
			fNewRight = min(fRight, math.Floor(fk+(fn-fi)*fs/fn+fsd))
			selectBox(arr, int(fNewLeft), int(fNewRight), int(fk), compare)
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
