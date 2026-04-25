package services

import (
	"github.com/Lekuruu/osutp/internal/common"
	"github.com/Lekuruu/osutp/internal/database"
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

func FetchRangePersonalBestScores(playerId int, offset int, limit int, sort string, state *common.State) ([]database.Score, error) {
	var scores []database.Score
	allowedStatuses := []int{1, 2} // Ranked and Approved
	result := state.Database.
		Preload("Player").
		Preload("Beatmap").
		Joins("JOIN beatmaps ON beatmaps.id = scores.beatmap_id").
		Where("scores.player_id = ?", playerId).
		Where("beatmaps.status IN ?", allowedStatuses).
		Order(sort).
		Offset(offset).
		Limit(limit).
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

func FetchBestScoresByBeatmap(beatmapId int, offset int, limit int, sort string, state *common.State) ([]database.Score, error) {
	var scores []database.Score
	preload := []string{"Player"}
	result := database.
		PreloadQuery(state.Database, preload).
		Where("beatmap_id = ?", beatmapId).
		Order(sort).
		Offset(offset).
		Limit(limit).
		Find(&scores)
	if result.Error != nil {
		return nil, result.Error
	}
	return scores, nil
}

func FetchTotalPersonalBestScores(playerId int, state *common.State) (int64, error) {
	var count int64
	allowedStatuses := []int{1, 2} // Ranked and Approved
	err := state.Database.
		Model(&database.Score{}).
		Joins("JOIN beatmaps ON beatmaps.id = scores.beatmap_id").
		Where("scores.player_id = ?", playerId).
		Where("beatmaps.status IN ?", allowedStatuses).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func FetchTotalScoresByBeatmap(beatmapId int, state *common.State) (int64, error) {
	var count int64
	err := state.Database.
		Model(&database.Score{}).
		Where("beatmap_id = ?", beatmapId).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
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
	_, err := DeleteScoreWithCount(id, state)
	return err
}

func DeleteScoreWithCount(id int, state *common.State) (int64, error) {
	result := state.Database.Delete(&database.Score{}, id)
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

func DeleteScoresByPlayer(playerId int, state *common.State) error {
	_, err := DeleteScoresByPlayerWithCount(playerId, state)
	return err
}

func DeleteScoresByPlayerWithCount(playerId int, state *common.State) (int64, error) {
	result := state.Database.Where("player_id = ?", playerId).Delete(&database.Score{})
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

func DeleteScoresByPlayerExcept(playerId int, retainedScoreIds []int, state *common.State) (int64, error) {
	var localScoreIds []int
	if err := state.Database.
		Model(&database.Score{}).
		Where("player_id = ?", playerId).
		Pluck("id", &localScoreIds).Error; err != nil {
		return 0, err
	}

	return deleteScoresExcept(localScoreIds, retainedScoreIds, state)
}

func DeleteScoresByBeatmap(beatmapId int, state *common.State) error {
	_, err := DeleteScoresByBeatmapWithCount(beatmapId, state)
	return err
}

func DeleteScoresByBeatmapWithCount(beatmapId int, state *common.State) (int64, error) {
	result := state.Database.Where("beatmap_id = ?", beatmapId).Delete(&database.Score{})
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

func DeleteScoresByBeatmapExcept(beatmapId int, retainedScoreIds []int, state *common.State) (int64, error) {
	var localScoreIds []int
	if err := state.Database.
		Model(&database.Score{}).
		Where("beatmap_id = ?", beatmapId).
		Pluck("id", &localScoreIds).Error; err != nil {
		return 0, err
	}

	return deleteScoresExcept(localScoreIds, retainedScoreIds, state)
}

func deleteScoresExcept(localScoreIds []int, retainedScoreIds []int, state *common.State) (int64, error) {
	retained := make(map[int]struct{}, len(retainedScoreIds))
	for _, id := range retainedScoreIds {
		retained[id] = struct{}{}
	}

	staleScoreIds := make([]int, 0)
	for _, id := range localScoreIds {
		if _, ok := retained[id]; !ok {
			staleScoreIds = append(staleScoreIds, id)
		}
	}

	var deleted int64
	for start := 0; start < len(staleScoreIds); start += 500 {
		end := min(start+500, len(staleScoreIds))

		result := state.Database.Delete(&database.Score{}, staleScoreIds[start:end])
		if result.Error != nil {
			return deleted, result.Error
		}
		deleted += result.RowsAffected
	}

	return deleted, nil
}
