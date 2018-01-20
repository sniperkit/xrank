[![GoDoc](http://godoc.org/github.com/mozu0/wltree?status.png)](http://godoc.org/github.com/mozu0/wltree)
[![Build Status](https://travis-ci.org/mozu0/wltree.svg?branch=master)](https://travis-ci.org/mozu0/wltree)
# wltree
Go library of [Wavelet Tree](http://en.wikipedia.org/wiki/Wavelet_Tree) that supports Rank and Select.
Wavelet Tree is a index on bytestring (or int64 keys) `s`, and can return the number of specific character in any substring of `s` in constant time (, or more precisely, O(log of number of distinct alphabets)).

## Example
```go
s := []byte("abracadabra")
wt := wltree.NewBytes(s)
wt.Rank('a', len(s)) //=> 5 (The number of 'a' in s.)
wt.Rank('a', 8) - wt.Rank('a', 3) //=> 3 (The number of 'a' in s[3:8] = "acada") 
wt.Select('a', 2 /* 0-origin, thus means 3rd */) //=> 5 (The index of the 3rd occurrence of 'a' in s)
```
