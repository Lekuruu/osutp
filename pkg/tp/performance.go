package tp

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

// PerformanceCalculationResult represents the computed performance of a score
type PerformanceCalculationResult struct {
	Total float64 `json:"total"`
	Speed float64 `json:"speed"`
	Aim   float64 `json:"aim"`
	Acc   float64 `json:"accuracy"`
}
