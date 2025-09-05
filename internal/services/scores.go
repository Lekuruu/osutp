package services

import (
	"github.com/Lekuruu/osutp-web/internal/common"
	"github.com/Lekuruu/osutp-web/internal/database"
)

func CreateScore(score *database.Score, state *common.State) error {
	result := state.Database.Create(score)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func FetchScoreById(id int, state *common.State) (*database.Score, error) {
	score := &database.Score{}
	result := state.Database.First(score, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return score, nil
}

func ScoreExists(id int, state *common.State) (bool, error) {
	var count int64
	result := state.Database.
		Model(&database.Score{}).
		Where("id = ?", id).
		Count(&count)
	if result.Error != nil {
		return false, result.Error
	}
	return count > 0, nil
}

func FetchPersonalBestScore(playerID int, beatmapID int, state *common.State) (*database.Score, error) {
	score := &database.Score{}
	result := state.Database.
		Where("player_id = ? AND beatmap_id = ?", playerID, beatmapID).
		Order("total_tp DESC").
		First(score)
	if result.Error != nil {
		return nil, result.Error
	}
	return score, nil
}

func DeleteScore(id int, state *common.State) error {
	result := state.Database.Delete(&database.Score{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
