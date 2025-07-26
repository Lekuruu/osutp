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

const flagImagePath = "../../web/static/images/flags/"

var printer = message.NewPrinter(message.MatchLanguage("en"))
var imageCache = make(map[string]image.Image)

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

func loadImageCached(location string) (image.Image, error) {
	if img, ok := imageCache[location]; ok {
		return img, nil
	}

	img, err := loadImage(location)
	if err != nil {
		return nil, err
	}

	imageCache[location] = img
	return img, nil
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

	textPrefix := "Rated "
	textCountryRank := banner.GetPlayer().CountryRankOrdinal()
	textOf := " of "
	textCountry := fmt.Sprintf("%s with ", banner.GetPlayer().Country())
	textTp := printer.Sprintf("%dtp", banner.GetPlayer().Tp())
	textDot := "."

	boldFont := banner.GetFont("bold")
	regularFont := banner.GetFont("regular")
	textPrefixWidth := font.MeasureString(regularFont, textPrefix).Ceil()
	textCountryRankWidth := font.MeasureString(boldFont, textCountryRank).Ceil()
	textOfWidth := font.MeasureString(regularFont, textOf).Ceil()
	textCountryWidth := font.MeasureString(regularFont, textCountry).Ceil()
	textTpWidth := font.MeasureString(boldFont, textTp).Ceil()
	textDotWidth := font.MeasureString(regularFont, textDot).Ceil()

	width := marginX*2 + textPrefixWidth + textCountryRankWidth + textOfWidth + textCountryWidth + textTpWidth + textDotWidth
	height := marginY*2 + lineHeight*2
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	renderText(
		textPrefix,
		regularFont,
		ColorDarkBlueRGBA, img,
		image.Point{X: marginX, Y: marginY + int(regularFont.Metrics().Ascent.Round())},
	)
	renderText(
		textCountryRank,
		boldFont,
		ColorBlackRGBA, img,
		image.Point{X: marginX + textPrefixWidth, Y: marginY + int(regularFont.Metrics().Ascent.Round())},
	)
	renderText(
		textOf,
		regularFont,
		ColorDarkBlueRGBA, img,
		image.Point{X: marginX + textPrefixWidth + textCountryRankWidth, Y: marginY + int(regularFont.Metrics().Ascent.Round())},
	)
	renderText(
		textCountry,
		regularFont,
		ColorDarkBlueRGBA, img,
		image.Point{X: marginX + textPrefixWidth + textCountryRankWidth + textOfWidth, Y: marginY + int(regularFont.Metrics().Ascent.Round())},
	)
	renderText(
		textTp,
		boldFont,
		ColorBlackRGBA, img,
		image.Point{X: marginX + textPrefixWidth + textCountryRankWidth + textOfWidth + textCountryWidth, Y: marginY + int(regularFont.Metrics().Ascent.Round())},
	)
	renderText(
		textDot,
		regularFont,
		ColorDarkBlueRGBA, img,
		image.Point{X: marginX + textPrefixWidth + textCountryRankWidth + textOfWidth + textCountryWidth + textTpWidth, Y: marginY + int(regularFont.Metrics().Ascent.Round())},
	)
	return img
}
