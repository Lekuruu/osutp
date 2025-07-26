package banners

import (
	"image"
	"os"
	"strings"
)

var flagImagePath = projectRootDirectory + "/static/images/flags/png/"
var projectRootDirectory = findRootDirectory()
var defaultFlagImage = loadDefaultFlagImage()

func loadFlagImage(country string) image.Image {
	country = strings.ToLower(country)
	flagImage, err := loadImageCached(flagImagePath + country + ".png")
	if err != nil {
		return defaultFlagImage
	}
	return flagImage
}

func loadDefaultFlagImage() image.Image {
	img, err := loadImageCached(flagImagePath + "unknown.png")
	if err != nil {
		return image.NewRGBA(image.Rect(0, 0, 0, 0))
	}
	return img
}

func findRootDirectory() string {
	workingDirectory, err := os.Getwd()
	if err != nil {
		panic("Failed to get working directory: " + err.Error())
	}

	// We are looking for the "web" directory
	targetFolder := "web"
	targetFolderPath := ""

	for {
		if _, err := os.Stat(workingDirectory + "/" + targetFolder); err == nil {
			return targetFolderPath + targetFolder
		}

		targetFolderPath = targetFolderPath + "../"
		parentDirectory := workingDirectory + "/.."
		if parentDirectory == workingDirectory {
			// We have reached the root directory
			break
		}

		workingDirectory = parentDirectory
	}

	panic("Failed to find the root directory containing the 'web' folder")
}
