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

type Exts []string
type ImgFmtExt map[string]Exts

var (
	fmtFrom   string
	fmtTo     string
	imgFmtExt ImgFmtExt
)

func initExt() ImgFmtExt {
	return ImgFmtExt{
		"jpeg": Exts{"jpg", "jpeg"},
		"png":  Exts{"png"},
		"gif":  Exts{"gif"},
	}
}

func Convert(dirName string, from string, to string) {
	if dirName == "" {
		panic("Directory name is not provided.")
	}

	imgFmtExt = initExt()
	fmtFrom = convFromExtToImgFmt(from)
	fmtTo = convFromExtToImgFmt(to)
	if fmtFrom == "" || fmtTo == "" {
		panic("Given image format is not supported.")
	}
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
	file, err := os.Create(generateOutputPath(path, fmtTo))
	if err != nil {
		panic(err)
	}
	defer file.Close()

	switch fmtTo {
	case "jpeg":
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

func generateOutputPath(path string, fmtTo string) string {
	// TODO: consider the case when the file of the same name already exists
	base := strings.TrimRight(filepath.Base(path), filepath.Ext(path))
	ext := convFromImgFmtToExt(fmtTo)
	return strings.Join([]string{base, ext}, ".")
}

func convFromExtToImgFmt(ext string) string {
	ext = strings.ToLower(ext)
	for imgFmt, fmtExts := range imgFmtExt {
		for _, fmtExt := range fmtExts {
			if ext == fmtExt {
				return imgFmt
			}
		}
	}
	return ""
}

func convFromImgFmtToExt(fmt string) string {
	fmt = strings.ToLower(fmt)
	for imgFmt, fmtExts := range imgFmtExt {
		if fmt == imgFmt {
			return fmtExts[0]
		}
	}
	return ""
}
