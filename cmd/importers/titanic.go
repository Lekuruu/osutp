package main

import (
	"fmt"

	"github.com/Lekuruu/osutp-web/internal/common"
	"github.com/Lekuruu/osutp-web/internal/importers/titanic"
	"github.com/Lekuruu/osutp-web/internal/services"
)

func main() {
	state := common.NewState()
	if state == nil {
		return
	}

	err := titanic.ImportBeatmapsByDifficulty(0, state)
	if err != nil {
		fmt.Printf("Error occurred while importing beatmaps: %v\n", err)
	}

	beatmapBatch, err := services.FetchBeatmapsByDifficulty(0, 1000, 0, []string{}, state)
	if err != nil {
		fmt.Printf("Error occurred while fetching beatmaps: %v\n", err)
		return
	}

	titanic.ImportOrUpdateLeaderboards(beatmapBatch, state)
	fmt.Println("Import completed.")
}
