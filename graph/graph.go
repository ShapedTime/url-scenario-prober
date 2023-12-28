package graph

type Graph[T ConcurrencySafeObj] struct {
	obj   *T
	nodes []*Graph[T]
}

type ConcurrencySafeObj interface {
	RLock()
	RUnlock()
	Lock()
	Unlock()
}

func NewGraph[T ConcurrencySafeObj](obj *T) *Graph[T] {
	return &Graph[T]{
		obj:   obj,
		nodes: []*Graph[T]{},
	}
}

func (g *Graph[T]) AddNode(node *Graph[T]) {
	g.nodes = append(g.nodes, node)
}

func (g *Graph[T]) GetNodes() []*Graph[T] {
	return g.nodes
}

func (g *Graph[T]) GetObj() *T {
	return g.obj
}

func (g *Graph[T]) SetObj(obj *T) {
	g.obj = obj
}

func (g *Graph[T]) BFSFirst(f func(obj T) bool) (*T, bool) {
	for i, _ := range g.nodes {
		(*g.nodes[i].obj).RLock()
		obj := *g.nodes[i].obj
		(*g.nodes[i].obj).RUnlock()

		if f(obj) {
			return g.nodes[i].obj, true
		}

		// If node has children, search them too
		if len(g.nodes[i].nodes) > 0 {
			return g.nodes[i].BFSFirst(f)
		}
	}

	return nil, false
}
