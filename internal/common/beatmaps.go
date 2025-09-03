package common

import (
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
		attributes[mod]["ApproachRate"] = float64(response.ApproachRate)
		attributes[mod]["OverallDifficulty"] = float64(response.OverallDifficulty)
		attributes[mod]["HpDrainRate"] = float64(response.HpDrainRate)
		attributes[mod]["CircleSize"] = float64(response.CircleSize)
		attributes[mod]["SpeedDifficulty"] = response.SpeedDifficulty
		attributes[mod]["AimDifficulty"] = response.AimDifficulty
		attributes[mod]["SpeedStars"] = response.SpeedStars
		attributes[mod]["AimStars"] = response.AimStars
		attributes[mod]["StarRating"] = response.StarRating
	}

	schema.DifficultyAttributes = attributes
	schema.AmountNormal = beatmap.NbCircles
	schema.AmountSliders = beatmap.NbSliders
	schema.AmountSpinners = beatmap.NbSpinners
	schema.MaxCombo = beatmap.MaxCombo
	state.Database.Save(schema)
	return nil
}
