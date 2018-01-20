package rtree

import (
	"fmt"
	"github.com/intdxdt/mbr"
	"github.com/intdxdt/math"
)

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
func kobjCmp(a interface{}, b interface{}) int {
	self, other := a.(*KObj), b.(*KObj)
	dx := self.dist - other.dist
	if math.FloatEqual(dx, 0.0) {
		return 0
	} else if dx < 0 {
		return -1
	}
	return 1
}
