package main

import (
	"flag"
	"github.com/justwy/ranking-engine/gameprocessor"
	"github.com/justwy/ranking-engine/ranking"
	"github.com/justwy/ranking-engine/scores"
	"github.com/justwy/ranking-engine/models"
	"fmt"
)

const K_FACTOR = 20

var teamById = map[string]models.Team{
	"yinwng": {
		Id: "yinwng",
		Score: 400,
	},
	"ryanmark": {
		Id: "ryanmark",
		Score: 400,
	},
}

func main() {
	input := flag.String("input", "/tmp/games.tsv", "input location of games tsv file")
	flag.Parse()

	gameProcessor := gameprocessor.NewTsvFileGameResultProcessor()

	gameResults := gameProcessor.Read(*input)

	scoreEngine := scores.NewEloScoreEngine(K_FACTOR)

	rankingEngine := ranking.NewRankEngine(scoreEngine)

	rankingEngine.UpdateScore(teamById, gameResults)

	teams := []models.Team{}

	for _, team := range teamById {
		teams = append(teams, team)
	}

	rankingEngine.Rank(teams)

	fmt.Println(teams)
}
