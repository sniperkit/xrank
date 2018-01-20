package webgraph

import "errors"

// Key: used for conversion between key and index into an array
// Incoming: key of the source node where the arrow comes from and it's weight
type Vertex struct {
	Key      string
	Incoming map[string]float64
}

// To make a graph first we gradually fill Keys by using AddNode. Keys hashes
// strings and gives us integer indexes. Then when we have all of them Fixate
// allocates Nodes. Each node is a Vertex which knows its original Key and all
// its neighbourly Incoming nodes.
// After Fixate we use AddArrow to add directed connections. The Incoming part
// of each node is like an adjency list. Elements of OutgoingCount will be
// incremented when an outgoing arrow is added by AddArrow.
// CalculateDefaultWeights will then fixate the values in OutgoingCount.
type Graph struct {
	Keys          map[string]int
	Nodes         []Vertex
	OutgoingCount []int
}

// New returns a webgraph. Use expectedNodeCount to hint at the total no. of
// nodes it will hold.
func New(expectedNodeCount int) (g *Graph) {
	g = new(Graph)
	g.Keys = make(map[string]int, expectedNodeCount)
	return
}

// AddNode can only be used before Fixate. It checks if it already has a
// node with key if so it returns an error, if not it adds it.
func (g *Graph) AddNode(key string) error {
	if _, present := g.Keys[key]; present {
		return errors.New("AddNode: Node already present")
	}
	g.Keys[key] = len(g.Keys)
	return nil
}

// Fixate uses the Keys map to learn the length necessary for Nodes and
// fills in the Key part of each Vertex in the Nodes slice
func (g *Graph) Fixate() {
	l := len(g.Keys)
	g.OutgoingCount = make([]int, l)
	g.Nodes = make([]Vertex, l)
	for key, index := range g.Keys {
		g.Nodes[index].Key = key
		g.Nodes[index].Incoming = make(map[string]float64)
	}
}

// Returns the number of nodes added before the call to Fixate
func (g *Graph) FixedLength() int {
	return len(g.Nodes)
}

// After Fixate one can use this function to translate string keys to indices
func (g *Graph) Key2idx(key string) (idx int) {
	idx = g.Keys[key]
	return
}

// After Fixate one can use this function to translate indices to string keys
func (g *Graph) Idx2key(idx int) (key string) {
	key = g.Nodes[idx].Key
	return
}

// AddArrow adds an arrow from 'k2' to 'k1' if they are both present nodes
// and if the arrow doesn't already exist
func (g *Graph) AddArrow(k2, k1 string) error {
	k1Idx, k1Present := g.Keys[k1]
	if !k1Present {
		return errors.New("AddArrow: Could not find 'k1' node")
	}
	k2Idx, k2Present := g.Keys[k2]
	if !k2Present {
		return errors.New("AddArrow: Could not find 'k2' node")
	}
	if _, arrowPresent := g.Nodes[k1Idx].Incoming[k2]; arrowPresent {
		return errors.New("AddArrow: Arrow already in place")
	}
	g.Nodes[k1Idx].Incoming[k2] = 0.0
	g.OutgoingCount[k2Idx]++
	return nil
}

// In a default webgraph the PageRank is equally devided over the outgoing
// links. So one of n outgoing arrows has weight 1/n.
func (g *Graph) CalculateDefaultWeights() error {
	tmp := make([]int, len(g.OutgoingCount)) // we need a copy of OutgoingCount
	if len(g.OutgoingCount) != copy(tmp, g.OutgoingCount) {
		return errors.New("CalculateDefaultWeights: Could not copy all of OutgoingCount")
	}

	// distribute PageRank to all nodes if there are no outgoing arrows
	for j := range tmp {
		if tmp[j] == 0 {
			// and add arrows to all nodes from g.Idx2key(j)
			from := g.Idx2key(j)
			for to := range g.Keys {
				g.AddArrow(from, to)
			}
		}
	}
	// set the weight
	for i := range g.Nodes {
		for k := range g.Nodes[i].Incoming {
			g.Nodes[i].Incoming[k] = 1.0 / (float64(g.OutgoingCount[g.Key2idx(k)]))
		}
	}
	return nil
}

// takes a weight w and dampens it
func (g *Graph) Dampen(w float64) float64 {
	const d float64 = 0.85
	e := (1.0 - d) / (float64(g.FixedLength()))
	return d*w + e
}

// IncomingWeightsVector gives a row of the Google matrix.
// Takes a nodekey 'key' and finds the corresponding incoming node weights,
// which it puts in the correct positions of an otherwise empty vector.
// Gives an error if 'key' has no node associated with it.
func (g *Graph) IncomingWeightsVector(key string) (wvec []float64, err error) {
	if idx, present := g.Keys[key]; present {
		wvec = make([]float64, g.FixedLength())
		for k, w := range g.Nodes[idx].Incoming {
			wvec[g.Keys[k]] = w
		}
		for i := range wvec {
			wvec[i] = g.Dampen(wvec[i])
		}
	} else {
		err = errors.New("IncomingWeightsVector: Could not find node with given key")
	}
	return
}

// Multiply treats the webgraph as a square matrix which will be right multiplied
// by a vector 'invec'. It returns an error when invec doesn't correspond to
// the size of the square matrix (g.FixedLength()).
func (g *Graph) Multiply(invec []float64) (outvec []float64, err error) {
	if len(invec) != g.FixedLength() {
		err = errors.New("Multiply: Input vector not the right size")
		return
	}

	outvec = make([]float64, g.FixedLength())
	for i := 0; i < len(outvec); i++ {
		weights, werr := g.IncomingWeightsVector(g.Idx2key(i)) // row of the matrix
		if werr != nil {
			outvec = nil
			err = werr
			return
		}

		outvec[i] = 0.0
		for j := 0; j < len(invec); j++ {
			outvec[i] += invec[j] * weights[j]
		}
	}
	return
}
