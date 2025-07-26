package tp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	osr "github.com/robloxxa/go-osr"
)

// Score represents a score record for a beatmap.
type Score struct {
	BeatmapFilename string `json:"beatmapFilename"`
	BeatmapChecksum string `json:"beatmapChecksum"`
	TotalScore      int    `json:"totalScore"`
	MaxCombo        int    `json:"maxCombo"`
	Amount300       int    `json:"amount300"`
	Amount100       int    `json:"amount100"`
	Amount50        int    `json:"amount50"`
	AmountMiss      int    `json:"amountMiss"`
	AmountGeki      int    `json:"amountGeki"`
	AmountKatu      int    `json:"amountKatu"`
	Mods            int    `json:"mods"`
}

func NewScoreFromReplay(replay *osr.Replay, beatmapFilename string) *Score {
	return &Score{
		BeatmapFilename: beatmapFilename,
		BeatmapChecksum: replay.BeatmapMD5,
		TotalScore:      int(replay.TotalScore),
		MaxCombo:        int(replay.Combo),
		Amount300:       int(replay.Count300),
		Amount100:       int(replay.Count100),
		Amount50:        int(replay.Count50),
		AmountMiss:      int(replay.CountMiss),
		AmountGeki:      int(replay.CountGeki),
		AmountKatu:      int(replay.CountKatu),
		Mods:            int(replay.Mods),
	}
}

// PerformanceCalculationResult represents the computed performance of a score
type PerformanceCalculationResult struct {
	Total float64 `json:"total"`
	Speed float64 `json:"speed"`
	Aim   float64 `json:"aim"`
	Acc   float64 `json:"accuracy"`
}

// PerformanceCalculationRequest represents a request to calculate the performance of a score.
type PerformanceCalculationRequest struct {
	Score      *Score                       `json:"score"`
	Difficulty *DifficultyCalculationResult `json:"difficulty"`
}

func (request *PerformanceCalculationRequest) Perform(serviceUrl string) (*PerformanceCalculationResult, error) {
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	jsonReader := strings.NewReader(string(jsonData))
	resp, err := http.Post(serviceUrl+"/performance", "application/json", jsonReader)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to calculate performance: %s", resp.Status)
	}

	var result PerformanceCalculationResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return &result, nil
}

func NewPerformanceCalculationRequest(score *Score, difficulty *DifficultyCalculationResult) *PerformanceCalculationRequest {
	return &PerformanceCalculationRequest{
		Score:      score,
		Difficulty: difficulty,
	}
}
