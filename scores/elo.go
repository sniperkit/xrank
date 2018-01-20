package scores

import "math"

type eloAlgo struct {
	kFactor float64
}

func (eloAlgo eloAlgo) Compute(winnerScore, loserScore float64) (float64, float64)  {

	return calculateElo(eloAlgo.kFactor, winnerScore, loserScore)
}

func calculateElo(kFactor, winnerScore, loserScore float64) (float64, float64) {
	exponent := math.Pow(10, float64((winnerScore - loserScore) / 400))
	ex := 1.0 / (1 + exponent)

	change := kFactor * ex

	return winnerScore + change, loserScore - change
}

