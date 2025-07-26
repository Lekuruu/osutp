package tp

import osu "github.com/natsukagami/go-osu-parser"

// DifficultyCalculationResult represents the computed difficulty and star ratings of a beatmap.
type DifficultyCalculationResult struct {
	AmountNormal      int     `json:"amountNormal"`
	AmountSliders     int     `json:"amountSliders"`
	AmountSpinners    int     `json:"amountSpinners"`
	MaxCombo          int     `json:"maxCombo"`
	SpeedDifficulty   float64 `json:"speedDifficulty"`
	AimDifficulty     float64 `json:"aimDifficulty"`
	SpeedStars        float64 `json:"speedStars"`
	AimStars          float64 `json:"aimStars"`
	StarRating        float64 `json:"starRating"`
	ApproachRate      float32 `json:"approachRate"`
	CircleSize        float32 `json:"circleSize"`
	HpDrainRate       float32 `json:"hpDrainRate"`
	OverallDifficulty float32 `json:"overallDifficulty"`
	SliderMultiplier  float64 `json:"sliderMultiplier"`
	SliderTickRate    float64 `json:"sliderTickRate"`
}

type DifficultyCalculationRequest struct {
	Beatmap    *BeatmapBase    `json:"beatmap"`
	HitObjects []HitObjectBase `json:"hitObjects"`
}

func NewDifficultyCalculationRequest(beatmap *BeatmapBase, hitObjects []HitObjectBase) *DifficultyCalculationRequest {
	return &DifficultyCalculationRequest{
		Beatmap:    beatmap,
		HitObjects: hitObjects,
	}
}

func NewDifficultyCalculationRequestFromBeatmap(beatmap osu.Beatmap) *DifficultyCalculationRequest {
	hitObjects := make([]HitObjectBase, len(beatmap.HitObjects))
	for i, obj := range beatmap.HitObjects {
		hitObjects[i] = *NewHitObjectBase(obj)
	}
	return &DifficultyCalculationRequest{
		Beatmap:    NewBeatmapBase(&beatmap),
		HitObjects: hitObjects,
	}
}
