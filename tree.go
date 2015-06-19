package lithosphere

type Tree struct {
	M    MergeFunc
	L    []LeafFunc
	Linp []interface{}
}

type LeafFunc func(inp interface{})
type MergeFunc func(size int)

func (t *Tree) DoTree(size int) {
	for k, f := range t.L {
		go f(t.Linp[k]) //all f's should push to channel
	}
	t.M(size) //M should listen to channel and process
}
