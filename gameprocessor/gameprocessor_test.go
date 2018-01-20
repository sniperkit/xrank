package gameprocessor

import (
	"testing"
	"github.com/justwy/ranking-engine/models"
	"fmt"
)

var TEAM1 = models.Team{
	Id: "yin",
	Score: 200,
}

var TEAM2 = models.Team{
	Id: "ryan",
	Score: 200,
}

func TestCorrectScores(t *testing.T) {
	fmt.Printf("%v\n", models.Tournament{"a", "b"})
}
