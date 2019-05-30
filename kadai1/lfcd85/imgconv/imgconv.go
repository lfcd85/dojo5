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

type ImgFmt string
type Ext string
type Exts []Ext
type ImgFmtExts map[ImgFmt]Exts

var (
	fmtFrom    ImgFmt
	fmtTo      ImgFmt
	imgFmtExts ImgFmtExts
)

func initExts() ImgFmtExts {
	return ImgFmtExts{
		"jpeg": Exts{"jpg", "jpeg"},
		"png":  Exts{"png"},
		"gif":  Exts{"gif"},
	}
}

func Convert(dirName string, from string, to string) {
	if dirName == "" {
		panic("Directory name is not provided.")
	}

	imgFmtExts = initExts()
	extFrom := Ext(strings.ToLower(from))
	extTo := Ext(strings.ToLower(to))
	fmtFrom = convFromExtToImgFmt(extFrom)
	fmtTo = convFromExtToImgFmt(extTo)
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
		if isFmtFrom := checkExt(info.Name(), fmtFrom); !info.IsDir() && isFmtFrom {
			convSingleFile(path)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
}

func checkExt(fileName string, fmtFrom ImgFmt) bool {
	fileExt := Ext(strings.TrimPrefix(filepath.Ext(fileName), "."))
	fileImgFmt := convFromExtToImgFmt(fileExt)
	return fileImgFmt == fmtFrom
}

func convSingleFile(path string) {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	img, fmtStr, err := image.Decode(file)
	if err != nil {
		fmt.Printf("%q will not be converted (%v)\n", path, err)
	}

	if ImgFmt(fmtStr) == fmtFrom {
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

func generateOutputPath(path string, fmtTo ImgFmt) string {
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
