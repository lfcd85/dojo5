package imgconv

import (
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"
)

var (
	fmtFrom string
	fmtTo   string
)

func Convert(dirName string, from string, to string) {
	if dirName == "" {
		panic("Directory name is not provided.")
	}

	fmtFrom = detectImgFmt(from)
	fmtTo = detectImgFmt(to)
	if fmtFrom == fmtTo {
		panic("Image formats before and after conversion are the same.")
	}

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

func detectImgFmt(fmt string) string {
	var detectedFmt string

	switch fmt = strings.ToLower(fmt); fmt {
	case "jpg", "jpeg":
		detectedFmt = "jpeg"
	case "png":
		detectedFmt = "png"
	case "gif":
		detectedFmt = "gif"
	default:
		panic("Given image format is not supported.")
	}

	return detectedFmt
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

	if format == fmtFrom {
		writeOutputFile(img, path)
	}
}

func writeOutputFile(img image.Image, path string) {
	ext := fmtTo
	file, err := os.Create(generateOutputPath(path, ext))
	if err != nil {
		panic(err)
	}
	defer file.Close()

	switch fmtTo {
	case "jpg", "jpeg":
		if err := jpeg.Encode(file, img, nil); err != nil {
			panic(err)
		}
	case "png":
		if err := png.Encode(file, img); err != nil {
			panic(err)
		}
	case "gif":
		if err := gif.Encode(file, img, nil); err != nil {
			panic(err)
		}
	}
}

func generateOutputPath(path string, ext string) string {
	// TODO: consider the case when the file of the same name already exists
	base := strings.TrimRight(filepath.Base(path), filepath.Ext(path))
	return strings.Join([]string{base, ext}, ".")
}
