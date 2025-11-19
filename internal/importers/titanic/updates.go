package titanic

import (
	"time"

	"github.com/Lekuruu/osutp/internal/common"
	"github.com/Lekuruu/osutp/internal/services"
)

func (importer *TitanicImporter) EnqueueUserUpdate(userID int, state *common.State) error {
	player, err := services.FetchPlayerById(userID, state)
	if err != nil {
		state.Logger.Logf("Failed to fetch player %d for update: %v", userID, err)
		return err
	}

	if player.IsUpdating {
		state.Logger.Logf("Player %d is already updating", player.ID)
		return nil
	}

	if player.WasRecentlyUpdated() {
		state.Logger.Logf("Player %d was recently updated, skipping", player.ID)
		return nil
	}

	state.Logger.Logf("Enqueuing update for player %d", player.ID)
	services.SetPlayerUpdatingStatus(player.ID, true, state)

	go func() {
		defer func() {
			// Handle any panics during import
			if r := recover(); r != nil {
				state.Logger.Logf("Recovered from panic while updating player %d: %v", player.ID, r)
				services.SetPlayerUpdatingStatus(player.ID, false, state)
			}
		}()

		_, importErr := importer.ImportUser(player.ID, state)
		if importErr != nil {
			state.Logger.Logf("Failed to import user %d: %v", player.ID, importErr)
			services.SetPlayerUpdatingStatus(player.ID, false, state)
			return
		}

		services.UpdatePlayerLastUpdate(player.ID, time.Now(), state)
		services.SetPlayerUpdatingStatus(player.ID, false, state)
		state.Logger.Logf("Finished updating player %d", player.ID)
	}()
	return nil
}
