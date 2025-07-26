package tests

import (
	"image/png"
	"os"
	"testing"

	_ "embed"

	"github.com/Lekuruu/osutp-web/pkg/banners"
)

func TestBannerRender(t *testing.T) {
	player := banners.NewPlayer("Cookiezi", "Korea", "KR", "../web/static/images/flags/kr.gif", 1, 1, 3790)
	banner := banners.NewCleanStyleCenteredBanner(player)
	if banner == nil {
		t.Fatal("Failed to create CleanStyleBanner")
	}

	img := banner.Render()
	if img == nil {
		t.Error("Render() returned nil image")
	}

	outFile, err := os.Create("test_banner.png")
	if err != nil {
		t.Fatalf("Failed to create output file: %v", err)
	}
	defer outFile.Close()

	err = png.Encode(outFile, img)
	if err != nil {
		t.Fatalf("Failed to encode image: %v", err)
	}
}
