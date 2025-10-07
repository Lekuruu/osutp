package tp

import (
	"math"
	"sort"

	osu "github.com/natsukagami/go-osu-parser"
)

// DifficultyType represents the type of difficulty being calculated
type DifficultyType int

const (
	DifficultyTypeSpeed DifficultyType = 0
	DifficultyTypeAim   DifficultyType = 1
)

const (
	starScalingFactor    = 0.045
	extremeScalingFactor = 0.5
	playfieldWidth       = 512.0

	// The weighting of each strain value decays to 0.9 * it's previous value
	decayWeight = 0.9

	// In milliseconds. For difficulty calculation we will only look at the highest strain value in each time interval of size STRAIN_STEP.
	// This is to eliminate higher influence of stream over aim by simply having more HitObjects with high strain.
	// The higher this value, the less strains there will be, indirectly giving long beatmaps an advantage.
	strainStep = 400.0
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

func (result *DifficultyCalculationResult) Level() float64 {
	aimLevel := result.AimLevel()
	speedLevel := result.SpeedLevel()
	return ((speedLevel + aimLevel) + math.Abs(speedLevel-aimLevel)) / 2.125
}

func (result *DifficultyCalculationResult) AimLevel() float64 {
	return approximateTpLevel(result.AimDifficulty)
}

func (result *DifficultyCalculationResult) SpeedLevel() float64 {
	return approximateTpLevel(result.SpeedDifficulty)
}

// CalculateDifficulty calculates the difficulty of a beatmap with given mods
func CalculateDifficulty(beatmap *osu.Beatmap, mods uint32) *DifficultyCalculationResult {
	// Adjust beatmap attributes, based on relevant mods
	timeRate := 1.0
	adjustDifficulty(beatmap, mods, &timeRate)

	// Fill our custom tpHitObject class, that carries additional information
	tpHitObjects := make([]*TpHitObject, 0, len(beatmap.HitObjects))
	circleRadius := playfieldWidth / 16.0 * (1.0 - 0.7*(beatmap.CircleSize-5.0)/5.0)

	for _, hitObject := range beatmap.HitObjects {
		tpHitObjects = append(tpHitObjects, NewTpHitObject(hitObject, circleRadius, beatmap))
	}

	// Sort tpHitObjects by StartTime of the HitObjects - just to make sure
	sort.Slice(tpHitObjects, func(i, j int) bool {
		return tpHitObjects[i].HitObject.StartTime < tpHitObjects[j].HitObject.StartTime
	})

	// Calculate strain values
	if !calculateStrainValues(tpHitObjects, timeRate) {
		return nil
	}

	// Calculate difficulties
	speedDifficulty := calculateDifficultyForType(tpHitObjects, DifficultyTypeSpeed, timeRate)
	aimDifficulty := calculateDifficultyForType(tpHitObjects, DifficultyTypeAim, timeRate)

	// OverallDifficulty is not considered in this algorithm and neither is HpDrainRate. That means, that in this form the algorithm determines how hard it physically is
	// to play the map, assuming, that too much of an error will not lead to a death.
	// It might be desirable to include OverallDifficulty into map difficulty, but in my personal opinion it belongs more to the weighting of the actual peformance
	// and is superfluous in the beatmap difficulty rating.
	// If it were to be considered, then I would look at the hit window of normal HitCircles only, since Sliders and Spinners are (almost) "free" 300s and take map length
	// into account as well.

	// The difficulty can be scaled by any desired metric.
	// In osu!tp it gets squared to account for the rapid increase in difficulty as the limit of a human is approached. (Of course it also gets scaled afterwards.)
	// It would not be suitable for a star rating, therefore:

	// The following is a proposal to forge a star rating from 0 to 5. It consists of taking the square root of the difficulty, since by simply scaling the easier
	// 5-star maps would end up with one star.
	speedStars := math.Sqrt(speedDifficulty) * starScalingFactor
	aimStars := math.Sqrt(aimDifficulty) * starScalingFactor

	// Again, from own observations and from the general opinion of the community a map with high speed and low aim (or vice versa) difficulty is harder,
	// than a map with mediocre difficulty in both. Therefore we can not just add both difficulties together, but will introduce a scaling that favors extremes.
	starRating := speedStars + aimStars + math.Abs(speedStars-aimStars)*extremeScalingFactor

	// Another approach to this would be taking Speed and Aim separately to a chosen power, which again would be equivalent. This would be more convenient if
	// the hit window size is to be considered as well.

	// Note: The star rating is tuned extremely tight! Airman (/b/104229) and Freedom Dive (/b/126645), two of the hardest ranked maps, both score ~4.66 stars.
	// Expect the easier kind of maps that officially get 5 stars to obtain around 2 by this metric. The tutorial still scores about half a star.
	// Tune by yourself as you please. ;)

	return &DifficultyCalculationResult{
		AmountNormal:      beatmap.NbCircles,
		AmountSliders:     beatmap.NbSliders,
		AmountSpinners:    beatmap.NbSpinners,
		MaxCombo:          beatmap.MaxCombo,
		SpeedDifficulty:   speedDifficulty,
		AimDifficulty:     aimDifficulty,
		SpeedStars:        speedStars,
		AimStars:          aimStars,
		StarRating:        starRating,
		ApproachRate:      float32(beatmap.ApproachRate),
		CircleSize:        float32(beatmap.CircleSize),
		HpDrainRate:       float32(beatmap.HPDrainRate),
		OverallDifficulty: float32(beatmap.OverallDifficulty),
		SliderTickRate:    float64(beatmap.SliderTickRate),
		SliderMultiplier:  beatmap.SliderMultiplier,
	}
}

func mapDifficultyRange(difficulty, min, mid, max float64) float64 {
	if difficulty > 5 {
		return mid + (max-mid)*(difficulty-5)/5
	}
	if difficulty < 5 {
		return mid - (mid-min)*(5-difficulty)/5
	}
	return mid
}

func adjustDifficulty(beatmap *osu.Beatmap, mods uint32, timeRate *float64) {
	if mods&HardRock != 0 {
		beatmap.OverallDifficulty = math.Min(beatmap.OverallDifficulty*1.4, 10)
		beatmap.CircleSize = math.Min(beatmap.CircleSize*1.3, 10)
		beatmap.HPDrainRate = math.Min(beatmap.HPDrainRate*1.4, 10)
		beatmap.ApproachRate = math.Min(beatmap.ApproachRate*1.4, 10)
	}

	if mods&Easy != 0 {
		beatmap.OverallDifficulty = math.Max(beatmap.OverallDifficulty/2, 0)
		beatmap.CircleSize = math.Max(beatmap.CircleSize/2, 0)
		beatmap.HPDrainRate = math.Max(beatmap.HPDrainRate/2, 0)
		beatmap.ApproachRate = math.Max(beatmap.ApproachRate/2, 0)
	}

	if mods&DoubleTime != 0 || mods&Nightcore != 0 {
		*timeRate = 1.5
		recalculateBeatmapDifficulty(beatmap, *timeRate)
	}

	if mods&HalfTime != 0 {
		*timeRate = 0.75
		recalculateBeatmapDifficulty(beatmap, *timeRate)
	}
}

func recalculateBeatmapDifficulty(beatmap *osu.Beatmap, timeRate float64) {
	preEmpt := mapDifficultyRange(beatmap.ApproachRate, 1800, 1200, 450) / timeRate
	hitWindow300 := mapDifficultyRange(beatmap.OverallDifficulty, 80, 50, 20) / timeRate
	beatmap.OverallDifficulty = -(hitWindow300 - 80.0) / 6.0
	if preEmpt > 1200 {
		beatmap.ApproachRate = (1800 - preEmpt) / 120
	} else {
		beatmap.ApproachRate = (1200-preEmpt)/150 + 5
	}
}

func calculateStrainValues(tpHitObjects []*TpHitObject, timeRate float64) bool {
	if len(tpHitObjects) == 0 {
		return false
	}

	// First hitObject starts at strain 1
	// 1 is the default for strain values, so we don't need to set it here
	for i := 1; i < len(tpHitObjects); i++ {
		tpHitObjects[i].CalculateStrains(tpHitObjects[i-1], timeRate)
	}

	return true
}

func calculateDifficultyForType(tpHitObjects []*TpHitObject, diffType DifficultyType, timeRate float64) float64 {
	actualStrainStep := strainStep * timeRate

	// Find the highest strain value within each strain step
	highestStrains := make([]float64, 0)
	intervalEndTime := actualStrainStep
	maximumStrain := 0.0

	var previousHitObject *TpHitObject

	for _, hitObject := range tpHitObjects {
		// While we are beyond the current interval push the
		// currently available maximum to our strain list
		for float64(hitObject.HitObject.StartTime) > intervalEndTime {
			highestStrains = append(highestStrains, maximumStrain)

			// The maximum strain of the next interval is not zero by default!
			// We need to take the last hitObject we encountered, take its strain
			// and apply the decay until the beginning of the next interval.
			if previousHitObject == nil {
				maximumStrain = 0
			} else {
				decay := math.Pow(decayBase[int(diffType)],
					(intervalEndTime-float64(previousHitObject.HitObject.StartTime))/1000)
				maximumStrain = previousHitObject.Strains[int(diffType)] * decay
			}

			// Go to the next time interval
			intervalEndTime += actualStrainStep
		}

		// Obtain maximum strain
		if hitObject.Strains[int(diffType)] > maximumStrain {
			maximumStrain = hitObject.Strains[int(diffType)]
		}

		previousHitObject = hitObject
	}

	// Build the weighted sum over the highest strains for each interval
	difficulty := 0.0
	weight := 1.0

	// Sort from highest to lowest strain
	sort.Float64s(highestStrains)

	// Reverse the slice
	for i, j := 0, len(highestStrains)-1; i < j; i, j = i+1, j-1 {
		highestStrains[i], highestStrains[j] = highestStrains[j], highestStrains[i]
	}

	for _, strain := range highestStrains {
		difficulty += weight * strain
		weight *= decayWeight
	}

	return difficulty
}

func approximateTpLevel(value float64) float64 {
	return 0.0877*value - 68.3
}
