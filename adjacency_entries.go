package main

type adjacencyEntries struct {
	inEdges      []int
	outEdgeCount int
}

func (a *adjacencyEntries) addInEdge(outIndex int) {
	a.inEdges = append(a.inEdges, outIndex)
}

func (a *adjacencyEntries) addOutEdge() {
	a.outEdgeCount++
}
