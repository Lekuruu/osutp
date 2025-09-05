package common

import (
	"math"

	"github.com/Lekuruu/osutp-web/internal/database"
	"github.com/Lekuruu/osutp-web/pkg/tp"
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

	attributes := database.DifficultyAttributes{}
	for _, mod := range modCombinations {
		request := tp.NewDifficultyCalculationRequestFromBeatmap(beatmap, int(mod))
		if request == nil {
			continue
		}

		response, err := request.Perform(state.Config.TpServiceUrl)
		if err != nil {
			return err
		}
		attributes[mod] = map[string]float64{}
		attributes[mod]["ApproachRate"] = round(float64(response.ApproachRate), 6)
		attributes[mod]["OverallDifficulty"] = round(float64(response.OverallDifficulty), 6)
		attributes[mod]["HpDrainRate"] = round(float64(response.HpDrainRate), 6)
		attributes[mod]["CircleSize"] = round(float64(response.CircleSize), 6)
		attributes[mod]["SpeedDifficulty"] = round(response.SpeedDifficulty, 6)
		attributes[mod]["AimDifficulty"] = round(response.AimDifficulty, 6)
		attributes[mod]["SpeedStars"] = round(response.SpeedStars, 6)
		attributes[mod]["AimStars"] = round(response.AimStars, 6)
		attributes[mod]["StarRating"] = round(response.StarRating, 6)
	}

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
