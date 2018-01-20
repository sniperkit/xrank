package wltree

import (
	"math/rand"
	"testing"
)

const maxSize = 512

var weights = []map[byte]int{{
	'a': 1,
	'c': 1,
	'g': 1,
	't': 1,
}, {
	'a': 1,
	'b': 1,
	'c': 2,
	'd': 3,
	'e': 5,
	'f': 8,
}}

func TestWltree(t *testing.T) {
	fails := 0

	for size := 0; size < maxSize; size++ {
		for _, ws := range weights {
			bs := random(size, ws)
			wt := NewBytes(bs)
			wti := NewInt64Keys(byteSlice(bs))

			var counts [256]int
			for i := 0; i <= len(bs) && fails < 30; i++ {
				for c := 0; c < 256; c++ {
					c := byte(c)
					if got, want := wt.Rank(c, i), counts[c]; got != want {
						t.Errorf("Bytes: %q.Rank(%v, %v) => got %v, want %v", bs, string(c), i, got, want)
						fails++
					}
					if got, want := wti.Rank(int64(c), i), counts[c]; got != want {
						t.Errorf("IntKeys: %q.Rank(%v, %v) => got %v, want %v", bs, string(c), i, got, want)
						fails++
					}
				}
				if i != len(bs) {
					c := bs[i]
					if got, want := wt.Select(c, counts[c]), i; got != want {
						t.Errorf("Bytes: %q.Select(%v, %v) => got %v, want %v", bs, string(c), counts[c], got, want)
						fails++
					}
					if got, want := wti.Select(int64(c), counts[c]), i; got != want {
						t.Errorf("IntKeys: %q.Select(%v, %v) => got %v, want %v", bs, string(c), counts[c], got, want)
						fails++
					}
					counts[bs[i]]++
				}
			}
		}
	}
}

func random(size int, weights map[byte]int) []byte {
	var (
		cs  []byte
		ws  []int
		sum = 0
	)
	for c, w := range weights {
		cs = append(cs, c)
		ws = append(ws, w)
		sum += w
	}

	var bs []byte
	for i := 0; i < size; i++ {
		r := rand.Intn(sum)
		for j, w := range ws {
			r -= w
			if r < 0 {
				bs = append(bs, cs[j])
			}
		}
	}

	return bs
}
