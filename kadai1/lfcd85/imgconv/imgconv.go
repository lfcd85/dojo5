package imgconv

import (
	"image"
	"image/png"
	"os"

	_ "image/jpeg"
)

func Convert(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		panic(err)
	}
	writeOutputFile(img)
}

func writeOutputFile(img image.Image) {
	file, err := os.Create("output.png")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	if err := png.Encode(file, img); err != nil {
		panic(err)
	}
}
