package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"math"
	"os"
	"runtime"
	"strconv"
	"time"
)

func check(err error) {
	if err != nil && err != io.EOF {
		panic(err)
	}
}

var filename = flag.String("file", "/Users/AnimotoOverstreet/go/bin/web-Google.txt", "file to parse on")

func scanInput(scanner *bufio.Scanner, nodeCount int) (adjacencyMatrix, map[int]int) {
	var matrix adjacencyMatrix
	matrix.initWithSize(nodeCount)

	indexMap := make(map[int]int)
	nextIndex := 0

	for scanner.Scan() {
		outIndex, err := strconv.Atoi(scanner.Text())
		check(err)
		if !scanner.Scan() {
			panic("unexpected failed scan")
		}
		inIndex, err := strconv.Atoi(scanner.Text())
		check(err)

		out, ok := indexMap[outIndex]
		if !ok {
			out = nextIndex
			indexMap[outIndex] = nextIndex
			nextIndex++
		}

		in, ok := indexMap[inIndex]
		if !ok {
			in = nextIndex
			indexMap[inIndex] = nextIndex
			nextIndex++
		}

		matrix.addEdge(out, in)
	}

	return matrix, indexMap
}

func setVals(slice []float64, val float64) {
	for i := range slice {
		slice[i] = val
	}
}

func sumVals(slice []float64) float64 {
	sum := 0.0
	for _, v := range slice {
		sum += v
	}
	return sum
}

func addToVals(slice []float64, val float64) {
	for i := range slice {
		slice[i] += val
	}
}

// file format rules:
// first line is an integer for the number of nodes
// all following lines are a pain of integers
// the first int is the outIndex, second in the inIndex
func main() {
	flag.Parse()
	fmt.Println(*filename)

	start := time.Now()

	file, err := os.Open(*filename)
	check(err)

	scanner := bufio.NewScanner(bufio.NewReader(file))
	scanner.Split(bufio.ScanWords)

	scanner.Scan()
	nodeCount, err := strconv.Atoi(scanner.Text())
	check(err)

	matrix, indexMap := scanInput(scanner, nodeCount)

	if scanner.Err() != nil {
		panic("reading input exited with non-EOF error")
	}

	file.Close()

	finish := time.Since(start)
	//reading file done
	fmt.Println("time to read file:", finish)

	start = time.Now()

	iteration := 0
	const epsilon = 0.0000000001
	const beta = 0.8
	rank := make([]float64, nodeCount)
	rankNew := make([]float64, nodeCount)
	done := false

	concurrencyFactor := runtime.GOMAXPROCS(0) * 64

	sectionSize := (nodeCount + concurrencyFactor - 1) / concurrencyFactor

	setVals(rank, 1.0/float64(nodeCount))

	for !done {
		matrix.generateNewRank(rank, rankNew, beta, concurrencyFactor)

		sum := sumVals(rankNew)

		jumpFactor := (1.0 - sum) / float64(nodeCount)

		addToVals(rankNew, jumpFactor)

		dones := make(chan bool)

		fail := make(chan bool)

		done = true

		for i := 0; i < concurrencyFactor; i++ {
			go func(i int) {
				bottom := sectionSize * i
				top := sectionSize * (i + 1)
				if i == concurrencyFactor-1 {
					top = nodeCount
				}
				for index := range rankNew[bottom:top] {
					if math.Abs(rankNew[index+bottom]-rank[index+bottom]) > epsilon {
						dones <- false
						fail <- true
						return
					}
					select {
					case <-fail:
						return
					default:
					}
				}
				dones <- true
			}(i)
		}

		for i := 0; i < concurrencyFactor; i++ {
			ok := <-dones
			if !ok {
				done = false
				break
			}
		}

		copy(rank, rankNew)
		iteration++
	}

	finish = time.Since(start)

	fmt.Println("time to run algorithm:", finish)

	fmt.Println("rank of node 99:", rank[indexMap[99]])
	fmt.Println("final iteration:", iteration)
}
