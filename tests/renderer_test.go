package tests

import (
	"image/png"
	"os"
	"testing"

	_ "embed"

	"github.com/Lekuruu/osutp-web/pkg/banners"
)

func TestBannerRenderRegular(t *testing.T) {
	player := banners.NewPlayer("Cookiezi", "Korea", "KR", "../web/static/images/flags/kr.gif", 1, 1, 3790)
	banner := banners.NewCleanStyleBanner(player)
	if banner == nil {
		t.Fatal("Failed to create CleanStyleBanner")
	}

	img := banner.Render()
	if img == nil {
		t.Error("Render() returned nil image")
	}

	outFile, err := os.Create("test_banner_cs.png")
	if err != nil {
		t.Fatalf("Failed to create output file: %v", err)
	}
	defer outFile.Close()

	err = png.Encode(outFile, img)
	if err != nil {
		t.Fatalf("Failed to encode image: %v", err)
	}
}

func TestBannerRenderCentered(t *testing.T) {
	player := banners.NewPlayer("Cookiezi", "Korea", "KR", "../web/static/images/flags/kr.gif", 1, 1, 3790)
	banner := banners.NewCleanStyleCenteredBanner(player)
	if banner == nil {
		t.Fatal("Failed to create CleanStyleCenteredBanner")
	}

	img := banner.Render()
	if img == nil {
		t.Error("Render() returned nil image")
	}

	outFile, err := os.Create("test_banner_csc.png")
	if err != nil {
		t.Fatalf("Failed to create output file: %v", err)
	}
	defer outFile.Close()

	err = png.Encode(outFile, img)
	if err != nil {
		t.Fatalf("Failed to encode image: %v", err)
	}
}

func TestBannerRenderOneLine(t *testing.T) {
	player := banners.NewPlayer("Cookiezi", "Korea", "KR", "../web/static/images/flags/kr.gif", 1, 1, 3790)
	banner := banners.NewCleanStyleOneLineBanner(player)
	if banner == nil {
		t.Fatal("Failed to create CleanStyleOneLineBanner")
	}

	img := banner.Render()
	if img == nil {
		t.Error("Render() returned nil image")
	}

	outFile, err := os.Create("test_banner_cso.png")
	if err != nil {
		t.Fatalf("Failed to create output file: %v", err)
	}
	defer outFile.Close()

	err = png.Encode(outFile, img)
	if err != nil {
		t.Fatalf("Failed to encode image: %v", err)
	}
}
