package lib

import (
	"github.com/hishboy/gocommons/lang"
)

// store git repo message
var Queue *lang.Queue

// define the size of concurrent goroutine
var ThreadPool *lang.Queue

// record the cloned git repository
// used for concurrent count word
var RepoQueue *lang.Queue

func init() {
	Queue = lang.NewQueue()
	RepoQueue = lang.NewQueue()
	initTheadPool(10)
}

func initTheadPool(size int) {
	ThreadPool = lang.NewQueue()
	for i := 0; i < size; i++ {
		ThreadPool.Push(i)
	}
}
