package tests

import (
	"testing"

	"github.com/Lekuruu/osutp-web/pkg/tp"
	osr "github.com/robloxxa/go-osr"
)

func TestPerformanceFreedomDive(t *testing.T) {
	performScorePerformanceTest(
		t,
		"performance_test_freedomdive.osr",
		"difficulty_test_freedomdive.osu",
	)
}

func TestPerformanceFreedomDiveAnother(t *testing.T) {
	performScorePerformanceTest(
		t,
		"performance_test_freedomdive_another.osr",
		"difficulty_test_freedomdive_another.osu",
	)
}

func TestPerformanceRemoteControl(t *testing.T) {
	performScorePerformanceTest(
		t,
		"performance_test_remotecontrol.osr",
		"difficulty_test_remotecontrol.osu",
	)
}

func TestPerformanceRogUnlimitation(t *testing.T) {
	performScorePerformanceTest(
		t,
		"performance_test_rogunlimitation.osr",
		"difficulty_test_rogunlimitation.osu",
	)
}

func TestPerformanceGimmeGimme(t *testing.T) {
	performScorePerformanceTest(
		t,
		"performance_test_gimmegimme.osr",
		"difficulty_test_gimmegimme.osu",
	)
}

func TestPerformanceMatzcore(t *testing.T) {
	performScorePerformanceTest(
		t,
		"performance_test_matzcore.osr",
		"difficulty_test_matzcore.osu",
	)
}

func TestPerformanceRedGoose(t *testing.T) {
	performScorePerformanceTest(
		t,
		"performance_test_redgoose.osr",
		"difficulty_test_redgoose.osu",
	)
}

func TestPerformanceAirman(t *testing.T) {
	performScorePerformanceTest(
		t,
		"performance_test_airman.osr",
		"difficulty_test_airman.osu",
	)
}

func TestPerformanceKillerSong(t *testing.T) {
	performScorePerformanceTest(
		t,
		"performance_test_killersong.osr",
		"difficulty_test_killersong.osu",
	)
}

func TestPerformanceBakaNA(t *testing.T) {
	performScorePerformanceTest(
		t,
		"performance_test_bakana.osr",
		"difficulty_test_bakana.osu",
	)
}

func TestPerformanceBoozehound(t *testing.T) {
	performScorePerformanceTest(
		t,
		"performance_test_boozehound.osr",
		"difficulty_test_boozehound.osu",
	)
}

func TestPerformnaceFourSeasonsOfLoneliness(t *testing.T) {
	performScorePerformanceTest(
		t,
		"performance_test_fourseasonsofloneliness.osr",
		"difficulty_test_fourseasonsofloneliness.osu",
	)
}

func TestPerformanceNyanNyan(t *testing.T) {
	performScorePerformanceTest(
		t,
		"performance_test_nyanyan.osr",
		"difficulty_test_nyanyan.osu",
	)
}

func TestPerformanceWahrheit(t *testing.T) {
	performScorePerformanceTest(
		t,
		"performance_test_wahrheit.osr",
		"difficulty_test_wahrheit.osu",
	)
}

func TestPerformancePoinsettia(t *testing.T) {
	performScorePerformanceTest(
		t,
		"performance_test_poinsettia.osr",
		"difficulty_test_poinsettia.osu",
	)
}

func performScorePerformanceTest(t *testing.T, replayFile string, beatmapFile string) {
	t.Helper()

	replay, err := osr.NewReplayFromFile(replayFile)
	if err != nil {
		t.Fatalf("failed to load replay (%s): %v", replayFile, err)
	}

	score := tp.NewScoreFromReplay(replay, beatmapFile)
	if score == nil {
		t.Fatalf("failed to create score from replay (%s)", replayFile)
	}

	beatmapDifficulty := performBeatmapDifficultyCalculation(t, beatmapFile, score.Mods)
	if beatmapDifficulty == nil {
		t.Fatalf("failed to calculate beatmap difficulty (%s)", beatmapFile)
	}

	request := tp.NewPerformanceCalculationRequest(score, beatmapDifficulty)
	result, err := request.Perform(*serviceUrl)
	if err != nil {
		t.Fatalf("failed to calculate performance (%s): %v", replayFile, err)
	}

	t.Logf(
		"Performance for %s: %.2f (Aim: %.2f, Speed: %.2f, Acc: %.2f)",
		replayFile,
		result.Total,
		result.Aim,
		result.Speed,
		result.Acc,
	)
}
