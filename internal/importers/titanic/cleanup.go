package titanic

import (
	"fmt"

	"github.com/Lekuruu/osutp/internal/common"
	"github.com/Lekuruu/osutp/internal/services"
	"github.com/Lekuruu/osutp/internal/updaters"
)

func (importer *TitanicImporter) cleanupMissingBeatmap(beatmapID int, state *common.State) error {
	return importer.cleanupBeatmap(beatmapID, "confirmed remote 404", state)
}

func (importer *TitanicImporter) cleanupMissingBeatmapset(beatmapsetID int, state *common.State) error {
	localBeatmaps, err := services.FetchBeatmapsBySetId(beatmapsetID, state)
	if err != nil {
		state.Logger.Logf("Failed to fetch local beatmaps for missing beatmapset %d: %v", beatmapsetID, err)
		return err
	}

	for _, beatmap := range localBeatmaps {
		if err := importer.cleanupBeatmap(beatmap.ID, fmt.Sprintf("confirmed remote beatmapset %d 404", beatmapsetID), state); err != nil {
			return err
		}
	}
	return nil
}

func (importer *TitanicImporter) cleanupBeatmapsetRemovedBeatmaps(beatmapsetID int, remoteBeatmapIDs map[int]struct{}, state *common.State) error {
	localBeatmaps, err := services.FetchBeatmapsBySetId(beatmapsetID, state)
	if err != nil {
		state.Logger.Logf("Failed to fetch local beatmaps for beatmapset %d reconciliation: %v", beatmapsetID, err)
		return err
	}

	for _, beatmap := range localBeatmaps {
		if _, ok := remoteBeatmapIDs[beatmap.ID]; ok {
			continue
		}
		if err := importer.cleanupBeatmap(beatmap.ID, fmt.Sprintf("missing from remote beatmapset %d snapshot", beatmapsetID), state); err != nil {
			return err
		}
	}
	return nil
}

func (importer *TitanicImporter) cleanupBeatmap(beatmapID int, reason string, state *common.State) error {
	deletedScores, err := services.DeleteScoresByBeatmapWithCount(beatmapID, state)
	if err != nil {
		state.Logger.Logf("Failed to cleanup scores for missing beatmap %d: %v", beatmapID, err)
		return err
	}
	deletedBeatmaps, err := services.DeleteBeatmapWithCount(beatmapID, state)
	if err != nil {
		state.Logger.Logf("Failed to delete missing beatmap %d: %v", beatmapID, err)
		return err
	}
	if deletedScores == 0 && deletedBeatmaps == 0 {
		return nil
	}

	state.Logger.Logf("Cleaned up local beatmap %d after %s (%d beatmaps, %d scores removed)", beatmapID, reason, deletedBeatmaps, deletedScores)
	return importer.recomputeRatingsAfterCleanup(state)
}

func (importer *TitanicImporter) cleanupUnavailableUser(userID int, reason string, state *common.State) error {
	deletedScores, err := services.DeleteScoresByPlayerWithCount(userID, state)
	if err != nil {
		state.Logger.Logf("Failed to cleanup scores for %s user %d: %v", reason, userID, err)
		return err
	}
	deletedPlayers, err := services.DeletePlayerWithCount(userID, state)
	if err != nil {
		state.Logger.Logf("Failed to delete %s user %d: %v", reason, userID, err)
		return err
	}
	if deletedScores == 0 && deletedPlayers == 0 {
		return nil
	}

	state.Logger.Logf("Cleaned up local data for %s user %d (%d players, %d scores removed)", reason, userID, deletedPlayers, deletedScores)
	return importer.recomputeRatingsAfterCleanup(state)
}

func (importer *TitanicImporter) cleanupUserIfUnavailable(userID int, state *common.State) (bool, error) {
	user, err := importer.fetchUserById(userID)
	if err != nil {
		return false, err
	}
	if user == nil {
		return true, importer.cleanupUnavailableUser(userID, "missing", state)
	}
	if user.IsRestricted() || user.IsDeactivated() {
		reason := "restricted"
		if user.IsDeactivated() {
			reason = "deactivated"
		}
		return true, importer.cleanupUnavailableUser(userID, reason, state)
	}
	return false, nil
}

func (importer *TitanicImporter) reconcileBeatmapScores(beatmapID int, remoteScoreIDs []int, state *common.State) error {
	deleted, err := services.DeleteScoresByBeatmapExcept(beatmapID, remoteScoreIDs, state)
	if err != nil {
		state.Logger.Logf("Failed to cleanup stale scores for beatmap %d: %v", beatmapID, err)
		return err
	}
	if deleted == 0 {
		return nil
	}

	state.Logger.Logf("Removed %d stale local scores for beatmap %d", deleted, beatmapID)
	return importer.recomputeRatingsAfterCleanup(state)
}

func (importer *TitanicImporter) reconcileUserTopPlays(userID int, remoteScoreIDs []int, state *common.State) error {
	deleted, err := services.DeleteScoresByPlayerExcept(userID, remoteScoreIDs, state)
	if err != nil {
		state.Logger.Logf("Failed to cleanup stale top plays for user %d: %v", userID, err)
		return err
	}
	if deleted == 0 {
		return nil
	}

	state.Logger.Logf("Removed %d stale local top plays for user %d", deleted, userID)
	return importer.recomputeRatingsAfterCleanup(state)
}

func (importer *TitanicImporter) recomputeRatingsAfterCleanup(state *common.State) error {
	if err := updaters.UpdatePlayerRatings(state); err != nil {
		state.Logger.Logf("Failed to recompute player ratings after cleanup: %v", err)
		return fmt.Errorf("failed to recompute player ratings after cleanup: %w", err)
	}
	return nil
}
