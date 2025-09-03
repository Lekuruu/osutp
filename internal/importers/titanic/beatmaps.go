package titanic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Lekuruu/osutp-web/internal/common"
	"github.com/Lekuruu/osutp-web/internal/database"
	"github.com/Lekuruu/osutp-web/internal/services"
)

type BeatmapSearchRequest struct {
	Language   *int    `json:"language,omitempty"`
	Genre      *int    `json:"genre,omitempty"`
	Mode       *int    `json:"mode,omitempty"`
	Uncleared  *bool   `json:"uncleared,omitempty"`
	Unplayed   *bool   `json:"unplayed,omitempty"`
	Cleared    *bool   `json:"cleared,omitempty"`
	Played     *bool   `json:"played,omitempty"`
	Query      *string `json:"query,omitempty"`
	Category   int     `json:"category"`
	Order      int     `json:"order"`
	Sort       int     `json:"sort"`
	Storyboard bool    `json:"storyboard"`
	Video      bool    `json:"video"`
	Titanic    bool    `json:"titanic"`
	Page       int     `json:"page"`
}

type BeatmapsetModel struct {
	ID                 int            `json:"id"`
	Title              *string        `json:"title,omitempty"`
	Artist             *string        `json:"artist,omitempty"`
	Creator            *string        `json:"creator,omitempty"`
	Source             *string        `json:"source,omitempty"`
	Tags               *string        `json:"tags,omitempty"`
	CreatorID          *int           `json:"creator_id,omitempty"`
	Status             int            `json:"status"`
	HasVideo           bool           `json:"has_video"`
	HasStoryboard      bool           `json:"has_storyboard"`
	Server             int            `json:"server"`
	Available          bool           `json:"available"`
	Enhanced           bool           `json:"enhanced"`
	CreatedAt          string         `json:"created_at"`
	ApprovedAt         *string        `json:"approved_at,omitempty"`
	LastUpdate         string         `json:"last_update"`
	OszFilesize        int            `json:"osz_filesize"`
	OszFilesizeNoVideo int            `json:"osz_filesize_novideo"`
	DisplayTitle       string         `json:"display_title"`
	LanguageID         int            `json:"language_id"`
	GenreID            int            `json:"genre_id"`
	Beatmaps           []BeatmapModel `json:"beatmaps"`
}

type BeatmapModel struct {
	ID           int     `json:"id"`
	SetID        int     `json:"set_id"`
	Mode         int     `json:"mode"`
	MD5          string  `json:"md5"`
	Status       int     `json:"status"`
	Version      string  `json:"version"`
	Filename     string  `json:"filename"`
	CreatedAt    string  `json:"created_at"`
	LastUpdate   string  `json:"last_update"`
	Playcount    int     `json:"playcount"`
	Passcount    int     `json:"passcount"`
	TotalLength  int     `json:"total_length"`
	DrainLength  int     `json:"drain_length"`
	MaxCombo     int     `json:"max_combo"`
	BPM          float64 `json:"bpm"`
	CS           float64 `json:"cs"`
	AR           float64 `json:"ar"`
	OD           float64 `json:"od"`
	HP           float64 `json:"hp"`
	Diff         float64 `json:"diff"`
	CountNormal  int     `json:"count_normal"`
	CountSlider  int     `json:"count_slider"`
	CountSpinner int     `json:"count_spinner"`
}

func (beatmap *BeatmapModel) ToSchema(beatmapset *BeatmapsetModel) *database.Beatmap {
	return &database.Beatmap{
		ID:                   beatmap.ID,
		SetID:                beatmap.SetID,
		Title:                dereferenceString(beatmapset.Title),
		Artist:               dereferenceString(beatmapset.Artist),
		Creator:              dereferenceString(beatmapset.Creator),
		Source:               dereferenceString(beatmapset.Source),
		Tags:                 dereferenceString(beatmapset.Tags),
		Version:              beatmap.Version,
		Status:               beatmap.Status,
		AR:                   beatmap.AR,
		OD:                   beatmap.OD,
		CS:                   beatmap.CS,
		AmountNormal:         beatmap.CountNormal,
		AmountSliders:        beatmap.CountSlider,
		AmountSpinners:       beatmap.CountSpinner,
		MaxCombo:             beatmap.MaxCombo,
		DifficultyAttributes: database.DifficultyAttributes{},
	}
}

const TitanicApiBaseurl = "https://api.titanic.sh"

func ImportBeatmapsByDifficulty(page int, state *common.State) error {
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

	results, err := performSearchRequest(request)
	if err != nil {
		return err
	}

	for _, beatmapset := range results {
		for _, beatmap := range beatmapset.Beatmaps {
			if beatmap.Mode != GameModeOsu {
				continue
			}
			if exists, _ := services.BeatmapExists(beatmap.ID, state); exists {
				continue
			}
			schema := beatmap.ToSchema(&beatmapset)
			result := state.Database.Create(schema)
			if result.Error != nil {
				return result.Error
			}

			file, err := fetchBeatmapFile(beatmap.ID)
			if err != nil {
				return err
			}

			err = common.UpdateBeatmapDifficulty(file, schema, state)
			if err != nil {
				return err
			}
			fmt.Printf("Imported Beatmap: '%s' (https://osu.titanic.sh/b/%d)\n", schema.FullName(), schema.ID)
		}
	}
	return nil
}

func performSearchRequest(request BeatmapSearchRequest) ([]BeatmapsetModel, error) {
	jsonData, _ := json.Marshal(request)
	url := TitanicApiBaseurl + "/beatmapsets/search"

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
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

func fetchBeatmapFile(beatmapId int) ([]byte, error) {
	url := fmt.Sprintf("%s/beatmaps/%d/file", TitanicApiBaseurl, beatmapId)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func dereferenceString(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}
