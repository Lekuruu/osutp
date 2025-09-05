package services

import (
	"fmt"

	"github.com/Lekuruu/osutp-web/internal/common"
	"github.com/Lekuruu/osutp-web/internal/database"
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

func FetchTotalBeatmaps(state *common.State) (int64, error) {
	var count int64
	err := state.Database.Model(&database.Beatmap{}).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func FetchBeatmapsByDifficulty(offset int, limit int, mods int, state *common.State) ([]*database.Beatmap, error) {
	var beatmaps []*database.Beatmap
	result := state.Database.Model(&beatmaps).
		Where("difficulty_attributes IS NOT NULL").
		Order(fmt.Sprintf("json_extract(difficulty_attributes, '$.%d.StarRating') DESC", mods)).
		Offset(offset).
		Limit(limit).
		Find(&beatmaps)
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
