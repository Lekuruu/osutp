package titanic

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Lekuruu/osutp/internal/common"
	"github.com/Lekuruu/osutp/internal/database"
	"github.com/Lekuruu/osutp/internal/services"
	"github.com/Lekuruu/osutp/pkg/tp"
	"gorm.io/gorm"
)

func (importer *TitanicImporter) ImportScore(scoreId int, state *common.State) (*database.Score, error) {
	url := fmt.Sprintf("%s/scores/%d", importer.ApiUrl, scoreId)
	resp, err := http.Get(url)
	if err != nil {
		// Check for any rate limit errors and wait if needed
		if strings.Contains(err.Error(), "429 Too Many Requests") {
			time.Sleep(time.Second * 60)
			return importer.ImportScore(scoreId, state)
		}
		return nil, err
	}
	defer resp.Body.Close()

	var score ScoreModel
	if err := json.NewDecoder(resp.Body).Decode(&score); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	beatmap, err := services.FetchBeatmapById(score.BeatmapID, state)
	if err != nil {
		return nil, err
	}
	if beatmap == nil {
		// Try to import the beatmap if it doesn't exist
		beatmap, err = importer.ImportBeatmap(score.BeatmapID, true, state)
		if err != nil {
			return nil, fmt.Errorf("failed to import beatmap %d for score %d: %v", score.BeatmapID, score.ID, err)
		}
	}

	return importer.importScoreFromModel(score, beatmap, state)
}

func (importer *TitanicImporter) importBeatmapLeaderboard(beatmapId int, state *common.State) error {
	scores, err := importer.performLeaderboardRequest(beatmapId, 0, state)
	if err != nil {
		return err
	}

	beatmap, err := services.FetchBeatmapById(beatmapId, state)
	if err != nil {
		return err
	}
	if beatmap == nil {
		return fmt.Errorf("beatmap with id %d not found in database", beatmapId)
	}

	for _, score := range scores {
		importer.importScoreFromModel(score, beatmap, state)
	}
	return nil
}

func (importer *TitanicImporter) performLeaderboardRequest(beatmapId int, offset int, state *common.State) ([]ScoreModel, error) {
	url := fmt.Sprintf("%s/beatmaps/%d/scores?offset=%d", importer.ApiUrl, beatmapId, offset)
	resp, err := http.Get(url)
	if err != nil {
		// Check for any rate limit errors and wait if needed
		if strings.Contains(err.Error(), "429 Too Many Requests") {
			time.Sleep(time.Second * 60)
			return importer.performLeaderboardRequest(beatmapId, offset, state)
		}
		return nil, err
	}
	defer resp.Body.Close()

	var scores ScoreCollectionModel
	if err := json.NewDecoder(resp.Body).Decode(&scores); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return scores.Scores, nil
}

func (importer *TitanicImporter) importScoreFromModel(score ScoreModel, beatmap *database.Beatmap, state *common.State) (*database.Score, error) {
	if score.Mode != 0 {
		// Only process osu! standard scores
		return nil, nil
	}
	if score.User.Restricted || !score.User.Activated {
		// Skip & delete restricted users
		services.DeleteScoresByPlayer(score.UserID, state)
		services.DeletePlayer(score.UserID, state)
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
