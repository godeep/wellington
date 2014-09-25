package sprite_sass

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSpriteLookup(t *testing.T) {

	imgs := ImageList{}
	imgs.Decode("test/139.jpg", "test/140.jpg")
	if f := imgs.Lookup("test/139.jpg"); f != 0 {
		t.Errorf("Invalid file location given found %d, expected %d", f, 0)
	}

	if f := imgs.Lookup("test/140.jpg"); f != 1 {
		t.Errorf("Invalid file location given found %d, expected %d", f, 1)
	}

	if f := imgs.Lookup("140"); f != 1 {
		t.Errorf("Invalid file location given found %d, expected %d", f, 1)
	}

	if f := imgs.Lookup("notafile.jpg"); f != -1 {
		t.Errorf("Found a file that doesn't exist")
	}
}

func TestSpriteCombine(t *testing.T) {
	imgs := ImageList{}
	imgs.Decode("test/139.jpg", "test/140.jpg")
	imgs.Vertical = true
	imgs.Combine()

	if imgs.Height() != 279 {
		t.Errorf("Invalid Height found %d, wanted %d", imgs.Height(), 279)
	}

	if imgs.Width() != 96 {
		t.Errorf("Invalid Width found %d, wanted %d", imgs.Width(), 192)
	}

	if x := imgs.X(1); x != 0 {
		t.Errorf("Invalid X found %d, wanted %d", x, 0)
	}

	if y := imgs.Y(1); y != 139 {
		t.Errorf("Invalid Y found %d, wanted %d", y, 139)
	}

	if e := -1; e != imgs.ImageWidth("150") {
		t.Errorf("Non-existant image width invalid"+
			"\n    was:%d\nexpected:%d",
			imgs.ImageWidth("150"), e)
	}

	if e := -1; e != imgs.ImageHeight("150") {
		t.Errorf("Non-existant image width invalid"+
			"\n    was:%d\nexpected:%d",
			imgs.ImageHeight("150"), e)
	}

	if e := ""; e != imgs.Dimensions("150") {
		t.Errorf("Non-existant image width invalid"+
			"\n    was:%s\nexpected:%s",
			imgs.Dimensions("150"), e)
	}

	//Quick cache check
	imgs.Combine()
	if imgs.Height() != 279 || imgs.Width() != 96 {
		t.Errorf("Cache invalid")
	}

	testFile, err := imgs.Export()

	defer func() {
		//Cleanup test files
		err := os.Remove(testFile)

		if err != nil {
			panic(err)
		}

	}()

	if err != nil {
		panic(err)
	}
}

// Test image dimension calls
func TestSpriteImageDimensions(t *testing.T) {
	imgs := ImageList{}
	imgs.Decode("test/*.png")

	if e := "width: 96px;\nheight: 139px"; e != imgs.Dimensions("139") {
		t.Errorf("Dimensions invalid was: %s\nexpected: %s\n",
			imgs.Dimensions("139"), e)
	}

	if e := 139; e != imgs.ImageHeight("139") {
		t.Errorf("Height invalid    was:%d\nexpected:%d",
			imgs.ImageHeight("139"), e)
	}

	if e := 96; e != imgs.ImageWidth("139") {
		t.Errorf("Height invalid was:%d\nexpected:%d",
			imgs.ImageWidth("139"), e)
	}

	if e := "-96px 0px"; imgs.Position("140") != e {
		t.Errorf("Invalid position found was: %s\nexpected:%s",
			imgs.Position("140"), e)
	}

	output := rerandom.ReplaceAllString(imgs.CSS("140"), "")
	if e := `url("./test") -96px 0px`; output != e {
		t.Errorf("Invalid CSS generated on test     was: %s\nexpected: %s",
			output, e)
	}

}

//Test file globbing
func TestSpriteGlob(t *testing.T) {
	imgs := ImageList{
		ImageDir: "test",
	}
	imgs.Decode("*.png")

	if f := imgs.Lookup("test/139.png"); f != 0 {
		t.Errorf("Invalid file location given found %d, expected %d", f, 0)
	}

	if f := imgs.Lookup("test/140.png"); f != 1 {
		t.Errorf("Invalid file location given found %d, expected %d", f, 1)
	}

	if f := imgs.Lookup("notafile.png"); f != -1 {
		t.Errorf("Found a file that doesn't exist")
	}
	outpath := rerandom.ReplaceAllString(imgs.OutFile, "")
	outfile := filepath.Base(outpath)
	if e := "image"; e != outfile {
		t.Errorf("Outfile misnamed \n     was: %s\nexpected: %s", outpath, e)
	}
	ext := filepath.Ext(imgs.OutFile)
	if e := ".png"; e != ext {
		t.Errorf("Outfile invalid extension\n    was: %s\nexpected: %s",
			ext, e)
	}
	imgs = ImageList{
		ImageDir:  ".",
		GenImgDir: "build/test",
	}
	imgs.Decode("test/*.png")
	outpath = rerandom.ReplaceAllString(imgs.OutFile, "")
	outfile = filepath.Base(outpath)

	if e := "test"; e != outfile {
		t.Errorf("Outfile misnamed \n     was: %s\nexpected: %s", outpath, e)
	}
	ext = filepath.Ext(imgs.OutFile)
	if e := ".png"; e != ext {
		t.Errorf("Outfile invalid extension\n    was: %s\nexpected: %s",
			ext, e)
	}
	if e := "build/test/test"; e != outpath {
		t.Errorf("Invalid path\n     was: %s\nexpected: %s", outpath, e)
	}
}

func TestSpriteDecode(t *testing.T) {
	//Should fail with unable to find file
	i := ImageList{}
	err := i.Decode("notafile")

	if err != nil {
		t.Errorf("Error thrown for non-existant file")
	}

	if len(i.Images) > 0 {
		t.Errorf("Found a non-existant file")
	}
}

func TestSpriteHorizontal(t *testing.T) {

	imgs := ImageList{}
	imgs.Decode("test/139.jpg", "test/140.jpg")
	imgs.Vertical = false
	imgs.Combine()

	if e := 140; imgs.Height() != e {
		t.Errorf("Invalid Height found %d, wanted %d", imgs.Height(), e)
	}

	if e := 192; imgs.Width() != e {
		t.Errorf("Invalid Width found %d, wanted %d", imgs.Width(), e)
	}

	if e := 96; imgs.X(1) != e {
		t.Errorf("Invalid X found %d, wanted %d", imgs.X(1), e)
	}

	if e := 0; imgs.Y(1) != e {
		t.Errorf("Invalid Y found %d, wanted %d", imgs.Y(1), e)
	}
}

func TestSpriteInline(t *testing.T) {
	imgs := ImageList{}
	imgs.Decode("test/pixel.png")
	imgs.Combine()
	bytes := imgs.inline()

	// Bytes are non-deterministic, so check length and move on
	if len(bytes) < 300 || len(bytes) > 350 {
		t.Errorf("Pixel blog data had an invalid length"+
			"\n     was: %d\nexpected: 300-350", len(bytes))
	}

	str := imgs.Inline()
	if len(str) < 400 || len(str) > 500 {
		t.Errorf("CSS length has an invalid length:%d expected: 400-500",
			len(str))
	}
}
