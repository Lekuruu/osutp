package tests

import (
	"strings"
	"testing"

	"github.com/Lekuruu/osutp-web/pkg/tp"
	osu "github.com/natsukagami/go-osu-parser"
)

func TestBeatmapRogUnlimitation(t *testing.T) {
	performBeatmapDifficultyCalculation(t, "difficulty_test_rogunlimitation.osu", 0)
}

func TestBeatmapFreedomDive(t *testing.T) {
	performBeatmapDifficultyCalculation(t, "difficulty_test_freedomdive.osu", 0)
}

func performBeatmapDifficultyCalculation(t *testing.T, beatmapFile string, mods int) *tp.DifficultyCalculationResult {
	beatmap, err := osu.ParseFile(beatmapFile)
	if err != nil {
		t.Fatalf("Failed to parse beatmap: %v", err)
	}

	request := tp.NewDifficultyCalculationRequestFromBeatmap(beatmap, mods)
	if request == nil {
		t.Fatal("Failed to create difficulty calculation request from beatmap")
	}

	response, err := request.Perform(*serviceUrl)
	if err != nil {
		// Check if the service was actually running
		if strings.Contains(err.Error(), "connection refused") {
			t.Skipf("Skipping test because the service is not running at %s", *serviceUrl)
			return nil
		}
		t.Fatalf("Failed to perform difficulty calculation request: %v", err)
	}

	t.Logf("Difficulty Calculation Result: %.2f\n", response.StarRating)
	return response
}
