package model

import (
	"fmt"
)

// structs match reddit API @ www.reddit.com/subreddit/about.json
type Subreddit struct {
	Name string
	Data SubredditData
}

type SubredditData struct {
	Accounts_active int
	Subscribers     int
	Display_name    string
	AverageActive   int
	PercentActive   float32

	/*  Unused fields
	Url                string
	Public_description string
	Description_html   string
	Header_img         string
	Header_title       string
	Title              string
	Description        string
	Name               string
	*/
}

func (sr Subreddit) String() string {
	return fmt.Sprintf("%s - Subs: %d, AvgActive: %d, PercentActive: %.05f",
		sr.Name, sr.Data.Subscribers, sr.Data.AverageActive, sr.Data.PercentActive)
}

func (sr Subreddit) CalcPercentActive() {

	sr.Data.PercentActive = 100.0 * float32(sr.Data.Accounts_active/sr.Data.Subscribers)
	fmt.Printf("%f\n", sr.Data.PercentActive)
}

// returns list of subreddits to watch
func GetSubreddits() []Subreddit {
	names := [...]string{
		"actionscript",
		"asm",
		"c_programming",
		"clojure",
		"cobol",
		"cpp",
		"csharp",
		"d_language",
		"delphi",
		"dotnet",
		"fortran",
		"fsharp",
		"golang",
		"haskell",
		"iosprogramming",
		"java",
		"javascript",
		"julia",
		"lisp",
		"lua",
		"matlab",
		"pascal",
		"perl",
		"php",
		"postscript",
		"python",
		"ruby",
		"rust",
		"scala",
		"sql",
		"swift",
		"visualbasic",
	}

	subreddits := []Subreddit{}

	for i, _ := range names {
		subreddits = append(subreddits, newSubreddit(names[i]))
	}

	return subreddits
}

// helper returns named subreddit, required for url in update package
func newSubreddit(name string) Subreddit {
	var sr Subreddit
	sr.Name = name
	return sr
}
