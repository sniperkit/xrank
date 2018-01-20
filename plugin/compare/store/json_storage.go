package store

import (
	"encoding/json"
	"github.com/mwstobo/rank/rankings"
	"io/ioutil"
	"os"
)

type RankingJson struct {
	Ranking []string `json:"ranking"`
}

type JsonStorage struct {
	filename string
}

func NewJsonStorage(filename string) *JsonStorage {
	return &JsonStorage{
		filename: filename,
	}
}

func (storage *JsonStorage) Import() ([]string, error) {
	rankingData, err := ioutil.ReadFile(storage.filename)
	if os.IsNotExist(err) {
		return make([]string, 0), nil
	} else if err != nil {
		return nil, err
	}

	rankingJson := &RankingJson{}
	err = json.Unmarshal(rankingData, rankingJson)
	if err != nil {
		return nil, err
	}

	if rankingJson.Ranking == nil {
		return make([]string, 0), nil
	}
	return rankingJson.Ranking, nil
}

func (storage *JsonStorage) Export(ranking rankings.Ranking) error {
	rankingJson := &RankingJson{}
	rankingJson.Ranking = make([]string, ranking.Length())
	for i := 0; i < ranking.Length(); i += 1 {
		rankingJson.Ranking[i] = ranking.Select(i)
	}

	rankingData, err := json.Marshal(rankingJson)
	rankingData = append(rankingData, '\n')
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(storage.filename, rankingData, 0644)
	if err != nil {
		return err
	}

	return nil
}
