package ranking

import (
	"fmt"
	"testing"
	"github.com/justwy/ranking-engine/models"
	"github.com/justwy/ranking-engine/scores"
)

var TEAM_BY_ID = map[string]models.Team{
	"team-one": {
		Id: "team-one",
		Score: 400,
	},
	"team-two": {
		Id: "team-two",
		Score: 200,
	},
	"team-three": {
		Id: "team-three",
		Score: 500,
	},
	"team-four": {
		Id: "team-four",
		Score: 300,
	},
}

var GAMES = []models.GameResult{
	{
		WinnerId: TEAM_BY_ID["team-one"].Id,
		LoserId: TEAM_BY_ID["team-two"].Id,
	},
	{
		WinnerId: TEAM_BY_ID["team-three"].Id,
		LoserId: TEAM_BY_ID["team-four"].Id,
	},
}

func TestRankEngine(t *testing.T) {
	rankEngine := NewRankEngine(scores.NewEloScoreEngine(20))

	rankEngine.UpdateScore(TEAM_BY_ID, GAMES);

	teams := []models.Team{}

	for _, value := range TEAM_BY_ID {
		teams = append(teams, value)
	}

	rankEngine.Rank(teams)

	fmt.Println(teams)
}

