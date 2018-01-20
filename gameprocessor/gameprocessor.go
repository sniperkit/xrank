package gameprocessor

import "github.com/justwy/ranking-engine/models"
import (
	"io/ioutil"
	"strings"
	"strconv"
	"fmt"
)

// GameResultReader reads game results from a uri.
// Note this is not stream reader. It loads all bytes behind the uri all into memory
// TODO - Need to change to stream reader if input is large
type GameResultProcessor interface {
	Read(uri string) []models.GameResult
}

// line format: winner_id	loser_id	tournament_id	sequence_id
type tsvFileProcessor struct {}

const IDX_WINNER_ID = 0
const IDX_LOSER_ID = 1
const IDX_TOURNAMENT_ID = 2
const IDX_SEQUENCE_ID = 3

func (processor tsvFileProcessor) Read(uri string) []models.GameResult {
	gameResults := []models.GameResult{}

	bytes, err := ioutil.ReadFile(uri)
	if err != nil {
		panic("cannot read game result file: " + uri)
	}

	games := string(bytes)

	fmt.Println(games)

	lines := strings.Split(games, "\n")
	for _, line := range lines {
		tokens := strings.Split(line, "\t")

		if len(tokens) < 4 {
			continue
		}

		sequence, _ := strconv.Atoi(tokens[IDX_SEQUENCE_ID])

		gameResult := models.GameResult{
			WinnerId: tokens[IDX_WINNER_ID],
			LoserId: tokens[IDX_LOSER_ID],
			Sequence: sequence,
			TournamentId: tokens[IDX_TOURNAMENT_ID],
		}

		gameResults = append(gameResults, gameResult)
	}

	return gameResults
}

func NewTsvFileGameResultProcessor() GameResultProcessor {
	return tsvFileProcessor{}
}
