package main

import (
	"log"
	"time"

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
			// There's a chance we are being rate limited, so let's wait here
			time.Sleep(time.Second * 60)
			continue
		}
		if amount == 0 {
			break
		}

		state.Logger.Logf("Imported %d beatmaps from page %d", amount, page)
	}

	updaters.UpdatePlayerRatings(state)
	state.Logger.Log("Import completed.")
}
