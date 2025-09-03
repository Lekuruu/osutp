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

const (
	GameModeOsu   = 0
	GameModeTaiko = 1
	GameModeCatch = 2
	GameModeMania = 3
)

const (
	BeatmapCategoryAny         = 0
	BeatmapCategoryLeaderboard = 1
	BeatmapCategoryRanked      = 2
	BeatmapCategoryQualified   = 3
	BeatmapCategoryLoved       = 4
	BeatmapCategoryApproved    = 5
	BeatmapCategoryPending     = 6
	BeatmapCategoryWIP         = 7
	BeatmapCategoryGraveyard   = 8
)

const (
	BeatmapSortByTitle      = 0
	BeatmapSortByArtist     = 1
	BeatmapSortByCreator    = 2
	BeatmapSortByDifficulty = 3
	BeatmapSortByRanked     = 4
	BeatmapSortByRating     = 5
	BeatmapSortByPlays      = 6
	BeatmapSortByCreated    = 7
	BeatmapSortByRelevance  = 8
)

const (
	BeatmapOrderDescending = 0
	BeatmapOrderAscending  = 1
)

const (
	BeatmapLanguageAny          = 0
	BeatmapLanguageUnspecified  = 1
	BeatmapLanguageEnglish      = 2
	BeatmapLanguageJapanese     = 3
	BeatmapLanguageChinese      = 4
	BeatmapLanguageInstrumental = 5
	BeatmapLanguageKorean       = 6
	BeatmapLanguageFrench       = 7
	BeatmapLanguageGerman       = 8
	BeatmapLanguageSwedish      = 9
	BeatmapLanguageSpanish      = 10
	BeatmapLanguageItalian      = 11
	BeatmapLanguageRussian      = 12
	BeatmapLanguagePolish       = 13
	BeatmapLanguageOther        = 14
)

const (
	BeatmapGenreAny         = 0
	BeatmapGenreUnspecified = 1
	BeatmapGenreVideoGame   = 2
	BeatmapGenreAnime       = 3
	BeatmapGenreRock        = 4
	BeatmapGenrePop         = 5
	BeatmapGenreOther       = 6
	BeatmapGenreNovelty     = 7
	BeatmapGenreHipHop      = 9
	BeatmapGenreElectronic  = 10
	BeatmapGenreMetal       = 11
	BeatmapGenreClassical   = 12
	BeatmapGenreFolk        = 13
	BeatmapGenreJazz        = 14
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

// BeatmapsetModel represents a beatmapset with its metadata
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

// BeatmapModel represents an individual beatmap inside a set
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

func dereferenceString(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}
