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

const (
	ColorDarkBlue = 0x1e1e64
	ColorBlack    = 0x000028
)

var (
	ColorDarkBlueRGBA = color.RGBA{0x1e, 0x1e, 0x64, 0xff}
	ColorBlackRGBA    = color.RGBA{0x00, 0x00, 0x28, 0xff}
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

	boldFont := banner.GetFont("bold")
	textWidth := font.MeasureString(boldFont, text).Ceil()

	width := marginX*2 + textWidth
	height := marginY*2 + lineHeight*2
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	renderText(
		text,
		boldFont,
		ColorDarkBlueRGBA, img,
		image.Point{X: marginX, Y: marginY + int(boldFont.Metrics().Ascent.Round())},
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

	regularFont := banner.GetFont("regular")
	textWidth := font.MeasureString(regularFont, text).Ceil()

	width := marginX*2 + textWidth
	height := marginY*2 + lineHeight*2
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	renderText(
		text,
		regularFont,
		ColorDarkBlueRGBA, img,
		image.Point{X: marginX, Y: marginY + int(regularFont.Metrics().Ascent.Round())},
	)
	return img
}
