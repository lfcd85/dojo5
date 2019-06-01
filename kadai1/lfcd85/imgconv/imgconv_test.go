package imgconv

import (
	"testing"
)

func assertEq(t *testing.T, actual interface{}, expected interface{}) {
	if actual != expected {
		t.Errorf("actual: %v, expected: %v", actual, expected)
	}
}

/* NOTE: This testing can work, but should be run in a sandbox.

func assertNil(t *testing.T, obj interface{}) {
	if obj != nil {
		t.Errorf("actual: not nil, expected: nil")
	}
}

func assertNotNil(t *testing.T, obj interface{}) {
	if obj == nil {
		t.Errorf("actual: nil, expected: not nil")
	}
}

func TestConvert(t *testing.T) {
	cases := []struct {
		from     string
		to       string
		expected bool
	}{
		{"jpeg", "png", true},
		{"png", "gif", true},
		{"gif", "jpeg", true},
		{"jpeg", "jpeg", false},
		{"rb", "go", false},
	}

	for _, c := range cases {
		err := Convert("../test", c.from, c.to)
		if c.expected == true {
			assertNil(t, err)
		} else {
			assertNotNil(t, err)
		}
	}
}

*/

func TestDetectImgFmts(t *testing.T) {
	cases := []struct {
		from            string
		to              string
		expectedFmtFrom ImgFmt
		expectedFmtTo   ImgFmt
	}{
		{
			"PNG",
			"jpg",
			ImgFmt("png"),
			ImgFmt("jpeg"),
		},
		{
			"gif",
			"JPEG",
			ImgFmt("gif"),
			ImgFmt("jpeg"),
		},
	}

	for _, c := range cases {
		detectImgFmts(c.from, c.to)
		assertEq(t, fmtFrom, c.expectedFmtFrom)
		assertEq(t, fmtTo, c.expectedFmtTo)
	}
}

func TestCheckExt(t *testing.T) {
	cases := []struct {
		fileName string
		fmtFrom  ImgFmt
		expected bool
	}{
		{"hoge.jpg", ImgFmt("jpeg"), true},
		{"fuga.png", ImgFmt("gif"), false},
		{"piyo.png", ImgFmt("png"), true},
		{"foo.js", ImgFmt("png"), false},
		{".JPEG", ImgFmt("jpeg"), true},
		{"jpeg", ImgFmt("jpeg"), false},
		{"foopng", ImgFmt("png"), false},
	}

	initExts()
	for _, c := range cases {
		fmtFrom = c.fmtFrom
		assertEq(t, checkExt(c.fileName), c.expected)
	}
}

func TestGenerateOutputPath(t *testing.T) {
	cases := []struct {
		path     string
		fmtTo    ImgFmt
		expected string
	}{
		{
			"path/to/hoge.jpg",
			ImgFmt("png"),
			"path/to/hoge.png",
		},
		{
			"./path/to/fuga.PNG",
			ImgFmt("jpeg"),
			"./path/to/fuga.jpg",
		},
		{
			"piyo.png",
			ImgFmt("gif"),
			"piyo.gif",
		},
	}

	initExts()
	for _, c := range cases {
		fmtTo = c.fmtTo
		assertEq(t, generateOutputPath(c.path), c.expected)
	}
}
