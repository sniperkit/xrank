package webgraph

import "testing"

func TestSimpleGraph(t *testing.T) {
	const size = 4
	g := New(size)

	g.AddNode("A")
	g.AddNode("B")
	g.AddNode("C")
	g.AddNode("D")
	err0 := g.AddNode("D")
	if err0 == nil {
		t.Fatal("AddNode allows adding duplicates")
	}

	g.Fixate()
	if g.FixedLength() != size {
		t.Fatal("size mismatch")
	}

	if g.Idx2key(g.Key2idx("A")) != "A" {
		t.Fatal("key to index translation or vice versa doesn't work")
	}
	if g.Idx2key(g.Key2idx("B")) != "B" {
		t.Fatal("key to index translation or vice versa doesn't work")
	}
	if g.Idx2key(g.Key2idx("C")) != "C" {
		t.Fatal("key to index translation or vice versa doesn't work")
	}
	if g.Idx2key(g.Key2idx("D")) != "D" {
		t.Fatal("key to index translation or vice versa doesn't work")
	}

	err1 := g.AddArrow("A", "NonExistent")
	if err1 == nil {
		t.Fatal("AddArrow allows for non-existent node arguments")
	}

	err2 := g.AddArrow("A", "B")
	if err2 != nil {
		t.Fatal("AddArrow is to strict, gives error were it shouldn't")
	}
	g.AddArrow("A", "C")
	g.AddArrow("D", "B")

	g.CalculateDefaultWeights()

	if v, e := g.IncomingWeightsVector("A"); e != nil || v[0] != 0.037500000000000006 || v[1] != 0.25 || v[2] != 0.25 || v[3] != 0.037500000000000006 {
		t.Logf("IWV: %v\n", v)
		t.Fatal("IncomingWeightsVector gives unexpected values")
	}
	if v, e := g.IncomingWeightsVector("B"); e != nil || v[0] != 0.4625 || v[1] != 0.25 || v[2] != 0.25 || v[3] != 0.8875 {
		t.Logf("IWV: %v\n", v)
		t.Fatal("IncomingWeightsVector gives unexpected values")
	}
	if v, e := g.IncomingWeightsVector("C"); e != nil || v[0] != 0.4625 || v[1] != 0.25 || v[2] != 0.25 || v[3] != 0.037500000000000006 {
		t.Logf("IWV: %v\n", v)
		t.Fatal("IncomingWeightsVector gives unexpected values")
	}
	if v, e := g.IncomingWeightsVector("D"); e != nil || v[0] != 0.037500000000000006 || v[1] != 0.25 || v[2] != 0.25 || v[3] != 0.037500000000000006 {
		t.Logf("IWV: %v\n", v)
		t.Fatal("IncomingWeightsVector gives unexpected values")
	}
	if v, e := g.IncomingWeightsVector("NonExistent"); v != nil || e == nil {
		t.Fatal("IncomingWeightsVector allows for non-existent nodes")
	}

	badinvec := make([]float64, g.FixedLength()+1)
	if _, e := g.Multiply(badinvec); e == nil {
		t.Fatal("Multiply allows multiplication with vectors of the wrong size")
	}

	invec := []float64{0.2, 0.3, 0.25, 0.25}
	outvec, merr := g.Multiply(invec)
	if merr != nil {
		t.Fatal("Multiply to restrictive, error on a good slice")
	}
	if outvec[0] != 0.154375 || outvec[1] != 0.451875 || outvec[2] != 0.239375 || outvec[3] != 0.154375 {
		t.Logf("outvec: %v\n", outvec)
		t.Fatal("Multiply gives unexpected values")
	}
}
