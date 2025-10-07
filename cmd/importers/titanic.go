package main

import (
	"github.com/Lekuruu/osutp/internal/common"
	"github.com/Lekuruu/osutp/internal/importers/titanic"
	"github.com/Lekuruu/osutp/internal/services"
)

func main() {
	state := common.NewState()
	if state == nil {
		return
	}

	// Update logger name
	state.Logger = common.NewLogger("titanic")

	err := titanic.ImportBeatmapsByDifficulty(0, state)
	if err != nil {
		state.Logger.Logf("Error occurred while importing beatmaps: %v", err)
		return
	}

	beatmapBatch, err := services.FetchBeatmapsByDifficulty(0, 1000, 0, []string{}, state)
	if err != nil {
		state.Logger.Logf("Error occurred while fetching beatmaps: %v", err)
		return
	}

	titanic.ImportOrUpdateLeaderboards(beatmapBatch, state)
	titanic.UpdatePlayerRatings(state)
	state.Logger.Log("Import completed.")
}
