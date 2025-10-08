package titanic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/Lekuruu/osutp/internal/common"
	"github.com/Lekuruu/osutp/internal/database"
	"github.com/Lekuruu/osutp/internal/services"
	"github.com/Lekuruu/osutp/internal/updaters"
)

func (importer *TitanicImporter) ImportBeatmap(beatmapID int, importLeaderboard bool, state *common.State) (*database.Beatmap, error) {
	beatmap, err := importer.fetchBeatmapById(beatmapID)
	if err != nil {
		// TODO: Delete existing beatmap, if it exists
		return nil, err
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

	if beatmapEntry.DifficultyAttributes != nil && !forcedRecalculation {
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
	jsonData, _ := json.Marshal(request)
	url := importer.ApiUrl + "/beatmapsets/search"

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		// Check for any rate limit errors and wait if needed
		if strings.Contains(err.Error(), "429 Too Many Requests") {
			time.Sleep(time.Second * 60)
			return importer.performSearchRequest(request, state)
		}
		return nil, err
	}
	defer resp.Body.Close()

	// Read response
	body, _ := io.ReadAll(resp.Body)

	var results []BeatmapsetModel
	if err := json.Unmarshal(body, &results); err != nil {
		return nil, err
	}
	return results, nil
}

func (importer *TitanicImporter) fetchBeatmapById(beatmapId int) (*BeatmapModel, error) {
	url := fmt.Sprintf("%s/beatmaps/%d", importer.ApiUrl, beatmapId)
	resp, err := http.Get(url)
	if err != nil {
		// Check for any rate limit errors and wait if needed
		if strings.Contains(err.Error(), "429 Too Many Requests") {
			time.Sleep(time.Second * 60)
			return importer.fetchBeatmapById(beatmapId)
		}
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var beatmap BeatmapModel
	if err := json.Unmarshal(body, &beatmap); err != nil {
		return nil, err
	}
	return &beatmap, nil
}

func (importer *TitanicImporter) fetchBeatmapFile(beatmapId int) ([]byte, error) {
	url := fmt.Sprintf("%s/beatmaps/%d/file", importer.ApiUrl, beatmapId)
	resp, err := http.Get(url)
	if err != nil {
		// Check for any rate limit errors and wait if needed
		if strings.Contains(err.Error(), "429 Too Many Requests") {
			time.Sleep(time.Second * 60)
			return importer.fetchBeatmapFile(beatmapId)
		}
		return nil, err
	}

	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func dereferenceString(s *string) (result string) {
	if s != nil {
		return *s
	}
	return ""
}
