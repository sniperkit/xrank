/*
Package wltree provides an implementation of Wavelet Tree.
See http://en.wikipedia.org/wiki/Wavelet_Tree for details.

Example

    s := []byte("abracadabra")
    wt := wltree.NewBytes(s)
    // The number of 'a' in s.
    wt.Rank('a', len(s)) //=> 5
    // The number of 'a' in s[3:8] = "acada"
    wt.Rank('a', 8) - wt.Rank('a', 3) //=> 3
    // The index of the 3rd occurrence of 'a' in s. 0-origin, thus 2 means 3rd.
    wt.Select('a', 2) //=> 5
*/
package wltree

import (
	"fmt"

	"github.com/mozu0/bitvector"
	"github.com/mozu0/huffman"
)

// Interface is a interface for arraylike elements that can be indexed by Wavelet Tree.
// Equal elements must have the same key, and different elements must have different keys.
type Interface interface {
	// Len returns the length of the arraylike.
	Len() int
	// Key returns the integer key of the i-th element.
	Key(i int) int64
}

// Int64Keys represents a Wavelet Tree on int64 keys.
type Int64Keys struct {
	nodes map[int64][]*bitvector.BitVector
	codes map[int64]string
}

// NewInt64Keys makes a Wavlet Tree from arraylike s whose elements can yield integer keys.
func NewInt64Keys(s Interface) *Int64Keys {
	w := &Int64Keys{
		nodes: make(map[int64][]*bitvector.BitVector),
		codes: make(map[int64]string),
	}

	// Generate huffman tree based on character occurrences in s.
	keyset, counts := freq(s)
	codes := huffman.FromInts(counts)
	for i, code := range codes {
		w.codes[keyset[i]] = code
	}

	// Count number of bits in each node of the wavelet tree.
	sizes := make(map[string]int)
	for i := range keyset {
		code := codes[i]
		count := counts[i]
		for j := range code {
			sizes[code[:j]] += count
		}
	}

	// Assign BitVector Builders to each wavelet tree node.
	builders := make(map[string]*bitvector.Builder)
	var keys []string
	for key, size := range sizes {
		keys = append(keys, key)
		builders[key] = bitvector.NewBuilder(size)
	}

	// Set bits in each BitVector Builder.
	index := make(map[string]int)
	for i, size := 0, s.Len(); i < size; i++ {
		k := s.Key(i)
		code := w.codes[k]
		for j := range code {
			if code[j] == '1' {
				builders[code[:j]].Set(index[code[:j]])
			}
			index[code[:j]]++
		}
	}

	// Build all BitVectors.
	bvs := make(map[string]*bitvector.BitVector)
	for key, builder := range builders {
		bvs[key] = builder.Build()
	}

	// For each charactor, register the path from wavelet tree root, through wavelet tree nodes, and
	// to the leaf.
	for i, k := range keyset {
		code := codes[i]
		for j := range code {
			w.nodes[k] = append(w.nodes[k], bvs[code[:j]])
		}
	}

	return w
}

// Rank returns the count of elements with the key in s[0:i].
func (w *Int64Keys) Rank(key int64, i int) int {
	code := w.codes[key]
	if code == "" {
		return 0
	}

	nodes := w.nodes[key]
	for j := range nodes {
		if code[j] == '1' {
			i = nodes[j].Rank1(i)
		} else {
			i = nodes[j].Rank0(i)
		}
	}
	return i
}

// Select returns i such that Rank(c, i) = r.
// i.e. it returns the index of r-th occurrence of the element with the key.
// Note that r is 0-origined, so wt.Select('a', 2) returns the index of the third 'a'.
func (w *Int64Keys) Select(key int64, r int) int {
	code := w.codes[key]
	if code == "" {
		panic(fmt.Sprintf("wltree: no such element with key %v in s.", key))
	}

	nodes := w.nodes[key]
	for j := len(nodes) - 1; j >= 0; j-- {
		if code[j] == '1' {
			r = nodes[j].Select1(r)
		} else {
			r = nodes[j].Select0(r)
		}
	}
	return r
}

// Bytes represents a Wavelet Tree on bytestring.
type Bytes struct {
	nodes [256][]*bitvector.BitVector
	codes [256]string
}

// NewBytes constructs a Wavelet Tree from bytestring.
func NewBytes(s []byte) *Bytes {
	intKeys := NewInt64Keys(byteSlice(s))
	b := &Bytes{}
	for i, nodes := range intKeys.nodes {
		b.nodes[i] = nodes
	}
	for i, code := range intKeys.codes {
		b.codes[i] = code
	}
	return b
}

// Rank returns the count of the character c in s[0:i].
func (w *Bytes) Rank(c byte, i int) int {
	code := w.codes[c]
	if code == "" {
		return 0
	}

	nodes := w.nodes[c]
	for j := range nodes {
		if code[j] == '1' {
			i = nodes[j].Rank1(i)
		} else {
			i = nodes[j].Rank0(i)
		}
	}
	return i
}

// Select returns i such that Rank(c, i) = r.
// i.e. it returns the index of r-th occurrence of the character c.
// Note that r is 0-origined, so wt.Select('a', 2) returns the index of the third 'a'.
func (w *Bytes) Select(c byte, r int) int {
	code := w.codes[c]
	if code == "" {
		panic(fmt.Sprintf("wltree: no such character %q in s.", string(c)))
	}

	nodes := w.nodes[c]
	for j := len(nodes) - 1; j >= 0; j-- {
		if code[j] == '1' {
			r = nodes[j].Select1(r)
		} else {
			r = nodes[j].Select0(r)
		}
	}
	return r
}

func freq(s Interface) (keyset []int64, counts []int) {
	freqs := make(map[int64]int)
	for i, size := 0, s.Len(); i < size; i++ {
		freqs[s.Key(i)]++
	}
	for k, w := range freqs {
		keyset = append(keyset, k)
		counts = append(counts, w)
	}
	return
}

type byteSlice []byte

func (b byteSlice) Len() int {
	return len(b)
}
func (b byteSlice) Key(i int) int64 {
	return int64(b[i])
}
