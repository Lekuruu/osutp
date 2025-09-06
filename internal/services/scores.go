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

func FetchPersonalBestScore(playerId int, beatmapID int, state *common.State) (*database.Score, error) {
	score := &database.Score{}
	result := state.Database.
		Where("player_id = ? AND beatmap_id = ?", playerId, beatmapID).
		Order("total_tp DESC").
		First(score)
	if result.Error != nil {
		return nil, result.Error
	}
	return score, nil
}

func FetchPersonalBestScores(playerId int, state *common.State) ([]database.Score, error) {
	var scores []database.Score
	allowedStatuses := []int{1, 2} // Ranked and Approved
	result := state.Database.
		Joins("JOIN beatmaps ON beatmaps.id = scores.beatmap_id").
		Where("scores.player_id = ?", playerId).
		Where("beatmaps.status IN ?", allowedStatuses).
		Order("scores.total_tp DESC").
		Find(&scores)
	if result.Error != nil {
		return nil, result.Error
	}
	return scores, nil
}

func FetchBestScores(offset int, limit int, sort string, state *common.State) ([]database.Score, error) {
	var scores []database.Score
	preload := []string{"Player", "Beatmap"}
	result := database.
		PreloadQuery(state.Database, preload).
		Order(sort).
		Offset(offset).
		Limit(limit).
		Find(&scores)
	if result.Error != nil {
		return nil, result.Error
	}
	return scores, nil
}

func FetchTotalScores(state *common.State) (int64, error) {
	var count int64
	err := state.Database.Model(&database.Score{}).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func DeleteScore(id int, state *common.State) error {
	result := state.Database.Delete(&database.Score{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func DeleteScoresByPlayer(playerId int, state *common.State) error {
	result := state.Database.Where("player_id = ?", playerId).Delete(&database.Score{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}
