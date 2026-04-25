package titanic

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Lekuruu/osutp/internal/common"
	"github.com/Lekuruu/osutp/internal/database"
	"github.com/Lekuruu/osutp/internal/services"
	"github.com/Lekuruu/osutp/pkg/tp"
	"gorm.io/gorm"
)

func (importer *TitanicImporter) ImportScore(scoreId int, state *common.State) (*database.Score, error) {
	url := fmt.Sprintf("%s/scores/%d", importer.ApiUrl, scoreId)
	var score ScoreModel
	if err := importer.GetJson(url, &score); err != nil {
		var statusErr *HttpStatusError
		if errors.As(err, &statusErr) && statusErr.statusCode == http.StatusNotFound {
			deleted, err := services.DeleteScoreWithCount(scoreId, state)
			if err != nil {
				state.Logger.Logf("Failed to cleanup missing score %d: %v", scoreId, err)
				return nil, err
			}
			if deleted > 0 {
				state.Logger.Logf("Cleaned up local score %d after confirmed remote 404", scoreId)
				if err := importer.recomputeRatingsAfterCleanup(state); err != nil {
					return nil, err
				}
			}
			return nil, nil
		}
		return nil, err
	}

	beatmap, err := services.FetchBeatmapById(score.BeatmapID, state)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if beatmap == nil {
		// Try to import the beatmap if it doesn't exist
		beatmap, err = importer.ImportBeatmap(score.BeatmapID, true, state)
		if err != nil {
			return nil, fmt.Errorf("failed to import beatmap %d for score %d: %v", score.BeatmapID, score.ID, err)
		}
		if beatmap == nil {
			return nil, nil
		}
	}

	return importer.importScoreFromModel(score, beatmap, state)
}

func (importer *TitanicImporter) importBeatmapLeaderboard(beatmapId int, state *common.State) error {
	beatmap, err := services.FetchBeatmapById(beatmapId, state)
	if err != nil {
		return err
	}
	if beatmap == nil {
		return fmt.Errorf("beatmap with id %d not found in database", beatmapId)
	}

	offset := 0
	limit := 100
	remoteScoreIDs := make([]int, 0)
	importErrors := make([]error, 0)

	for {
		scores, err := importer.performLeaderboardRequest(beatmapId, offset, limit, state)
		if err != nil {
			return err
		}

		if len(scores) == 0 {
			break
		}

		for _, score := range scores {
			remoteScoreIDs = append(remoteScoreIDs, score.ID)
			if _, err := importer.importScoreFromModel(score, beatmap, state); err != nil {
				state.Logger.Logf("Failed to import leaderboard score %d for beatmap %d: %v", score.ID, beatmapId, err)
				importErrors = append(importErrors, err)
			}
		}

		if len(scores) < limit {
			break
		}
		offset += limit
	}

	if err := importer.reconcileBeatmapScores(beatmap.ID, remoteScoreIDs, state); err != nil {
		return err
	}

	if len(importErrors) > 0 {
		return fmt.Errorf("failed to import %d scores for beatmap %d: %w", len(importErrors), beatmapId, errors.Join(importErrors...))
	}

	if err := services.UpdateBeatmapLastScoreUpdate(beatmap.ID, time.Now(), state); err != nil {
		return err
	}
	return nil
}

func (importer *TitanicImporter) performLeaderboardRequest(beatmapId int, offset int, limit int, state *common.State) ([]ScoreModel, error) {
	url := fmt.Sprintf("%s/beatmaps/%d/scores?offset=%d&limit=%d", importer.ApiUrl, beatmapId, offset, limit)
	var scores ScoreCollectionModel
	if err := importer.GetJson(url, &scores); err != nil {
		return nil, err
	}

	return scores.Scores, nil
}

func (importer *TitanicImporter) importScoreFromModel(score ScoreModel, beatmap *database.Beatmap, state *common.State) (*database.Score, error) {
	if score.Mode != 0 {
		// Only process osu! standard scores
		return nil, nil
	}
	if score.User.Restricted || !score.User.Activated {
		// Skip & delete restricted/deactivated users
		reason := "restricted"
		if !score.User.Activated {
			reason = "deactivated"
		}
		if err := importer.cleanupUnavailableUser(score.UserID, reason, state); err != nil {
			return nil, err
		}
		return nil, nil
	}

	scoreExists, err := services.ScoreExists(score.ID, state)
	if err != nil {
		return nil, fmt.Errorf("failed to check if score exists: %v", err)
	}
	if scoreExists {
		return nil, nil
	}

	schema := score.ToSchema()
	if schema.Relaxing() {
		// Skip relax/autopilot scores
		return nil, nil
	}

	difficulty, err := beatmap.DifficultyCalculationResult(schema.DifficultyMods())
	if err != nil {
		// Beatmap most likely has no difficulty attributes so we try to update it
		importer.ImportBeatmap(beatmap.ID, false, state)
		return nil, fmt.Errorf("failed to get difficulty calculation result: %v", err)
	}

	tpScore := schema.CalculationRequest(difficulty)
	result := tp.CalculatePerformance(difficulty, tpScore)
	if result == nil {
		return nil, fmt.Errorf("failed to calculate performance")
	}

	schema.TotalTp = result.Total
	schema.AimTp = result.Aim
	schema.SpeedTp = result.Speed
	schema.AccTp = result.Acc

	user, err := services.FetchPlayerById(score.UserID, state)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("failed to fetch user: %v", err)
	}

	if user == nil {
		user, err = importer.ImportUser(score.UserID, state)
		if err != nil {
			return nil, fmt.Errorf("failed to import user %d: %v", score.UserID, err)
		}
		if user == nil {
			// User does not exist or is restricted.
			return nil, nil
		}
	}

	personalBest, err := services.FetchPersonalBestScore(user.ID, beatmap.ID, state)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("failed to fetch personal best score: %v", err)
	}

	if personalBest != nil && personalBest.TotalTp > schema.TotalTp {
		// We don't have a new personal best
		return nil, nil
	}

	// Delete old personal best
	if personalBest != nil {
		if err := services.DeleteScore(personalBest.ID, state); err != nil {
			return nil, fmt.Errorf("failed to delete old personal best score: %v", err)
		}
	}

	if err := services.CreateScore(schema, state); err != nil {
		return nil, fmt.Errorf("failed to create score: %v", err)
	}

	return schema, nil
}
