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
		bestScores, err := services.FetchPersonalBestScores(player.ID, state)
		if err != nil {
			fmt.Println("Failed to fetch scores for player", player.ID, ":", err)
			continue
		}

		if err := common.UpdatePlayerRating(player, bestScores, state); err != nil {
			fmt.Println("Failed to process player", player.ID, ":", err)
			continue
		}
	}
	return nil
}
