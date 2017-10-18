package rtree

func nodeAtIndex(a []*Node, i int) *Node {
    var n = len(a)
    if i > n - 1 || i < 0 || n == 0 {
        return nil
    }
    return a[i]
}

func popNode(a []*Node) (*Node, []*Node) {
    var v *Node
    var n int
    if len(a) == 0 {
        return nil, a
    }
    n = len(a) - 1
    v, a[n] = a[n], nil
    return v, a[:n]
}

func popIndex(indxs *[]int) int {
    var n, v int
    a := *indxs
    n = len(a) - 1
    if n < 0 {
        return 0
    }
    v, a[n] = a[n], 0
    *indxs = a[:n]
    return v
}

//remove node at given index from node slice.
func removeNode(a []*Node, i int) []*Node {
    n := len(a) - 1
    if i > n {
        return a
    }
    a, a[n] = append(a[:i], a[i + 1:]...), nil
    return a
}

//condense node and its path from the root
func (tree *RTree) condense(path []*Node) {
    var parent *Node
    // go through the path, removing empty nodes and updating bboxes
    for i := len(path) - 1; i >= 0; i-- {
        if len(path[i].children) == 0 {
            //go through the path, removing empty nodes and updating bboxes
            if i > 0 {
                parent = path[i - 1]
                index := sliceIndex(len(parent.children), func(s int) bool {
                    return path[i] == parent.children[s]
                })
                parent.children = removeNode(parent.children, index)

            } else {
                //root is empty, rest root
                tree.Clear()
            }
        } else {
            calcBBox(path[i])
        }
    }
}


//Remove Item from RTree
//NOTE:if item is a bbox , then first found bbox is removed
func (tree *RTree) Remove(item BoxObj) *RTree {
    if item == nil {
        return tree
    }
    _, isnode := item.(*Node)

    var node *Node = tree.Data
    var parent *Node
    var bbox = item.BBox()
    var path = make([]*Node, 0)
    var indexes = make([]int, 0)
    var i, index int
    var goingUp bool

    // depth-first iterative self traversal
    for (node != nil ) || (len(path) > 0) {
        if node == nil {
            // go up
            node, path = popNode(path)
            parent = nodeAtIndex(path, len(path) - 1)
            i = popIndex(&indexes)
            goingUp = true
        }

        if node.leaf {
            // check current node
            //index = node.children.indexOf(item)
            index = sliceIndex(len(node.children), func(i int) bool {
                if isnode {
                    return node.children[i] == item
                }
                return node.children[i].bbox.Equals(item.BBox())
            })

            //if found
            if index != -1 {
                //item found, remove the item and condense self upwards
                //node.children.splice(index, 1)
                node.children = removeNode(node.children, index)
                path = append(path, node)
                tree.condense(path)
                return tree
            }

        }

        if !goingUp && !node.leaf && contains(node.bbox, bbox) {
            // go down
            path = append(path, node)
            indexes = append(indexes, i)
            i = 0
            parent = node
            node = node.children[0]
        } else if parent != nil {
            // go right
            i++
            node = nodeAtIndex(parent.children, i);
            goingUp = false
        } else {
            node = nil
        } // nothing found
    }
    return tree
}


