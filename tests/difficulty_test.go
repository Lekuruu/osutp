package tests

import (
	"testing"

	"github.com/Lekuruu/osutp/pkg/tp"
	osu "github.com/natsukagami/go-osu-parser"
)

func TestBeatmapRogUnlimitation(t *testing.T) {
	performBeatmapDifficultyCalculation(t, "difficulty_test_rogunlimitation.osu", 0)
}

func TestBeatmapFreedomDive(t *testing.T) {
	performBeatmapDifficultyCalculation(t, "difficulty_test_freedomdive.osu", 0)
}

func TestBeatmapFreedomDiveAnother(t *testing.T) {
	performBeatmapDifficultyCalculation(t, "difficulty_test_freedomdive_another.osu", 0)
}

func TestBeatmapRemoteControl(t *testing.T) {
	performBeatmapDifficultyCalculation(t, "difficulty_test_remotecontrol.osu", 0)
}

func TestBeatmapGimmeGimme(t *testing.T) {
	performBeatmapDifficultyCalculation(t, "difficulty_test_gimmegimme.osu", 0)
}

func TestBeatmapMatzcore(t *testing.T) {
	performBeatmapDifficultyCalculation(t, "difficulty_test_matzcore.osu", 0)
}

func TestBeatmapRedGoose(t *testing.T) {
	performBeatmapDifficultyCalculation(t, "difficulty_test_redgoose.osu", 0)
}

func TestBeatmapAirman(t *testing.T) {
	performBeatmapDifficultyCalculation(t, "difficulty_test_airman.osu", 0)
}

func TestBeatmapKillerSong(t *testing.T) {
	performBeatmapDifficultyCalculation(t, "difficulty_test_killersong.osu", 0)
}

func TestBeatmapBakaNA(t *testing.T) {
	performBeatmapDifficultyCalculation(t, "difficulty_test_bakana.osu", 0)
}

func TestBeatmapBoozehound(t *testing.T) {
	performBeatmapDifficultyCalculation(t, "difficulty_test_boozehound.osu", 0)
}

func TestBeatmapFourSeasonsOfLoneliness(t *testing.T) {
	performBeatmapDifficultyCalculation(t, "difficulty_test_fourseasonsofloneliness.osu", 0)
}

func TestBeatmapNyanNyan(t *testing.T) {
	performBeatmapDifficultyCalculation(t, "difficulty_test_nyanyan.osu", 0)
}

func TestBeatmapWahrheit(t *testing.T) {
	performBeatmapDifficultyCalculation(t, "difficulty_test_wahrheit.osu", 0)
}

func performBeatmapDifficultyCalculation(t *testing.T, beatmapFile string, mods uint32) *tp.DifficultyCalculationResult {
	t.Helper()

	beatmap, err := osu.ParseFile(beatmapFile)
	if err != nil {
		t.Fatalf("Failed to parse beatmap: %v", err)
	}

	response := tp.CalculateDifficulty(&beatmap, mods)
	if response == nil {
		t.Fatal("Failed to calculate difficulty")
	}

	t.Logf(
		"Difficulty Calculation Result: %.2f* / %.2f (Aim: %.2f* / %.2f, Speed: %.2f* / %.2f)\n",
		response.StarRating,
		response.Level(),
		response.AimStars,
		response.AimLevel(),
		response.SpeedStars,
		response.SpeedLevel(),
	)
	return response
}
