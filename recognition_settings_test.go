package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"
)

func TestCropImageRectReturnsRequestedArea(t *testing.T) {
	source := image.NewNRGBA(image.Rect(10, 20, 14, 24))
	fillTestRect(source, source.Bounds(), color.NRGBA{R: 10, G: 20, B: 30, A: 255})
	fillTestRect(source, image.Rect(11, 21, 13, 23), color.NRGBA{R: 200, G: 100, B: 50, A: 255})

	cropped, err := cropImageRect(source, image.Rect(11, 21, 13, 23))
	if err != nil {
		t.Fatalf("cropImageRect returned error: %v", err)
	}

	if cropped.Bounds().Dx() != 2 || cropped.Bounds().Dy() != 2 {
		t.Fatalf("expected 2x2 crop, got %v", cropped.Bounds())
	}

	assertPixelColor(t, cropped, 0, 0, color.NRGBA{R: 200, G: 100, B: 50, A: 255})
	assertPixelColor(t, cropped, 1, 1, color.NRGBA{R: 200, G: 100, B: 50, A: 255})
}

func TestCropImageRectClampsToVisibleIntersection(t *testing.T) {
	source := image.NewNRGBA(image.Rect(100, 100, 104, 104))
	fillTestRect(source, source.Bounds(), color.NRGBA{R: 1, G: 2, B: 3, A: 255})

	cropped, err := cropImageRect(source, image.Rect(102, 102, 106, 106))
	if err != nil {
		t.Fatalf("cropImageRect returned error: %v", err)
	}

	if cropped.Bounds().Dx() != 2 || cropped.Bounds().Dy() != 2 {
		t.Fatalf("expected visible intersection to be 2x2, got %v", cropped.Bounds())
	}

	assertPixelColor(t, cropped, 0, 0, color.NRGBA{R: 1, G: 2, B: 3, A: 255})
	assertPixelColor(t, cropped, 1, 1, color.NRGBA{R: 1, G: 2, B: 3, A: 255})
}

func TestCropImageRectReturnsErrorWhenNoVisibleIntersection(t *testing.T) {
	source := image.NewNRGBA(image.Rect(0, 0, 10, 10))

	_, err := cropImageRect(source, image.Rect(20, 20, 30, 30))
	if err == nil {
		t.Fatal("expected cropImageRect to fail when target has no visible intersection")
	}
}

func TestGetRecruitmentTagCatalogIncludesStandardGroups(t *testing.T) {
	app := NewApp()
	catalog := app.GetRecruitmentTagCatalog()
	if len(catalog.Groups) != 4 {
		t.Fatalf("expected 4 recruitment tag groups, got %d", len(catalog.Groups))
	}
	if catalog.Groups[0].Label != "职业" {
		t.Fatalf("expected first group to be 职业, got %s", catalog.Groups[0].Label)
	}
	if !isRecruitmentTag("近卫") || !isRecruitmentTag("高级资深干员") {
		t.Fatal("expected standard recruitment tags to be registered")
	}
}

func TestDecodeRecognitionTemplateSupportsLegacySingleStateRegions(t *testing.T) {
	payload, err := json.Marshal(legacyRecognitionTemplate{
		ID:             "legacy-template",
		Hwnd:           1,
		Title:          "legacy",
		ClassName:      "window",
		Width:          100,
		Height:         100,
		ScreenshotPath: "window.png",
		CreatedAt:      "2026-03-25T00:00:00Z",
		Regions: []legacyRecognitionRegion{
			{
				ID:            "region-01",
				Label:         "近卫",
				X:             0.1,
				Y:             0.2,
				Width:         0.3,
				Height:        0.4,
				ReferencePath: "regions/region-01.png",
			},
		},
	})
	if err != nil {
		t.Fatalf("failed to marshal legacy payload: %v", err)
	}

	template, err := decodeRecognitionTemplate(payload)
	if err != nil {
		t.Fatalf("decodeRecognitionTemplate returned error: %v", err)
	}
	if len(template.Regions) != 1 || len(template.Regions[0].States) != 1 {
		t.Fatalf("expected legacy template to normalize to one state, got %+v", template.Regions)
	}
	if template.Regions[0].States[0].Tag != "近卫" {
		t.Fatalf("expected legacy tag to be preserved when valid, got %s", template.Regions[0].States[0].Tag)
	}
}

func TestMatchRecognitionRegionStatesReturnsMatchedTags(t *testing.T) {
	templateDir := t.TempDir()
	img := image.NewNRGBA(image.Rect(0, 0, 2, 2))
	fillTestRect(img, img.Bounds(), color.NRGBA{R: 255, G: 0, B: 0, A: 255})
	imgBytes := mustEncodePNGForTest(t, img)

	if err := os.MkdirAll(filepath.Join(templateDir, "regions", "region-01"), 0o755); err != nil {
		t.Fatalf("failed to create region dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(templateDir, "regions", "region-01", "state-01.png"), imgBytes, 0o644); err != nil {
		t.Fatalf("failed to write state image: %v", err)
	}

	matched := matchRecognitionRegionStates(templateDir, img, []RecognitionRegionState{
		{ID: "state-01", Tag: "近卫", ReferencePath: "regions/region-01/state-01.png"},
		{ID: "state-02", Tag: "狙击", ReferencePath: "regions/region-01/missing.png"},
	})
	if len(matched) != 1 {
		t.Fatalf("expected one matched state, got %#v", matched)
	}
	if matched[0].Tag != "近卫" {
		t.Fatalf("expected matched tag 近卫, got %s", matched[0].Tag)
	}
}

func TestNormalizeRegionInputRejectsMissingTaggedState(t *testing.T) {
	_, err := normalizeRegionInput(RecognitionRegionInput{
		ID:     "region-01",
		Label:  "测试区域",
		X:      0,
		Y:      0,
		Width:  0.5,
		Height: 0.5,
		States: []RecognitionRegionStateInput{{ID: "state-01", Tag: "", ImagePNG: base64.StdEncoding.EncodeToString([]byte("x"))}},
	})
	if err == nil {
		t.Fatal("expected normalizeRegionInput to reject state without recruitment tag")
	}
}

func fillTestRect(img *image.NRGBA, rect image.Rectangle, fill color.NRGBA) {
	visibleRect := img.Bounds().Intersect(rect)
	for y := visibleRect.Min.Y; y < visibleRect.Max.Y; y++ {
		for x := visibleRect.Min.X; x < visibleRect.Max.X; x++ {
			img.SetNRGBA(x, y, fill)
		}
	}
}

func assertPixelColor(t *testing.T, img *image.NRGBA, x, y int, expected color.NRGBA) {
	t.Helper()

	actual := img.NRGBAAt(x, y)
	if actual != expected {
		t.Fatalf("expected pixel (%d,%d) to be %#v, got %#v", x, y, expected, actual)
	}
}

func mustEncodePNGForTest(t *testing.T, img image.Image) []byte {
	t.Helper()

	var buffer bytes.Buffer
	if err := png.Encode(&buffer, img); err != nil {
		t.Fatalf("failed to encode png: %v", err)
	}
	return buffer.Bytes()
}
