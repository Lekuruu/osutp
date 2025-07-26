package banners

import (
	"os"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

func LoadDefaultFont(options *truetype.Options) (font.Face, error) {
	return LoadTrueTypeFont("../web/static/fonts/tahoma.ttf", options)
}

func LoadDefaultFontBold(options *truetype.Options) (font.Face, error) {
	return LoadTrueTypeFont("../web/static/fonts/tahoma-bold.ttf", options)
}

func LoadTrueTypeFont(path string, options *truetype.Options) (font.Face, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	fnt, err := truetype.Parse(b)
	if err != nil {
		return nil, err
	}
	face := truetype.NewFace(fnt, options)
	return face, nil
}

func DefaultFontOptions(size float64) *truetype.Options {
	return &truetype.Options{
		Size:              size,
		DPI:               72,
		Hinting:           font.HintingFull,
		GlyphCacheEntries: 1024,
	}
}
