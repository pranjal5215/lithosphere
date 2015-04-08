package lithosphere

type Tree struct {
	M MergeFunc
	L []LeafFunc
}

type LeafFunc func()
type MergeFunc func(size int)

func (t *Tree) DoTree(size int) {
	for _, f := range t.L {
		go f() //all f's should push to channel
	}
	t.M(size) //M should listen to channel and process
}
