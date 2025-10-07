package common

import (
	"math"
	"sync"

	"github.com/Lekuruu/osutp/internal/database"
	"github.com/Lekuruu/osutp/pkg/tp"
	osu "github.com/natsukagami/go-osu-parser"
)

var modCombinations []uint32 = []uint32{
	tp.NoMod, tp.HardRock, tp.DoubleTime, tp.HalfTime, tp.Easy,
	tp.HardRock | tp.DoubleTime, tp.HardRock | tp.HalfTime, tp.Easy | tp.DoubleTime, tp.Easy | tp.HalfTime,
}

func UpdateBeatmapDifficulty(file []byte, schema *database.Beatmap, state *State) error {
	beatmap, err := osu.ParseBytes(file)
	if err != nil {
		return err
	}

	// Calculate difficulty attributes for all mod
	// combinations in parallel to speed up the process
	attributes := database.DifficultyAttributes{}
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, mod := range modCombinations {
		wg.Add(1)

		go func(mod uint32) {
			defer wg.Done()

			result := tp.CalculateDifficulty(&beatmap, uint32(mod))
			if result == nil {
				return
			}

			modAttrs := map[string]float64{
				"ApproachRate":      round(float64(result.ApproachRate), 6),
				"OverallDifficulty": round(float64(result.OverallDifficulty), 6),
				"HpDrainRate":       round(float64(result.HpDrainRate), 6),
				"CircleSize":        round(float64(result.CircleSize), 6),
				"SpeedDifficulty":   round(result.SpeedDifficulty, 6),
				"AimDifficulty":     round(result.AimDifficulty, 6),
				"SpeedStars":        round(result.SpeedStars, 6),
				"AimStars":          round(result.AimStars, 6),
				"StarRating":        round(result.StarRating, 6),
			}

			mu.Lock()
			attributes[mod] = modAttrs
			mu.Unlock()
		}(mod)
	}

	// Wait for difficulty calculations to finish
	wg.Wait()

	schema.DifficultyAttributes = attributes
	schema.AmountNormal = beatmap.NbCircles
	schema.AmountSliders = beatmap.NbSliders
	schema.AmountSpinners = beatmap.NbSpinners
	schema.MaxCombo = beatmap.MaxCombo
	state.Database.Save(schema)
	return nil
}

func round(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}
