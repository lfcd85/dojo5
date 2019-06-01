// Package imgconv provides a recursive conversion of images in the directory.
package imgconv

import (
	"errors"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"
)

// ImgFmtExts is a map of the image formats and its extensions
type ImgFmtExts map[ImgFmt]Exts

// Exts is a slice of image extensions
type Exts []Ext

// Ext is a image extension
type Ext string

// ImgFmt is a image format
type ImgFmt string

var (
	fmtFrom    ImgFmt
	fmtTo      ImgFmt
	imgFmtExts ImgFmtExts
)

// Convert recursively seeks a given directory and converts images from and to given formats.
func Convert(dirName string, from string, to string) error {
	if dirName == "" {
		return errors.New("Directory name is not provided.")
	}

	detectImgFmts(from, to)
	if fmtFrom == "" || fmtTo == "" {
		return errors.New("Given image format is not supported.")
	}
	if fmtFrom == fmtTo {
		return errors.New("Image formats before and after conversion are the same.")
	}

	err := filepath.Walk(dirName, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		err = convSingleFile(path, info)
		return err
	})
	return err
}

func initExts() {
	imgFmtExts = ImgFmtExts{
		"jpeg": Exts{"jpg", "jpeg"},
		"png":  Exts{"png"},
		"gif":  Exts{"gif"},
	}
}

func detectImgFmts(from string, to string) {
	initExts()
	extFrom := Ext(strings.ToLower(from))
	extTo := Ext(strings.ToLower(to))
	fmtFrom = convFromExtToImgFmt(extFrom)
	fmtTo = convFromExtToImgFmt(extTo)
}

func checkExt(fileName string) bool {
	fileExtStr := strings.TrimPrefix(filepath.Ext(fileName), ".")
	fileExt := Ext(strings.ToLower(fileExtStr))
	fileImgFmt := convFromExtToImgFmt(fileExt)
	return fileImgFmt == fmtFrom
}

func convSingleFile(path string, info os.FileInfo) error {
	if isFmtFrom := checkExt(info.Name()); info.IsDir() || !isFmtFrom {
		return nil
	}

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	img, fmtStr, err := image.Decode(file)
	if err != nil {
		fmt.Printf("%q is skipped (%v)\n", path, err)
		return nil
	}
	if ImgFmt(fmtStr) != fmtFrom {
		return nil
	}

	err = writeOutputFile(img, path)
	return err
}

func writeOutputFile(img image.Image, path string) error {
	file, err := os.Create(generateOutputPath(path))
	if err != nil {
		return err
	}
	defer file.Close()

	switch fmtTo {
	case "jpeg":
		if err := jpeg.Encode(file, img, nil); err != nil {
			return err
		}
	case "png":
		if err := png.Encode(file, img); err != nil {
			return err
		}
	case "gif":
		if err := gif.Encode(file, img, nil); err != nil {
			return err
		}
	}
	return nil
}

func generateOutputPath(path string) string {
	dirAndBase := strings.TrimRight(path, filepath.Ext(path))
	ext := convFromImgFmtToExt(fmtTo)
	return strings.Join([]string{dirAndBase, string(ext)}, ".")
}

func convFromExtToImgFmt(ext Ext) ImgFmt {
	for imgFmt, fmtExts := range imgFmtExts {
		for _, fmtExt := range fmtExts {
			if ext == fmtExt {
				return imgFmt
			}
		}
	}
	return ""
}

func convFromImgFmtToExt(fmt ImgFmt) Ext {
	for imgFmt, fmtExts := range imgFmtExts {
		if fmt == imgFmt {
			return fmtExts[0]
		}
	}
	return ""
}
