package titanic

import (
	"time"

	"github.com/Lekuruu/osutp/internal/common"
	"github.com/Lekuruu/osutp/internal/services"
)

func (importer *TitanicImporter) EnqueueUserUpdate(userID int, state *common.State) error {
	player, err := services.FetchPlayerById(userID, state)
	if err != nil {
		return err
	}

	if player.IsUpdating {
		return nil
	}

	if player.WasRecentlyUpdated() {
		return nil
	}

	services.SetPlayerUpdatingStatus(player.ID, true, state)

	go func() {
		defer func() {
			// Handle any panics during import
			if r := recover(); r != nil {
				services.SetPlayerUpdatingStatus(player.ID, false, state)
			}
		}()

		_, importErr := importer.ImportUser(player.ID, state)

		if importErr != nil {
			services.SetPlayerUpdatingStatus(player.ID, false, state)
			return
		}

		services.UpdatePlayerLastUpdate(player.ID, time.Now(), state)
		services.SetPlayerUpdatingStatus(player.ID, false, state)
	}()
	return nil
}
