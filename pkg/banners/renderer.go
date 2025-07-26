package banners

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"os"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"golang.org/x/text/message"
)

var printer = message.NewPrinter(message.MatchLanguage("en"))

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
	textPrefix := "Global osu!tp rank for "
	textUsername := fmt.Sprintf("%s: ", banner.GetPlayer().Username())
	textRank := printer.Sprintf("#%d", banner.GetPlayer().GlobalRank())

	boldFont := banner.GetFont("bold")
	boldFontLarge := banner.GetFont("bold_large")
	textPrefixWidth := font.MeasureString(boldFont, textPrefix).Ceil()
	textUsernameWidth := font.MeasureString(boldFont, textUsername).Ceil()
	textRankWidth := font.MeasureString(boldFontLarge, textRank).Ceil()

	width := marginX*2 + textPrefixWidth + textUsernameWidth + textRankWidth
	height := marginY*2 + lineHeight*2
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	renderText(
		textPrefix,
		boldFont,
		ColorDarkBlueRGBA, img,
		image.Point{X: marginX, Y: marginY + int(boldFont.Metrics().Ascent.Round())},
	)
	renderText(
		textUsername,
		boldFont,
		ColorBlackRGBA, img,
		image.Point{X: marginX + textPrefixWidth, Y: marginY + int(boldFont.Metrics().Ascent.Round())},
	)
	renderText(
		textRank,
		boldFontLarge,
		ColorBlackRGBA, img,
		image.Point{X: marginX + textPrefixWidth + textUsernameWidth, Y: marginY + int(boldFontLarge.Metrics().Ascent.Round()) - 4},
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
