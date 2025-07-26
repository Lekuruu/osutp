package banners

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"os"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

func loadImage(location string) (image.Image, error) {
	// Load the flag image from the specified location
	img, err := os.Open(location)
	if err != nil {
		return nil, err
	}
	defer img.Close()

	// Decode the image
	decodedImage, _, err := image.Decode(img)
	if err != nil {
		return nil, err
	}
	return decodedImage, nil
}

func renderText(text string, face font.Face, color color.Color, img draw.Image, position image.Point) {
	// Draw the text onto the image
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(color),
		Face: face,
	}
	d.Dot = fixed.P(position.X, position.Y)
	d.DrawString(text)
}

func renderGlobal(banner Banner) image.Image {
	const (
		lineHeight = 12
		marginX    = 20
		marginY    = 10
	)
	text := fmt.Sprintf(
		"Global osu!tp rank for %s: #%d",
		banner.GetPlayer().Username(), banner.GetPlayer().GlobalRank(),
	)
	textWidth := font.MeasureString(banner.GetFont(), text).Ceil()

	width := marginX*2 + textWidth
	height := marginY*2 + lineHeight*2
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// TODO: Make this more accurate to the original
	renderText(
		text,
		banner.GetFont(),
		color.White, img,
		image.Point{X: marginX, Y: marginY + int(banner.GetFont().Metrics().Ascent.Round())},
	)

	return img
}

func renderCountry(banner Banner) image.Image {
	const (
		lineHeight       = 12
		marginX          = 20
		marginY          = 10
		flagHeight       = 24
		spacingAfterFlag = 8
	)
	text := fmt.Sprintf(
		"Rated %s of %s with %dtp.",
		banner.GetPlayer().CountryRankOrdinal(),
		banner.GetPlayer().Country(),
		banner.GetPlayer().Tp(),
	)
	textWidth := font.MeasureString(banner.GetFont(), text).Ceil()

	width := marginX*2 + textWidth
	height := marginY*2 + lineHeight*2
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// TODO: Make this more accurate to the original
	renderText(
		text,
		banner.GetFont(),
		color.White, img,
		image.Point{X: marginX, Y: marginY + int(banner.GetFont().Metrics().Ascent.Round())},
	)

	return img
}
