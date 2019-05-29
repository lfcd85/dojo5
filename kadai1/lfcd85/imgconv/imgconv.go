package imgconv

import (
	"fmt"
	"image"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"strings"

	_ "image/jpeg"
)

func Convert(dirName string) {
	err := filepath.Walk(dirName, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			convSingleFile(path)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
}

func convSingleFile(path string) {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	img, format, err := image.Decode(file)
	if err != nil {
		fmt.Printf("%q will not be converted (%v)\n", path, err)
	}

	if format == "jpeg" {
		writeOutputFile(img, path)
	}
}

func detectImageFormat(r io.Reader) (string, error) {
	_, format, err := image.DecodeConfig(r)
	return format, err
}

func writeOutputFile(img image.Image, path string) {
	ext := "png"
	file, err := os.Create(generateOutputPath(path, ext))
	if err != nil {
		panic(err)
	}
	defer file.Close()

	if err := png.Encode(file, img); err != nil {
		panic(err)
	}
}

func generateOutputPath(path string, ext string) string {
	// TODO: consider the case when the file of the same name already exists
	base := strings.TrimRight(filepath.Base(path), filepath.Ext(path))
	return strings.Join([]string{base, ext}, ".")
}
