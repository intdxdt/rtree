package rtree

import (
	"fmt"
	"github.com/intdxdt/mbr"
)

//KObj instance struct
type KObj struct {
	node   *node
	MBR    *mbr.MBR
	IsItem bool
	Dist   float64
}

func (kobj *KObj) GetItem() BoxObj {
	return kobj.node.item
}

//String representation of knn object
func (kobj *KObj) String() string {
	return fmt.Sprintf("%v -> %v", kobj.node.bbox.String(), kobj.Dist)
}

//Compare - cmp interface
func kobjCmp(a interface{}, b interface{}) int {
	var self, other = a.(*KObj), b.(*KObj)
	var dx = self.Dist - other.Dist
	var r = 1
	if feq(dx, 0) {
		r = 0
	} else if dx < 0 {
		r = -1
	}
	return r
}
