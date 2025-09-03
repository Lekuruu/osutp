package services

import (
	"github.com/Lekuruu/osutp-web/internal/common"
	"github.com/Lekuruu/osutp-web/internal/database"
)

func BeatmapExists(id int, state *common.State) (bool, error) {
	var count int64
	err := state.Database.Model(&database.Beatmap{}).Where("id = ?", id).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
