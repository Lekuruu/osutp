package common

import (
	"github.com/Lekuruu/osutp-web/internal/database"
)

func UpdatePlayerRating(player *database.Player, bestScores []database.Score, state *State) error {
	tpValues := make([]float64, 0, len(bestScores))
	aimValues := make([]float64, 0, len(bestScores))
	speedValues := make([]float64, 0, len(bestScores))
	accValues := make([]float64, 0, len(bestScores))

	for _, score := range bestScores {
		tpValues = append(tpValues, score.TotalTp)
		speedValues = append(speedValues, score.SpeedTp)
		aimValues = append(aimValues, score.AimTp)
		accValues = append(accValues, score.AccTp)
	}

	player.TotalTp = calculateWeightedRating(tpValues)
	player.AimTp = calculateWeightedRating(aimValues)
	player.SpeedTp = calculateWeightedRating(speedValues)
	player.AccuracyTp = calculateWeightedRating(accValues)
	return nil
}

func calculateWeightedRating(values []float64) float64 {
	if len(values) == 0 {
		return 0.0
	}

	total := 0.0
	factor := 1.0
	baseFactor := 1.0

	for _, value := range values {
		total += value*factor + 0.25*baseFactor
		factor *= 0.95
		baseFactor *= 0.9994
	}

	return total
}
