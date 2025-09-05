package titanic

import (
	"fmt"
	"time"

	"github.com/Lekuruu/osutp-web/internal/common"
	"github.com/Lekuruu/osutp-web/internal/database"
	"github.com/Lekuruu/osutp-web/internal/services"
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

	for _, player := range players {
		if err := ProcessPlayerRating(player.ID, state); err != nil {
			fmt.Println("Failed to process player", player.ID, ":", err)
			continue
		}
		if err := state.Database.Save(player).Error; err != nil {
			fmt.Println("Failed to save player", player.ID, ":", err)
			continue
		}
	}
	return nil
}

func ProcessPlayerRating(playerId int, state *common.State) error {
	player, err := services.FetchPlayerById(playerId, state)
	if err != nil {
		return err
	}

	bestScores, err := services.FetchPersonalBestScores(playerId, state)
	if err != nil {
		return err
	}

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
