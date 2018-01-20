package lib

import (
	//"fmt"
	"strconv"
	"sync"
)

const (
	//github api v3
	API_V3 = "https://api.github.com"

	API_REPO_SEARCH = "/search/repositories"
)

type Repo struct {
	Language string
	Limit    int
}

type RepoResponse struct {
	TotalCount int64 `json:"total_count"`
	Items      []Git `json:"items"`
}

func (this *Repo) Search() {
	params := make(map[string]string)
	params["sort"] = "stars"
	params["order"] = "desc"
	params["q"] = "language:go stars:5..20"
	params["per_page"] = strconv.Itoa(this.Limit)
	//fmt.Println(params)

	//url := "https://api.github.com/search/repositories?sort=stars&order=desc&q=language:go"
	request := &Request{}
	var repoResp RepoResponse
	request.Get(API_V3+API_REPO_SEARCH, params, &repoResp)
	//request.Get(url, nil)
	//fmt.Println(repoResp)

	repositories := repoResp.Items
	size := len(repositories)
	if size > 0 {
		for _, git := range repositories {
			Queue.Push(git)
		}
	}
	done := make(chan bool, 1)
	var wait sync.WaitGroup
	wait.Add(1)
	go WordCount(done, &wait)
	Start(size, done)
	wait.Wait()
}
