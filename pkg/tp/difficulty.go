package tp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	osu "github.com/natsukagami/go-osu-parser"
)

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

// DifficultyCalculationRequest represents a request to calculate the difficulty of a beatmap.
type DifficultyCalculationRequest struct {
	Beatmap    *BeatmapBase    `json:"beatmap"`
	HitObjects []HitObjectBase `json:"hitObjects"`
	Mods       int             `json:"mods"`
}

func (request *DifficultyCalculationRequest) Perform(serviceUrl string) (*DifficultyCalculationResult, error) {
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	jsonReader := strings.NewReader(string(jsonData))
	resp, err := http.Post(serviceUrl+"/difficulty", "application/json", jsonReader)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to calculate difficulty: %s", resp.Status)
	}

	var result DifficultyCalculationResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	// NOTE: osu!tp.Service is not able to calculate this properly at the moment
	result.MaxCombo = request.Beatmap.DifficultyMaxCombo
	return &result, nil
}

func NewDifficultyCalculationRequest(beatmap *BeatmapBase, hitObjects []HitObjectBase) *DifficultyCalculationRequest {
	return &DifficultyCalculationRequest{
		Beatmap:    beatmap,
		HitObjects: hitObjects,
	}
}

func NewDifficultyCalculationRequestFromBeatmap(beatmap osu.Beatmap, mods int) *DifficultyCalculationRequest {
	hitObjects := make([]HitObjectBase, len(beatmap.HitObjects))
	for i, obj := range beatmap.HitObjects {
		hitObjects[i] = *NewHitObjectBase(obj)
	}
	return &DifficultyCalculationRequest{
		Beatmap:    NewBeatmapBase(&beatmap),
		HitObjects: hitObjects,
		Mods:       mods,
	}
}
