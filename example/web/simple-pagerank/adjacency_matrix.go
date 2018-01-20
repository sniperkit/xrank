package main

type adjacencyMatrix []*adjacencyEntries

func (a *adjacencyMatrix) initWithSize(size int) {
	(*a) = adjacencyMatrix(make([]*adjacencyEntries, size))
	for i := range *a {
		(*a)[i] = &adjacencyEntries{make([]int, 0, 10), 0}
	}
}

func (a adjacencyMatrix) addEdge(out int, in int) {
	a[out].addOutEdge()
	a[in].addInEdge(out)
}

func (a adjacencyMatrix) size() int {
	return len(a)
}

//x is the concurrency factor
//note: will error if number of nodes is too small relative to x
func (a adjacencyMatrix) generateNewRank(rank []float64, rankNew []float64, beta float64, x int) {
	dones := make(chan bool)

	sectionSize := (a.size() + x - 1) / x

	for i := 0; i < x; i++ {
		go func(i int) {
			bottom := sectionSize * i
			top := sectionSize * (i + 1)
			if i == x-1 {
				top = a.size()
			}
			for index, entry := range a[bottom:top] {
				r := 0.0
				for _, inEdge := range entry.inEdges {
					r += beta * rank[inEdge] / float64(a[inEdge].outEdgeCount)
				}
				rankNew[index+bottom] = r
			}
			dones <- true
		}(i)
	}

	for i := 0; i < x; i++ {
		<-dones
	}
}
