package tests

import (
	"testing"

	"github.com/Lekuruu/osutp-web/pkg/tp"
	osr "github.com/robloxxa/go-osr"
)

func TestPerformanceFreedomDive(t *testing.T) {
	replay, err := osr.NewReplayFromFile("performance_test_freedomdive.osr")
	if err != nil {
		t.Fatalf("failed to load replay: %v", err)
	}

	score := tp.NewScoreFromReplay(replay, "difficulty_test_freedomdive.osu")
	if score == nil {
		t.Fatal("failed to create score from replay")
	}

	beatmapDifficulty := performBeatmapDifficultyCalculation(t, "difficulty_test_freedomdive.osu", score.Mods)
	if beatmapDifficulty == nil {
		t.Fatal("failed to calculate beatmap difficulty")
	}

	// Perform performance calculation
	request := tp.NewPerformanceCalculationRequest(score, beatmapDifficulty)
	result, err := request.Perform(*serviceUrl)
	if err != nil {
		t.Fatalf("failed to calculate performance: %v", err)
	}

	t.Logf("Performance: %.2f (Aim: %.2f, Speed: %.2f, Acc: %.2f)", result.Total, result.Aim, result.Speed, result.Acc)
}

func TestPerformanceRemoteControl(t *testing.T) {
	replay, err := osr.NewReplayFromFile("performance_test_remotecontrol.osr")
	if err != nil {
		t.Fatalf("failed to load replay: %v", err)
	}

	score := tp.NewScoreFromReplay(replay, "difficulty_test_remotecontrol.osu")
	if score == nil {
		t.Fatal("failed to create score from replay")
	}

	beatmapDifficulty := performBeatmapDifficultyCalculation(t, "difficulty_test_remotecontrol.osu", score.Mods)
	if beatmapDifficulty == nil {
		t.Fatal("failed to calculate beatmap difficulty")
	}

	// Perform performance calculation
	request := tp.NewPerformanceCalculationRequest(score, beatmapDifficulty)
	result, err := request.Perform(*serviceUrl)
	if err != nil {
		t.Fatalf("failed to calculate performance: %v", err)
	}

	t.Logf("Performance: %.2f (Aim: %.2f, Speed: %.2f, Acc: %.2f)", result.Total, result.Aim, result.Speed, result.Acc)
}
