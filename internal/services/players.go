package services

import (
	"github.com/Lekuruu/osutp-web/internal/common"
	"github.com/Lekuruu/osutp-web/internal/database"
)

func PlayerUser(player *database.Player, state *common.State) error {
	result := state.Database.Create(player)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func FetchPlayerById(playerId int, state *common.State) (*database.Player, error) {
	player := &database.Player{}
	result := state.Database.First(player, playerId)
	if result.Error != nil {
		return nil, result.Error
	}
	return player, nil
}

func FetchPlayerByName(name string, state *common.State) (*database.Player, error) {
	player := &database.Player{}
	result := state.Database.Where("LOWER(name) = LOWER(?)", name).First(player)
	if result.Error != nil {
		return nil, result.Error
	}
	return player, nil
}

func FetchAllPlayers(state *common.State) ([]*database.Player, error) {
	var players []*database.Player
	result := state.Database.Find(&players)
	if result.Error != nil {
		return nil, result.Error
	}
	return players, nil
}

func PlayerExists(playerId int, state *common.State) (bool, error) {
	var count int64
	if err := state.Database.Model(&database.Player{}).Where("id = ?", playerId).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func DeletePlayer(playerId int, state *common.State) error {
	result := state.Database.Delete(&database.Player{}, playerId)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
