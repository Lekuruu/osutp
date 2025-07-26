package tests

import (
	"flag"
	"strings"
	"testing"

	"github.com/Lekuruu/osutp-web/pkg/tp"
	osu "github.com/natsukagami/go-osu-parser"
)

var serviceUrl = flag.String("service_url", "http://localhost:5028", "URL of the tp service")

func TestBeatmapRogUnlimitation(t *testing.T) {
	performBeatmapDifficultyCalculation(t, "tp_test_rogunlimitation.osu")
}

func TestBeatmapFreedomDive(t *testing.T) {
	performBeatmapDifficultyCalculation(t, "tp_test_freedomdive.osu")
}

func performBeatmapDifficultyCalculation(t *testing.T, beatmapFile string) {
	beatmap, err := osu.ParseFile(beatmapFile)
	if err != nil {
		t.Fatalf("Failed to parse beatmap: %v", err)
	}

	request := tp.NewDifficultyCalculationRequestFromBeatmap(beatmap)
	if request == nil {
		t.Fatal("Failed to create difficulty calculation request from beatmap")
	}

	response, err := request.Perform(*serviceUrl)
	if err != nil {
		// Check if the service was actually running
		if strings.Contains(err.Error(), "connection refused") {
			t.Skipf("Skipping test because the service is not running at %s", *serviceUrl)
			return
		}
		t.Fatalf("Failed to perform difficulty calculation request: %v", err)
	}

	t.Logf("Difficulty Calculation Result: %.2f\n", response.StarRating)
}
