package tp

import (
	"strings"

	osu "github.com/natsukagami/go-osu-parser"
)

// BeatmapBase holds difficulty settings and metadata for a beatmap.
type BeatmapBase struct {
	// Metadata
	Artist        string `json:"artist"`
	ArtistUnicode string `json:"artistUnicode"`
	Tags          string `json:"tags"`
	Title         string `json:"title"`
	TitleUnicode  string `json:"titleUnicode"`

	// Difficulty Settings
	DifficultyApproachRate     float64 `json:"difficultyApproachRate"`
	DifficultyCircleSize       float64 `json:"difficultyCircleSize"`
	DifficultyHpDrainRate      float64 `json:"difficultyHpDrainRate"`
	DifficultyOverall          float64 `json:"difficultyOverall"`
	DifficultySliderMultiplier float64 `json:"difficultySliderMultiplier"`
	DifficultySliderTickRate   float64 `json:"difficultySliderTickRate"`
	DifficultyMaxCombo         int     `json:"difficultyMaxCombo"`
}

func NewBeatmapBase(beatmap *osu.Beatmap) *BeatmapBase {
	if beatmap.ApproachRate <= 0 {
		// Some older beatmaps don't have AR set, so we set it to OD instead
		beatmap.ApproachRate = beatmap.OverallDifficulty
	}

	return &BeatmapBase{
		Artist:                     beatmap.Artist,
		ArtistUnicode:              beatmap.ArtistUnicode,
		Tags:                       strings.Join(beatmap.Tags, " "),
		Title:                      beatmap.Title,
		TitleUnicode:               beatmap.TitleUnicode,
		DifficultyApproachRate:     beatmap.ApproachRate,
		DifficultyCircleSize:       beatmap.CircleSize,
		DifficultyHpDrainRate:      beatmap.HPDrainRate,
		DifficultyOverall:          beatmap.OverallDifficulty,
		DifficultySliderMultiplier: beatmap.SliderMultiplier,
		DifficultySliderTickRate:   float64(beatmap.SliderTickRate),
		DifficultyMaxCombo:         beatmap.MaxCombo,
	}
}
