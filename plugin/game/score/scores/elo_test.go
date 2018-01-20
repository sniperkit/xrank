package scores

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

var winnerScore, loserScore float64 = 200, 200

func TestCorrectScores(t *testing.T) {
	eloAlgo := NewEloScoreEngine(20)

	score1, score2 := eloAlgo.Compute(winnerScore, loserScore)
	assert.Equal(t, 210, int(score1))
	assert.Equal(t, 190, int(score2))
}
