package main

import (
	"fmt"
	"github.com/wangming1993/verb-rank/lib"
)

func init() {
	fmt.Println("VERB-RANK INIT...")
}

func main() {
	golang()

	//wc()

	//redis()
}

func golang() {
	repo := &lib.Repo{
		Language: "go",
		Limit:    3,
	}
	repo.Search()
}

func wc() {
	dir := lib.CLONE_PATH
	lib.ReadDir(dir)
}

func redis() {
	reply, err := lib.Get("redis")

	fmt.Println(reply, err)
}
