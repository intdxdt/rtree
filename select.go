package rtree

import "math"

// sort an array so that items come in groups of n unsorted items,
// with groups sorted between each other and
// combines selection algorithm with binary divide & conquer approach.
func multiSelect(arr []BoxObj, left, right, n int, compare compareNode) {
    stack := make([]int, 2)
    stack[0], stack[1] = left, right
    var mid int

    for len(stack) > 0 {
        right, stack = popInt(stack)
        left, stack = popInt(stack)

        if right - left <= n {
            continue
        }

        _mid := float64(right - left) / float64(n)
        mid = left + int(math.Ceil(_mid / 2.0)) * n

        selectBox(arr, left, right, mid, compare)
        stack = appendInts(stack, left, mid, mid, right)
    }
}


// sort array between left and right (inclusive) so that the smallest k elements come first (unordered)
func selectBox(arr []BoxObj, left, right, k int, compare compareNode) {
    var i, j int
    var _n, _i, _newLeft, _newRight, _left, _right, _sn float64
    var _z, _s, _sd float64
    _left, _right, _k := float64(left), float64(right), float64(k)
    var t BoxObj

    for right > left {
        if right - left > 500 {
            _n = _right - _left + 1.0
            _i = _k - _left + 1.0
            _z = math.Log(_n)
            _s = 0.5 * math.Exp(2 * _z / 3.0)
            _sn = 1.0
            if (_i - _n / 2) < 0 {
                _sn = -1.0
            }
            _sd = 0.5 * math.Sqrt(_z * _s * (_n - _s) / _n) * (_sn)
            _newLeft = max(_left, math.Floor(_k - _i * _s / _n + _sd))
            _newRight = min(_right, math.Floor(_k + (_n - _i) * _s / _n + _sd))
            selectBox(arr, int(_newLeft), int(_newRight), int(_k), compare)
        }

        t = arr[k]
        i = left
        j = right

        swapItem(arr, left, k)
        if compare(&arr[right], &t) > 0 {
            swapItem(arr, left, right)
        }

        for i < j {
            swapItem(arr, i, j)
            i++
            j--
            for compare(&arr[i], &t) < 0 {
                i++
            }
            for compare(&arr[j], &t) > 0 {
                j--
            }
        }

        if compare(&arr[left], &t) == 0 {
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

