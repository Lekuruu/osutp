package services

import (
	"github.com/Lekuruu/osutp-web/internal/common"
	"github.com/Lekuruu/osutp-web/internal/database"
)

func CreateUser(player *database.Player, state *common.State) error {
	result := state.Database.Create(player)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func FetchUserById(userID int, state *common.State) (*database.Player, error) {
	player := &database.Player{}
	result := state.Database.First(player, userID)
	if result.Error != nil {
		return nil, result.Error
	}
	return player, nil
}

func FetchUserByName(name string, state *common.State) (*database.Player, error) {
	player := &database.Player{}
	result := state.Database.Where("name = ?", name).First(player)
	if result.Error != nil {
		return nil, result.Error
	}
	return player, nil
}

func UserExists(userID int, state *common.State) (bool, error) {
	var count int64
	if err := state.Database.Model(&database.Player{}).Where("id = ?", userID).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
