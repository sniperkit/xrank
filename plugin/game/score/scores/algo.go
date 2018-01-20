package scores

type ScoreEngine interface {
	Compute(winnerScore, loserScore float64) (float64, float64)
}

func NewEloScoreEngine(kFactor float64) ScoreEngine {
	return eloAlgo{
		kFactor: kFactor,
	}
}
