package tp

import (
	"math"

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
	Mods            uint32 `json:"mods"`
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
		Mods:            uint32(replay.Mods),
	}
}

func (s *Score) IsRelaxing() bool {
	return s.Mods&Relax != 0 || s.Mods&Autopilot != 0
}

func (s *Score) IsAutoplay() bool {
	return s.Mods&Autoplay != 0
}

func (s *Score) TotalHits() int {
	return s.Amount300 + s.Amount100 + s.Amount50 + s.AmountMiss
}

func (s *Score) TotalSuccessfulHits() int {
	return s.Amount300 + s.Amount100 + s.Amount50
}

func (s *Score) Accuracy() float64 {
	totalHits := s.TotalHits()
	if totalHits <= 0 {
		return 0.0
	}
	accuracy := float64(300*s.Amount300+100*s.Amount100+50*s.Amount50) / float64(totalHits*300)
	return math.Max(math.Min(accuracy, 1.0), 0.0)
}

// PerformanceCalculationResult represents the computed tp performance of a score
type PerformanceCalculationResult struct {
	Total float64 `json:"total"`
	Speed float64 `json:"speed"`
	Aim   float64 `json:"aim"`
	Acc   float64 `json:"acc"`
}

// CalculatePerformance calculates the tp performance of a score
func CalculatePerformance(difficulty *DifficultyCalculationResult, score *Score) *PerformanceCalculationResult {
	if score.IsRelaxing() || score.IsAutoplay() {
		return &PerformanceCalculationResult{}
	}

	// This is being adjusted to keep the final pp value scaled
	// around what it used to be when changing things
	multiplier := 1.1

	// Custom multipliers for NoFail and SpunOut
	if score.Mods&NoFail != 0 {
		multiplier *= 0.9
	}
	if score.Mods&SpunOut != 0 {
		multiplier *= 0.95
	}

	aim := computeAimValue(difficulty, score)
	speed := computeSpeedValue(difficulty, score)
	acc := computeAccValue(difficulty, score)

	attributes := math.Pow(aim, 1.1) + math.Pow(speed, 1.1) + math.Pow(acc, 1.1)
	total := math.Pow(attributes, 1.0/1.1) * multiplier

	return &PerformanceCalculationResult{
		Total: total,
		Aim:   aim,
		Speed: speed,
		Acc:   acc,
	}
}

func computeAimValue(difficulty *DifficultyCalculationResult, score *Score) float64 {
	aimValue := math.Pow(5.0*math.Max(1.0, float64(difficulty.AimStars)/0.0358)-4.0, 3.0) / 100000.0

	// Longer maps are worth more
	aimValue *= 1 + 0.1*math.Min(1.0, float64(score.TotalHits())/1500.0)

	// Penalize misses exponentially
	// This mainly fixes TAG4 maps and the likes until a per-hitobject solution is available
	aimValue *= math.Pow(0.97, float64(score.AmountMiss))

	// Combo scaling
	if difficulty.MaxCombo > 0 {
		aimValue *= math.Min(1.0, math.Pow(float64(score.MaxCombo), 0.8)/math.Pow(float64(difficulty.MaxCombo), 0.8))
	}

	approachRateFactor := 1.0

	if difficulty.ApproachRate > 10.0 {
		approachRateFactor += 0.3 * (float64(difficulty.ApproachRate) - 10.0)
	} else if difficulty.ApproachRate < 8.0 {
		// Hidden is worth more with lower AR
		if score.Mods&Hidden != 0 {
			approachRateFactor += 0.02 * (8.0 - float64(difficulty.ApproachRate))
		} else {
			approachRateFactor += 0.01 * (8.0 - float64(difficulty.ApproachRate))
		}
	}

	aimValue *= approachRateFactor

	// Hidden Bonus
	if score.Mods&Hidden != 0 {
		aimValue *= 1.18
	}

	// Flashlight Bonus
	if score.Mods&Flashlight != 0 {
		aimValue *= 1.36
	}

	// Scale aim value with accuracy, slightly
	aimValue *= 0.5 + score.Accuracy()/2.0

	// It is important to also consider overall difficulty when doing that
	aimValue *= 0.98 + math.Pow(float64(difficulty.OverallDifficulty), 2)/2500

	return aimValue
}

func computeSpeedValue(difficulty *DifficultyCalculationResult, score *Score) float64 {
	speedValue := math.Pow(5.0*math.Max(1.0, float64(difficulty.SpeedStars)/0.0358)-4.0, 3.0) / 100000.0

	// Longer maps are worth more
	speedValue *= 1 + 0.1*math.Min(1.0, float64(score.TotalHits())/1500.0)

	// Penalize misses exponentially
	// This mainly fixes TAG4 maps and the likes until a per-hitobject solution is available
	speedValue *= math.Pow(0.97, float64(score.AmountMiss))

	// Combo scaling
	if difficulty.MaxCombo > 0 {
		speedValue *= math.Min(1.0, math.Pow(float64(score.MaxCombo), 0.8)/math.Pow(float64(difficulty.MaxCombo), 0.8))
	}

	// Scale speed value with accuracy, slightly
	speedValue *= 0.5 + score.Accuracy()/2.0

	// It is important to also consider overall difficulty when doing that
	speedValue *= 0.98 + math.Pow(float64(difficulty.OverallDifficulty), 2)/2500

	return speedValue
}

func computeAccValue(difficulty *DifficultyCalculationResult, score *Score) float64 {
	// This percentage only considers HitCircles of any value
	// In this part of the calculation we focus on hitting the timing hit window
	betterAccuracyPercentage := 0.0

	if difficulty.AmountNormal > 0 {
		betterAccuracyPercentage = float64((score.Amount300-(score.TotalHits()-difficulty.AmountNormal))*6+score.Amount100*2+score.Amount50) /
			float64(difficulty.AmountNormal*6)
	}

	// It is possible to reach a negative accuracy with this formula. Cap it at zero - zero points
	if betterAccuracyPercentage < 0 {
		betterAccuracyPercentage = 0
	}

	// Lots of arbitrary values from testing.
	// Considering to use derivation from perfect accuracy in a probabilistic manner - assume normal distribution
	accValue := math.Pow(
		math.Pow(1.3, float64(difficulty.OverallDifficulty))*math.Pow(betterAccuracyPercentage, 15)/2,
		1.6,
	) * 8.3

	// Bonus for many hitcircles - it's harder to keep good accuracy up for longer
	accValue *= math.Min(1.15, math.Pow(float64(difficulty.AmountNormal)/1000.0, 0.3))

	// Hidden Bonus
	if score.Mods&Hidden != 0 {
		accValue *= 1.02
	}

	// Flashlight Bonus
	if score.Mods&Flashlight != 0 {
		accValue *= 1.02
	}

	return accValue
}
