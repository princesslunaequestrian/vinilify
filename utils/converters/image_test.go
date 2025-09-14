package converters

import (
	"image"
	"path/filepath"
	"testing"

	"github.com/fogleman/gg"
	"github.com/nfnt/resize"
	"github.com/princesslunaequestrian/vinilify/utils"
)

func TestCropAndRotateImage(t *testing.T) {
	testAssetsPath := filepath.Join(utils.GetRoot(), "utils", "converters", "test_assets")
	testResultsPath := filepath.Join(utils.GetRoot(), "utils", "converters", "test_results")

	img, err := gg.LoadImage(filepath.Join(testAssetsPath, "t.png"))
	if err != nil {
		t.Error(err)
	}

	img = resize.Resize(1000, 1000, CropAndRotateImage(img, float64(50)), resize.Lanczos3)
	dc := gg.NewContext(1000, 1000)
	dc.DrawImage(img, 0, 0)
	dc.SavePNG(filepath.Join(testResultsPath, "test1.png"))

}

func TestStackImage(t *testing.T) {
	testAssetsPath := filepath.Join(utils.GetRoot(), "utils", "converters", "test_assets")
	testResultsPath := filepath.Join(utils.GetRoot(), "utils", "converters", "test_results")

	img1, err := gg.LoadImage(filepath.Join(testAssetsPath, "t.png"))
	if err != nil {
		t.Error(err)
	}

	back1, _ := gg.LoadImage(filepath.Join(utils.GetAssets(), "Images", "Disk.png"))

	b1s := resize.Resize(1000, 1000, back1, resize.Lanczos3)
	i1s := resize.Resize(500, 500, img1, resize.Lanczos3)

	dc := gg.NewContext(1000, 1000)
	dc.DrawImage(i1s, 250, 250)
	dc.DrawImage(b1s, 0, 0)
	dc.SavePNG(filepath.Join(testResultsPath, "test2.png"))
}

func TestStackImages(t *testing.T) {
	testAssetsPath := filepath.Join(utils.GetRoot(), "utils", "converters", "test_assets")
	testResultsPath := filepath.Join(utils.GetRoot(), "utils", "converters", "test_results")

	img1, _ := LoadAndResizeImage(filepath.Join(testAssetsPath, "t.png"), 500, 500)
	img2, _ := LoadAndResizeImage(filepath.Join(testResultsPath, "test.png"), 250, 250)
	img3, _ := LoadAndResizeImage(filepath.Join(utils.GetAssets(), "Images", "Disk.png"), 1000, 1000)

	imgs := []image.Image{
		img3,
		img1,
		img2,
	}

	coords := []Coord{
		{x: 0, y: 0},
		{x: 250, y: 250},
		{x: 375, y: 375},
	}

	dc, err := StackImages(gg.NewContext(1000, 1000), imgs, coords)
	if err != nil {
		t.Error(err)
	}
	dc.SavePNG(filepath.Join(testResultsPath, "test_stack.png"))
}

func TestAssembleImage(t *testing.T) {
	testInputPath := "../../users/286946560/image.jpg"
	testOutPath := "../../users/286946560"

	err := AssembleImages(testInputPath, testOutPath)
	if err != nil {
		t.Error(err)
	}
}
