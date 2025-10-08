package main

import (
	"log"

	"github.com/Lekuruu/osutp/internal/common"
	"github.com/Lekuruu/osutp/internal/importers"
	"github.com/Lekuruu/osutp/internal/updaters"
)

func main() {
	state := common.NewState()
	if state == nil {
		return
	}

	importer, err := importers.NewImporter(state.Config)
	if err != nil {
		log.Fatalf("Failed to create importer: %v", err)
		return
	}

	// Update logger name
	state.Logger = common.NewLogger("titanic")

	for page := 0; true; page++ {
		amount, err := importer.ImportBeatmapsByDifficulty(page, state)
		if err != nil {
			state.Logger.Logf("Error occurred while importing beatmaps: %v", err)
			continue
		}

		// Update player rankings after each page
		updaters.UpdatePlayerRatings(state)

		if amount == 0 {
			// No more beatmaps to import
			break
		}
	}

	state.Logger.Log("Import completed.")
}
