package titanic

import (
	"strings"
	"time"

	"github.com/Lekuruu/osutp/internal/database"
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
	ID           int              `json:"id"`
	SetID        int              `json:"set_id"`
	Mode         int              `json:"mode"`
	MD5          string           `json:"md5"`
	Status       int              `json:"status"`
	Version      string           `json:"version"`
	Filename     string           `json:"filename"`
	CreatedAt    string           `json:"created_at"`
	LastUpdate   string           `json:"last_update"`
	Playcount    int              `json:"playcount"`
	Passcount    int              `json:"passcount"`
	TotalLength  int              `json:"total_length"`
	DrainLength  int              `json:"drain_length"`
	MaxCombo     int              `json:"max_combo"`
	BPM          float64          `json:"bpm"`
	CS           float64          `json:"cs"`
	AR           float64          `json:"ar"`
	OD           float64          `json:"od"`
	HP           float64          `json:"hp"`
	Diff         float64          `json:"diff"`
	CountNormal  int              `json:"count_normal"`
	CountSlider  int              `json:"count_slider"`
	CountSpinner int              `json:"count_spinner"`
	Beatmapset   *BeatmapsetModel `json:"beatmapset,omitempty"`
}

func (beatmap *BeatmapModel) ToSchema(beatmapset *BeatmapsetModel) *database.Beatmap {
	createdAt, err := time.Parse("2006-01-02T15:04:05", beatmap.CreatedAt)
	if err != nil {
		createdAt = time.Now().UTC()
	}

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
		CreatedAt:            createdAt,
		DifficultyAttributes: database.DifficultyAttributes{},
	}
}

type ScoreCollectionModel struct {
	Total  int          `json:"total"`
	Scores []ScoreModel `json:"scores"`
}

type ScoreModel struct {
	ID            int       `json:"id"`
	UserID        int       `json:"user_id"`
	BeatmapID     int       `json:"beatmap_id"`
	SubmittedAt   string    `json:"submitted_at"`
	Mode          int       `json:"mode"`
	StatusPP      int       `json:"status_pp"`
	StatusScore   int       `json:"status_score"`
	ClientVersion int       `json:"client_version"`
	PP            float64   `json:"pp"`
	PPv1          float64   `json:"ppv1"`
	Acc           float64   `json:"acc"`
	TotalScore    int       `json:"total_score"`
	MaxCombo      int       `json:"max_combo"`
	Mods          int       `json:"mods"`
	Perfect       bool      `json:"perfect"`
	Passed        bool      `json:"passed"`
	Pinned        bool      `json:"pinned"`
	Count300      int       `json:"n300"`
	Count100      int       `json:"n100"`
	Count50       int       `json:"n50"`
	CountMiss     int       `json:"nMiss"`
	CountGeki     int       `json:"nGeki"`
	CountKatu     int       `json:"nKatu"`
	Grade         string    `json:"grade"`
	ReplayViews   int       `json:"replay_views"`
	Failtime      *int      `json:"failtime,omitempty"`
	User          UserModel `json:"user"`
}

func (score *ScoreModel) ToSchema() *database.Score {
	createdAt, err := time.Parse(time.RFC3339, score.SubmittedAt)
	if err != nil {
		createdAt = time.Now().UTC()
	}

	return &database.Score{
		ID:         score.ID,
		BeatmapID:  score.BeatmapID,
		PlayerID:   score.UserID,
		TotalScore: score.TotalScore,
		MaxCombo:   score.MaxCombo,
		Mods:       uint32(score.Mods),
		FullCombo:  score.Perfect,
		Grade:      score.Grade,
		Accuracy:   score.Acc,
		Amount300:  score.Count300,
		Amount100:  score.Count100,
		Amount50:   score.Count50,
		AmountGeki: score.CountGeki,
		AmountKatu: score.CountKatu,
		AmountMiss: score.CountMiss,
		CreatedAt:  createdAt,
	}
}

type UserModel struct {
	ID             int     `json:"id"`
	Name           string  `json:"name"`
	Country        string  `json:"country"`
	CreatedAt      string  `json:"created_at"`
	LatestActivity string  `json:"latest_activity"`
	Restricted     bool    `json:"restricted"`
	Activated      bool    `json:"activated"`
	PreferredMode  int     `json:"preferred_mode"`
	Playstyle      int     `json:"playstyle"`
	Banner         *string `json:"banner,omitempty"`
	Website        *string `json:"website,omitempty"`
	Discord        *string `json:"discord,omitempty"`
	Twitter        *string `json:"twitter,omitempty"`
	Location       *string `json:"location,omitempty"`
	Interests      *string `json:"interests,omitempty"`
}

func (user *UserModel) ToSchema() *database.Player {
	createdAt, err := time.Parse("2006-01-02T15:04:05", user.CreatedAt)
	if err != nil {
		createdAt = time.Now().UTC()
	}

	return &database.Player{
		ID:        user.ID,
		Name:      user.Name,
		CreatedAt: createdAt,
		Country:   strings.ToUpper(user.Country),
	}
}
