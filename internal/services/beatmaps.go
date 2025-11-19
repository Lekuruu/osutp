package services

import (
	"fmt"
	"time"

	"github.com/Lekuruu/osutp/internal/common"
	"github.com/Lekuruu/osutp/internal/database"
)

func CreateBeatmap(beatmap *database.Beatmap, state *common.State) error {
	result := state.Database.Create(beatmap)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func FetchBeatmapById(id int, state *common.State) (*database.Beatmap, error) {
	beatmap := &database.Beatmap{}
	result := state.Database.First(beatmap, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return beatmap, nil
}

func FetchBeatmapsBySetId(setId int, state *common.State) ([]database.Beatmap, error) {
	var beatmaps []database.Beatmap
	result := state.Database.Where("set_id = ?", setId).Find(&beatmaps)
	if result.Error != nil {
		return nil, result.Error
	}

	return beatmaps, nil
}

func FetchTotalBeatmaps(filters []string, state *common.State) (int64, error) {
	var count int64
	query := state.Database.Model(&database.Beatmap{})
	for _, filter := range filters {
		query = query.Where(filter)
	}
	err := query.Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func FetchBeatmapsByDifficulty(offset int, limit int, mods uint32, filters []string, state *common.State) ([]*database.Beatmap, error) {
	var beatmaps []*database.Beatmap
	query := state.Database.Model(&beatmaps).
		Where("difficulty_attributes IS NOT NULL").
		Order(fmt.Sprintf("json_extract(difficulty_attributes, '$.%d.StarRating') DESC", mods)).
		Offset(offset).
		Limit(limit)

	for _, filter := range filters {
		query = query.Where(filter)
	}

	result := query.Find(&beatmaps)
	if result.Error != nil {
		return nil, result.Error
	}
	return beatmaps, nil
}

func BeatmapExists(id int, state *common.State) (bool, error) {
	var count int64
	err := state.Database.Model(&database.Beatmap{}).Where("id = ?", id).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func UpdateBeatmapStatus(id int, status int, state *common.State) error {
	result := state.Database.Model(&database.Beatmap{}).Where("id = ?", id).Update("status", status)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func UpdateBeatmapLastScoreUpdate(id int, timestamp time.Time, state *common.State) error {
	result := state.Database.Model(&database.Beatmap{}).Where("id = ?", id).Update("last_score_update", timestamp)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
