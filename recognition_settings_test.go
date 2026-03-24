package main

import (
	"image"
	"image/color"
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
