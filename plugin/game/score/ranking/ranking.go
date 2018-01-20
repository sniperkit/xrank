package ranking

import (
	"github.com/justwy/ranking-engine/models"
	"github.com/justwy/ranking-engine/scores"
	"sort"
)

type RankEngine interface {
	UpdateScore(teamById map[string]models.Team, games []models.GameResult)
	Rank(teams []models.Team)
}

type rankEngineImpl struct {
	scoreEngine scores.ScoreEngine
}

func NewRankEngine(scoreEngine scores.ScoreEngine) RankEngine {
	return rankEngineImpl{
		scoreEngine: scoreEngine,
	}
}

func (rankEngine rankEngineImpl) Rank(teams []models.Team) {
	sort.Slice(teams, func(i, j int) bool {
		return teams[i].Score > teams[j].Score
	})
}

// If a team in a game is not in the list of teams, then panic
func (rankEngine rankEngineImpl) UpdateScore(teamById map[string]models.Team, games []models.GameResult) {
	for _, game := range games {
		winner, exist := teamById[game.WinnerId];
		if !exist {
			panic("team is not registered in the ranking system. " + game.WinnerId)
		}

		loser, exist := teamById[game.LoserId];
		if !exist {
			panic("team is not registered in the ranking system. " + game.LoserId)
		}

		newWinnerScore, newLoserScore := rankEngine.scoreEngine.Compute(winner.Score, loser.Score)
		winner.Score = newWinnerScore
		loser.Score = newLoserScore
		teamById[game.WinnerId] = winner
		teamById[game.LoserId] = loser
	}
}
