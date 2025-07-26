package banners

import (
	"image"
	"image/draw"

	"golang.org/x/image/font"
)

type Banner interface {
	Id() int
	Name() string
	Render() image.Image

	GetPlayer() Player
	SetPlayer(player Player)
	GetFont(key string) font.Face
	SetFont(key string, font font.Face)
}

type BaseBanner struct {
	player Player
	fonts  map[string]font.Face
}

func (b *BaseBanner) Id() int {
	return 0
}

func (b *BaseBanner) Name() string {
	return "Banner"
}

func (b *BaseBanner) GetPlayer() Player {
	return b.player
}

func (b *BaseBanner) SetPlayer(player Player) {
	b.player = player
}

func (b *BaseBanner) GetFont(key string) font.Face {
	return b.fonts[key]
}

func (b *BaseBanner) SetFont(key string, font font.Face) {
	b.fonts[key] = font
}

type CleanStyleBanner struct {
	BaseBanner
}

func (b *CleanStyleBanner) Id() int {
	return 0
}

func (b *CleanStyleBanner) Name() string {
	return "Clean-Style"
}

func (b *CleanStyleBanner) Render() image.Image {
	const yOffset = -20
	top := renderGlobal(b)
	bottom := renderCountry(b)

	// Use bounds of top and bottom to calculate new width & height
	width := top.Bounds().Dx()
	height := top.Bounds().Dy() + bottom.Bounds().Dy() + yOffset
	imagePointOffset := image.Point{Y: top.Bounds().Dy() + yOffset}

	// Render new image with top and bottom parts
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(img, top.Bounds(), top, image.Point{}, draw.Over)
	draw.Draw(img, bottom.Bounds().Add(imagePointOffset), bottom, image.Point{}, draw.Over)
	return img
}

func NewCleanStyleBanner(player Player) *CleanStyleBanner {
	defaultFont, err := LoadDefaultFont(DefaultFontOptions(20))
	if err != nil {
		return nil
	}
	defaultFontBold, err := LoadDefaultFontBold(DefaultFontOptions(20))
	if err != nil {
		return nil
	}
	banner := &CleanStyleBanner{
		BaseBanner: BaseBanner{
			player: player,
			fonts:  make(map[string]font.Face),
		},
	}
	banner.SetFont("regular", defaultFont)
	banner.SetFont("bold", defaultFontBold)
	return banner
}

type CleanStyleCenteredBanner struct {
	BaseBanner
}

func (b *CleanStyleCenteredBanner) Id() int {
	return 1
}

func (b *CleanStyleCenteredBanner) Name() string {
	return "Clean-Style (Centered)"
}

func (b *CleanStyleCenteredBanner) Render() image.Image {
	const yOffset = -20
	top := renderGlobal(b)
	bottom := renderCountry(b)

	// Calculate new width & height based on top and bottom images
	width := top.Bounds().Dx()
	if bottom.Bounds().Dx() > width {
		width = bottom.Bounds().Dx()
	}
	height := top.Bounds().Dy() + bottom.Bounds().Dy() + yOffset
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Center top
	topX := (width - top.Bounds().Dx()) / 2
	draw.Draw(img, top.Bounds().Add(image.Point{X: topX, Y: 0}), top, image.Point{}, draw.Over)

	// Center bottom
	bottomX := (width - bottom.Bounds().Dx()) / 2
	bottomY := top.Bounds().Dy() + yOffset
	draw.Draw(img, bottom.Bounds().Add(image.Point{X: bottomX, Y: bottomY}), bottom, image.Point{}, draw.Over)

	return img
}

func NewCleanStyleCenteredBanner(player Player) *CleanStyleCenteredBanner {
	defaultFont, err := LoadDefaultFont(DefaultFontOptions(20))
	if err != nil {
		return nil
	}
	defaultFontBold, err := LoadDefaultFontBold(DefaultFontOptions(20))
	if err != nil {
		return nil
	}
	banner := &CleanStyleCenteredBanner{
		BaseBanner: BaseBanner{
			player: player,
			fonts:  make(map[string]font.Face),
		},
	}
	banner.SetFont("regular", defaultFont)
	banner.SetFont("bold", defaultFontBold)
	return banner
}

type CleanStyleOneLineBanner struct {
	BaseBanner
}

func (b *CleanStyleOneLineBanner) Id() int {
	return 2
}

func (b *CleanStyleOneLineBanner) Name() string {
	return "Clean-Style (1-Line)"
}

func (b *CleanStyleOneLineBanner) Render() image.Image {
	const xOffset = -20
	top := renderGlobal(b)
	bottom := renderCountry(b)

	// Calculate new width & height based on top and bottom images
	width := top.Bounds().Dx() + bottom.Bounds().Dx() + xOffset
	height := top.Bounds().Dy()

	// Render side-by-side
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(img, top.Bounds(), top, image.Point{}, draw.Over)
	draw.Draw(img, bottom.Bounds().Add(image.Point{X: top.Bounds().Dx() + xOffset, Y: 0}), bottom, image.Point{}, draw.Over)
	return img
}

func NewCleanStyleOneLineBanner(player Player) *CleanStyleOneLineBanner {
	defaultFont, err := LoadDefaultFont(DefaultFontOptions(20))
	if err != nil {
		return nil
	}
	defaultFontBold, err := LoadDefaultFontBold(DefaultFontOptions(20))
	if err != nil {
		return nil
	}
	banner := &CleanStyleOneLineBanner{
		BaseBanner: BaseBanner{
			player: player,
			fonts:  make(map[string]font.Face),
		},
	}
	banner.SetFont("regular", defaultFont)
	banner.SetFont("bold", defaultFontBold)
	return banner
}
