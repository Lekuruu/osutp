package updaters

import (
	"sort"
	"time"

	"github.com/Lekuruu/osutp/internal/common"
	"github.com/Lekuruu/osutp/internal/database"
	"github.com/Lekuruu/osutp/internal/services"
)

func UpdatePlayerRatings(state *common.State) error {
	players, err := services.FetchAllPlayers(state)
	if err != nil {
		return err
	}

	// Update player TP ratings based on their best scores
	if err := updatePlayerTpRatings(players, state); err != nil {
		return err
	}

	// Update player ranks after ratings are updated
	if err := updatePlayerRanks(players, state); err != nil {
		return err
	}
	return nil
}

// Update each player's rating based on their best scores
func updatePlayerTpRatings(players []*database.Player, state *common.State) error {
	for _, player := range players {
		bestScores, err := services.FetchPersonalBestScores(player.ID, state)
		if err != nil {
			state.Logger.Log("Failed to fetch scores for player", player.ID, ":", err)
			continue
		}

		if err := updateTpRatingForPlayer(player, bestScores, state); err != nil {
			state.Logger.Log("Failed to process player", player.ID, ":", err)
			continue
		}

		state.Logger.Logf("Updated player '%s' with total tp: %.2f (#%d)", player.Name, player.TotalTp, player.GlobalRank)
	}

	// Save updated player data
	if err := state.Database.Save(&players).Error; err != nil {
		return err
	}
	return nil
}

// Update player ranks based on their total TP
func updatePlayerRanks(players []*database.Player, state *common.State) error {
	sort.Slice(players, func(i, j int) bool {
		return players[i].TotalTp > players[j].TotalTp
	})
	playersByCountry := make(map[string][]*database.Player)

	// Update global rank & rank delta
	for index, player := range players {
		previousRank := player.GlobalRank
		player.GlobalRank = index + 1
		rankDelta := previousRank - player.GlobalRank
		player.RecentRankChange = rankDelta
		player.LastUpdate = time.Now().UTC()
		playersByCountry[player.Country] = append(playersByCountry[player.Country], player)

		if previousRank == 0 {
			// No rank change for new players
			player.RecentRankChange = 0
		}
	}

	// Update country ranks
	for _, countryPlayers := range playersByCountry {
		sort.Slice(countryPlayers, func(i, j int) bool {
			return countryPlayers[i].TotalTp > countryPlayers[j].TotalTp
		})
		for index, player := range countryPlayers {
			player.CountryRank = index + 1
		}
	}

	if err := state.Database.Save(&players).Error; err != nil {
		return err
	}
	return nil
}

func updateTpRatingForPlayer(player *database.Player, bestScores []database.Score, state *common.State) error {
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
