package titanic

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Lekuruu/osutp/internal/common"
	"github.com/Lekuruu/osutp/internal/database"
	"github.com/Lekuruu/osutp/internal/services"
	"github.com/Lekuruu/osutp/internal/updaters"
)

func (importer *TitanicImporter) ImportBeatmapset(beatmapsetID int, importLeaderboard bool, state *common.State) ([]*database.Beatmap, error) {
	beatmapset, err := importer.fetchBeatmapsetById(beatmapsetID)
	if err != nil {
		return nil, err
	}
	if beatmapset == nil {
		if err := importer.cleanupMissingBeatmapset(beatmapsetID, state); err != nil {
			return nil, err
		}
		return nil, nil
	}

	var importedBeatmaps []*database.Beatmap
	remoteBeatmapIDs := make(map[int]struct{}, len(beatmapset.Beatmaps))

	for _, beatmap := range beatmapset.Beatmaps {
		remoteBeatmapIDs[beatmap.ID] = struct{}{}
		beatmap.Beatmapset = beatmapset
		beatmapObject, err := importer.importBeatmapFromModel(&beatmap, true, state)
		if err != nil {
			state.Logger.Logf("Error importing beatmap %d: %v", beatmap.ID, err)
			continue
		}
		importedBeatmaps = append(importedBeatmaps, beatmapObject)

		if !importLeaderboard {
			continue
		}
		err = importer.importBeatmapLeaderboard(beatmap.ID, state)
		if err != nil {
			state.Logger.Logf("Error importing leaderboard for beatmap %d: %v", beatmap.ID, err)
			continue
		}
	}

	if err := importer.cleanupBeatmapsetRemovedBeatmaps(beatmapsetID, remoteBeatmapIDs, state); err != nil {
		return importedBeatmaps, err
	}
	return importedBeatmaps, nil
}

func (importer *TitanicImporter) ImportBeatmap(beatmapID int, importLeaderboard bool, state *common.State) (*database.Beatmap, error) {
	beatmap, err := importer.fetchBeatmapById(beatmapID)
	if err != nil {
		return nil, err
	}
	if beatmap == nil {
		if err := importer.cleanupMissingBeatmap(beatmapID, state); err != nil {
			return nil, err
		}
		return nil, nil
	}

	beatmapObject, err := importer.importBeatmapFromModel(beatmap, true, state)
	if err != nil {
		return nil, err
	}
	if !importLeaderboard {
		return beatmapObject, nil
	}

	err = importer.importBeatmapLeaderboard(beatmapID, state)
	if err != nil {
		return nil, err
	}
	return beatmapObject, nil
}

func (importer *TitanicImporter) ImportBeatmapsByDifficulty(page int, state *common.State) (int, error) {
	mode := GameModeOsu
	request := BeatmapSearchRequest{
		Category:   BeatmapCategoryLeaderboard,
		Order:      BeatmapOrderDescending,
		Sort:       BeatmapSortByDifficulty,
		Storyboard: false,
		Video:      false,
		Titanic:    false,
		Mode:       &mode,
		Page:       page,
	}

	results, err := importer.performSearchRequest(request, state)
	if err != nil {
		return 0, err
	}

	for _, beatmapset := range results {
		for _, beatmap := range beatmapset.Beatmaps {
			beatmap.Beatmapset = &beatmapset
			beatmapObject, err := importer.importBeatmapFromModel(&beatmap, false, state)
			if err != nil {
				state.Logger.Logf("Error importing beatmap %d: %v", beatmap.ID, err)
				continue
			}
			if beatmapObject == nil {
				continue
			}

			err = importer.importBeatmapLeaderboard(beatmap.ID, state)
			if err != nil {
				state.Logger.Logf("Error importing leaderboard for beatmap %d: %v", beatmap.ID, err)
				continue
			}
		}
	}

	state.Logger.Logf("Imported %d beatmaps from page %d", len(results), page)
	return len(results), nil
}

func (importer *TitanicImporter) ImportBeatmapsByDate(page int, state *common.State) (int, error) {
	mode := GameModeOsu
	request := BeatmapSearchRequest{
		Category:   BeatmapCategoryLeaderboard,
		Order:      BeatmapOrderDescending,
		Sort:       BeatmapSortByCreated,
		Storyboard: false,
		Video:      false,
		Titanic:    false,
		Mode:       &mode,
		Page:       page,
	}

	results, err := importer.performSearchRequest(request, state)
	if err != nil {
		return 0, err
	}

	for _, beatmapset := range results {
		for _, beatmap := range beatmapset.Beatmaps {
			beatmap.Beatmapset = &beatmapset
			beatmapObject, err := importer.importBeatmapFromModel(&beatmap, false, state)
			if err != nil {
				state.Logger.Logf("Error importing beatmap %d: %v", beatmap.ID, err)
				continue
			}
			if beatmapObject == nil {
				continue
			}

			err = importer.importBeatmapLeaderboard(beatmap.ID, state)
			if err != nil {
				state.Logger.Logf("Error importing leaderboard for beatmap %d: %v", beatmap.ID, err)
				continue
			}
		}
	}

	state.Logger.Logf("Imported %d beatmaps from page %d", len(results), page)
	return len(results), nil
}

func (importer *TitanicImporter) importBeatmapFromModel(beatmap *BeatmapModel, forcedRecalculation bool, state *common.State) (*database.Beatmap, error) {
	if beatmap.Mode != GameModeOsu {
		// We only support osu!standard for now
		return nil, nil
	}

	if beatmap.Beatmapset == nil {
		beatmapset, err := importer.fetchBeatmapsetById(beatmap.SetID)
		if err != nil {
			return nil, err
		}
		if beatmapset == nil {
			beatmapset = &BeatmapsetModel{}
		}
		beatmap.Beatmapset = beatmapset
	}

	// Check for existing beatmap entry
	beatmapEntry, err := services.FetchBeatmapById(beatmap.ID, state)
	if err != nil && err.Error() != "record not found" {
		return nil, err
	}

	if beatmapEntry == nil {
		// Create new beatmap entry if it doesn't exist
		beatmapEntry = beatmap.ToSchema(beatmap.Beatmapset)
		err := services.CreateBeatmap(beatmapEntry, state)
		if err != nil {
			return nil, err
		}
	}

	// TODO: Update beatmap metadata if beatmap was updated
	//		 We could do this by checking for the hash of the beatmap file

	// For now, we just update the status
	err = services.UpdateBeatmapStatus(beatmap.ID, beatmap.Status, state)
	if err != nil {
		return nil, err
	}

	if beatmapEntry.HasDifficultyAttributes() && !forcedRecalculation {
		// Skip if difficulty attributes already exist and recalculation is not forced
		state.Logger.Logf("Skipping difficulty calculation for Beatmap: '%s' (%s/b/%d)", beatmapEntry.FullName(), importer.WebUrl, beatmapEntry.ID)
		return beatmapEntry, nil
	}

	file, err := importer.fetchBeatmapFile(beatmap.ID)
	if err != nil {
		return nil, err
	}

	err = updaters.UpdateBeatmapDifficulty(file, beatmapEntry, state)
	if err != nil {
		return nil, err
	}

	state.Logger.Logf("Imported Beatmap: '%s' (%s/b/%d)", beatmapEntry.FullName(), importer.WebUrl, beatmapEntry.ID)
	return beatmapEntry, nil
}

func (importer *TitanicImporter) performSearchRequest(request BeatmapSearchRequest, state *common.State) ([]BeatmapsetModel, error) {
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	url := importer.ApiUrl + "/beatmapsets/search"
	var results []BeatmapsetModel
	if err := importer.PostJson(url, jsonData, &results); err != nil {
		return nil, err
	}
	return results, nil
}

func (importer *TitanicImporter) fetchBeatmapById(beatmapId int) (*BeatmapModel, error) {
	url := fmt.Sprintf("%s/beatmaps/%d", importer.ApiUrl, beatmapId)
	var beatmap BeatmapModel
	if err := importer.GetJson(url, &beatmap); err != nil {
		var statusErr *HttpStatusError
		if errors.As(err, &statusErr) && statusErr.statusCode == http.StatusNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &beatmap, nil
}

func (importer *TitanicImporter) fetchBeatmapsetById(beatmapsetId int) (*BeatmapsetModel, error) {
	url := fmt.Sprintf("%s/beatmapsets/%d", importer.ApiUrl, beatmapsetId)
	var beatmapset BeatmapsetModel
	if err := importer.GetJson(url, &beatmapset); err != nil {
		var statusErr *HttpStatusError
		if errors.As(err, &statusErr) && statusErr.statusCode == http.StatusNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &beatmapset, nil
}

func (importer *TitanicImporter) fetchBeatmapFile(beatmapId int) ([]byte, error) {
	url := fmt.Sprintf("%s/beatmaps/%d/file", importer.ApiUrl, beatmapId)
	return importer.GetBytes(url)
}

func dereferenceString(s *string) (result string) {
	if s != nil {
		return *s
	}
	return ""
}
