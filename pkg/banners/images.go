package banners

import (
	"image"
	"os"

	"golang.org/x/image/draw"
)

var imageCache = make(map[string]image.Image)

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

func resizeImage(img image.Image, scale float64) image.Image {
	if scale <= 0 {
		return img
	}

	newWidth := int(float64(img.Bounds().Dx()) * scale)
	newHeight := int(float64(img.Bounds().Dy()) * scale)

	resizedImg := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
	draw.CatmullRom.Scale(resizedImg, resizedImg.Bounds(), img, img.Bounds(), draw.Over, nil)

	return resizedImg
}
