package analytics

import (
	"fmt"
	"github.com/steaz/reddilytics/model"
	"sort"
)

// load most recent data.

func Analyze() {
	subreddits := model.LoadAll()

	sort.Sort(BySubs(subreddits))
	fmt.Println("Sorted by Subscribers: \n")
	for i, sr := range subreddits {
		fmt.Printf("%d. %s (%d)\n", i+1, sr.Name, sr.Data.Subscribers)
	}
	/*
		sort.Sort(ByActivity(subreddits))
		fmt.Println("\n\nSorted By Average Percent Activity: \n")
		for i, sr := range subreddits {
			fmt.Printf("%d. %s (%.03f %%)\n", i+1, sr.Name, sr.Data.PercentActive)
		}
	*/
}

type BySubs []model.Subreddit

func (a BySubs) Len() int           { return len(a) }
func (a BySubs) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a BySubs) Less(i, j int) bool { return a[i].Data.Subscribers > a[j].Data.Subscribers }

type ByActivity []model.Subreddit

func (a ByActivity) Len() int           { return len(a) }
func (a ByActivity) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByActivity) Less(i, j int) bool { return a[i].Data.PercentActive > a[j].Data.PercentActive }
