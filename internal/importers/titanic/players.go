package titanic

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Lekuruu/osutp/internal/common"
	"github.com/Lekuruu/osutp/internal/database"
	"github.com/Lekuruu/osutp/internal/services"
	"gorm.io/gorm"
)

func (importer *TitanicImporter) ImportUser(userID int, state *common.State) (*database.Player, error) {
	user, err := importer.fetchUserById(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, nil
	}
	if user.Restricted || !user.Activated {
		// Skip & delete restricted users
		services.DeleteScoresByPlayer(userID, state)
		services.DeletePlayer(userID, state)
		return nil, nil
	}

	return importer.importUserFromModel(*user, state)
}

func (importer *TitanicImporter) ImportUsersFromRankings(page int, state *common.State) (int, error) {
	return 0, fmt.Errorf("not implemented")
}

func (importer *TitanicImporter) importUserFromModel(user UserModel, state *common.State) (*database.Player, error) {
	// Check for existing player entry
	userEntry, err := services.FetchPlayerById(user.ID, state)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if userEntry != nil {
		importer.importUserTopPlays(user, state)
		return userEntry, nil
	}

	userEntry = user.ToSchema()
	if err := services.CreatePlayer(userEntry, state); err != nil {
		return nil, err
	}

	importer.importUserTopPlays(user, state)
	return userEntry, nil
}

func (importer *TitanicImporter) importUserTopPlays(user UserModel, state *common.State) error {
	offset := 0
	limit := 50

	for {
		scores, err := importer.fetchUserTopPlays(user.ID, 0, offset, limit)
		if err != nil {
			return fmt.Errorf("failed to fetch top plays for user %d: %v", user.ID, err)
		}

		if len(scores.Scores) == 0 {
			// No more scores to import
			break
		}

		for _, score := range scores.Scores {
			beatmap, err := services.FetchBeatmapById(score.BeatmapID, state)
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				continue
			}

			if beatmap == nil {
				// Try to import the beatmap if it doesn't exist
				beatmap, err = importer.ImportBeatmap(score.BeatmapID, false, state)
				if err != nil {
					continue
				}
			}

			_, err = importer.importScoreFromModel(score, beatmap, state)
			if err != nil {
				state.Logger.Logf("Failed to import score %d: %v", score.ID, err)
				continue
			}
		}

		// Check if we have more scores to fetch
		if len(scores.Scores) < limit {
			break
		}

		offset += limit
	}

	return nil
}

func (importer *TitanicImporter) fetchUserTopPlays(userID int, mode int, offset int, limit int) (*ScoreCollectionModel, error) {
	url := fmt.Sprintf("%s/users/%d/top/%d?offset=%d&limit=%d", importer.ApiUrl, userID, mode, offset, limit)
	var scores ScoreCollectionModel
	if err := importer.GetJson(url, &scores); err != nil {
		return nil, err
	}

	return &scores, nil
}

func (importer *TitanicImporter) fetchUserById(userID int) (*UserModel, error) {
	url := fmt.Sprintf("%s/users/%d", importer.ApiUrl, userID)
	var user UserModel
	if err := importer.GetJson(url, &user); err != nil {
		var statusErr *HttpStatusError
		if errors.As(err, &statusErr) && statusErr.statusCode == http.StatusNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}
