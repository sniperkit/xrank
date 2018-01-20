package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/google/go-github/github"
)

var (
	flags struct {
		sort     string
		lang     string
		pageSize int
	}
)

func init() {
	flag.StringVar(&flags.sort, "s", "stars", "sort, valid options: stars | forks")
	flag.StringVar(&flags.lang, "l", "go", "language, valid options: go | c | php")
	flag.IntVar(&flags.pageSize, "n", 50, "limit repos count")

	flag.Parse()
}

func main() {
	client := github.NewClient(nil)

	opt := &github.SearchOptions{}
	opt.Sort = flags.sort
	opt.Order = "desc"
	opt.PerPage = flags.pageSize
	result, _, err := client.Search.Repositories("language:"+flags.lang, opt)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%6s %5s %30s  %s\n", "star", "fork", "name", "desc")
	fmt.Printf("%6s %5s %30s  %s\n", strings.Repeat("=", 6),
		strings.Repeat("=", 5),
		strings.Repeat("=", 30),
		strings.Repeat("=", 30),
	)
	for _, repo := range result.Repositories {
		fmt.Printf("%6d %5d %30s  %s\n", *repo.StargazersCount, *repo.ForksCount,
			*repo.FullName, *repo.Description)
	}

	rate, _, _ := client.RateLimit()
	fmt.Printf("Rate Limit: %d/%d\n", rate.Remaining, rate.Limit)
}
