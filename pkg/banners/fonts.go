package banners

import (
	"os"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

var TahomaFontPath = "../web/static/fonts/tahoma.ttf"
var TahomaBoldFontPath = "../web/static/fonts/tahoma-bold.ttf"

func LoadDefaultFont(size float64) (font.Face, error) {
	return LoadTrueTypeFont(TahomaFontPath, defaultFontOptions(size))
}

func LoadDefaultFontBold(size float64) (font.Face, error) {
	return LoadTrueTypeFont(TahomaBoldFontPath, defaultFontOptions(size))
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

func defaultFontOptions(size float64) *truetype.Options {
	return &truetype.Options{
		Size:              size,
		DPI:               72,
		Hinting:           font.HintingFull,
		GlyphCacheEntries: 1024,
	}
}
