package titanic

import (
	"sort"
	"time"

	"github.com/Lekuruu/osutp/internal/common"
	"github.com/Lekuruu/osutp/internal/database"
	"github.com/Lekuruu/osutp/internal/services"
)

type UserModel struct {
	ID             int     `json:"id"`
	Name           string  `json:"name"`
	Country        string  `json:"country"`
	CreatedAt      string  `json:"created_at"`
	LatestActivity string  `json:"latest_activity"`
	Restricted     bool    `json:"restricted"`
	Activated      bool    `json:"activated"`
	PreferredMode  int     `json:"preferred_mode"`
	Playstyle      int     `json:"playstyle"`
	Banner         *string `json:"banner,omitempty"`
	Website        *string `json:"website,omitempty"`
	Discord        *string `json:"discord,omitempty"`
	Twitter        *string `json:"twitter,omitempty"`
	Location       *string `json:"location,omitempty"`
	Interests      *string `json:"interests,omitempty"`
}

func (user *UserModel) ToSchema() *database.Player {
	createdAt, err := time.Parse("2006-01-02T15:04:05", user.CreatedAt)
	if err != nil {
		createdAt = time.Now().UTC()
	}

	return &database.Player{
		ID:        user.ID,
		Name:      user.Name,
		Country:   user.Country,
		CreatedAt: createdAt,
	}
}

func UpdatePlayerRatings(state *common.State) error {
	players, err := services.FetchAllPlayers(state)
	if err != nil {
		return err
	}

	// Update each player's rating based on their best scores
	for _, player := range players {
		bestScores, err := services.FetchPersonalBestScores(player.ID, state)
		if err != nil {
			state.Logger.Log("Failed to fetch scores for player", player.ID, ":", err)
			continue
		}

		if err := common.UpdatePlayerRating(player, bestScores, state); err != nil {
			state.Logger.Log("Failed to process player", player.ID, ":", err)
			continue
		}

		state.Logger.Logf("Updated player '%s' with total tp: %.2f (#%d)", player.Name, player.TotalTp, player.GlobalRank)
	}

	// Save updated player data
	if err := state.Database.Save(&players).Error; err != nil {
		return err
	}

	// Sort players by total tp
	sort.Slice(players, func(i, j int) bool {
		return players[i].TotalTp > players[j].TotalTp
	})
	playersByCountry := make(map[string][]*database.Player)

	// Update global rank & rank delta
	state.Logger.Log("Updating global ranks...")

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
	state.Logger.Log("Updated country ranks...")

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
