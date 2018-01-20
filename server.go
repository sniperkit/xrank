package main

import (
	"io/ioutil"
	"net/http"

	"encoding/json"
	"fmt"

	"time"

	"log"

	"code.google.com/p/leveldb-go/leveldb"
	"code.google.com/p/leveldb-go/leveldb/db"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
)

const (
	gitHubStarURL = "https://api.github.com/search/repositories?q=language:%s&sort=star&order=desc"
	gitHubUserURL = "https://github.com/%s"

	snapshotKey = "lang:%s_data:%s"
)

// Ranking struct GitHub Star Ranking.
type Ranking struct {
	Item
	Rank     int64
	LastRank int64
	//	FullName string
	//	StargazersCount int64
	LastStargazersCount int64
	UserURL             string
	//	OwnerAvatarURL string
	//	HtmlURL string
	//	UpdatedAt       string
	//	CreatedAt       string
}

// parse time default.
func (r Ranking) ParseTime(at string) time.Time {
	t, _ := time.Parse(time.RFC3339, at)
	return t
}

// format time default.
func (r Ranking) FormatTime(at string) string {
	return r.ParseTime(at).Format("2006-01-02")
}

// ResponseRanking struct GitHub star Ranking for each Language.
type ResponseRanking struct {
	Language string
	Rankings []Ranking
}

func doTop(r render.Render) {
	// render
	r.HTML(200, "index", nil)
}

func doRanking(params martini.Params, r render.Render) {

	lang := params["language"]
	rankings, err := readGitHubStarRanking(lang)
	if err != nil {
		r.Error(400)
	}

	// render
	r.HTML(200, "ranking", ResponseRanking{
		Language: lang,
		Rankings: rankings,
	})
}

func createSnapshotKey(lang string, t time.Time) string {
	return fmt.Sprintf(snapshotKey, lang, t.Format("2006-01-02"))
}

func doSnapshot(logger *log.Logger, params martini.Params, r render.Render) {
	lang := params["language"]

	// levelDBへ保存
	level, err := leveldb.Open("snapshot", &db.Options{})
	if err != nil {
		r.Error(400)
	}
	defer level.Close()

	key := createSnapshotKey(lang, time.Now())
	logger.Println("key: ", key)
	if _, err := level.Get([]byte(key), &db.ReadOptions{}); err != nil {
		res, err := http.Get(fmt.Sprintf(gitHubStarURL, lang))
		if err != nil {
			r.Error(400)
		}
		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			r.Error(400)
		}
		if err := level.Set([]byte(key), data, &db.WriteOptions{}); err != nil {
			r.Error(400)
		}
	}

	r.JSON(200, nil)
}

func main() {
	m := martini.Classic()
	m.Use(martini.Static("public"))
	m.Use(render.Renderer(render.Options{
		Directory:  "templates",
		Extensions: []string{".tmpl"},
	}))

	// Router
	m.Get("/", doTop)
	m.Get("/ranking/:language", doRanking)
	m.Get("/snapshot/:language", doSnapshot)

	m.Run()
}

func readGitHubStarRanking(lang string) ([]Ranking, error) {

	res, err := http.Get(fmt.Sprintf(gitHubStarURL, lang))
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	// parse
	var reps Repositories
	if err := json.Unmarshal(data, &reps); err != nil {
		return nil, err
	}

	var snapshot Repositories
	level, err := leveldb.Open("snapshot", &db.Options{})
	if err != nil {
		return nil, err
	}
	defer level.Close()

	key := createSnapshotKey(lang, time.Now().AddDate(0, 0, -1))
	if d, err := level.Get([]byte(key), &db.ReadOptions{}); err == nil {
		if err := json.Unmarshal(d, &snapshot); err != nil {
			return nil, err
		}
	}

	return newRankings(&reps, &snapshot), nil
}

func newRankings(reps, snapshot *Repositories) []Ranking {

	rankings := make([]Ranking, len(reps.Items))
	for idx, item := range reps.Items {
		var lastRank int64
		var lastStargazersCount int64
		for sIdx, sItem := range snapshot.Items {
			if sItem.FullName != item.FullName {
				continue
			}
			lastRank = int64(sIdx + 1)
			lastStargazersCount = sItem.StargazersCount
			break
		}

		rankings[idx] = Ranking{
			Item:                item,
			Rank:                int64(idx + 1),
			LastRank:            lastRank,
			LastStargazersCount: lastStargazersCount,
			UserURL:             fmt.Sprintf(gitHubUserURL, item.Owner.Login),
		}
	}
	return rankings
}
