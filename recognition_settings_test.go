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
	"time"
)

func TestResolveRecognitionWindowInstancePrefersRecordedHandle(t *testing.T) {
	originalAlive := isWindowHandleAliveFunc
	originalCandidates := resolveRecognitionWindowCandidatesFunc
	defer func() {
		isWindowHandleAliveFunc = originalAlive
		resolveRecognitionWindowCandidatesFunc = originalCandidates
	}()

	candidatesCalled := false
	isWindowHandleAliveFunc = func(hwnd uintptr) bool {
		return hwnd == 1001
	}
	resolveRecognitionWindowCandidatesFunc = func(title string, className string) ([]RecognitionWindowCandidate, error) {
		candidatesCalled = true
		return nil, nil
	}

	hwnd, failureReason, err := resolveRecognitionWindowInstance(RecognitionTemplate{
		Hwnd:      1001,
		Title:     "target",
		ClassName: "game-window",
	})
	if err != nil {
		t.Fatalf("resolveRecognitionWindowInstance returned error: %v", err)
	}
	if failureReason != publicRecruitmentRecognitionFailureNone {
		t.Fatalf("expected no failure reason, got %s", failureReason)
	}
	if hwnd != 1001 {
		t.Fatalf("expected recorded handle to be reused, got %d", hwnd)
	}
	if candidatesCalled {
		t.Fatal("expected candidate resolution to be skipped when recorded handle is still alive")
	}
}

func TestResolveRecognitionWindowInstanceFallsBackToUniqueMetadataMatch(t *testing.T) {
	originalAlive := isWindowHandleAliveFunc
	originalCandidates := resolveRecognitionWindowCandidatesFunc
	defer func() {
		isWindowHandleAliveFunc = originalAlive
		resolveRecognitionWindowCandidatesFunc = originalCandidates
	}()

	isWindowHandleAliveFunc = func(hwnd uintptr) bool {
		return false
	}
	resolveRecognitionWindowCandidatesFunc = func(title string, className string) ([]RecognitionWindowCandidate, error) {
		return []RecognitionWindowCandidate{
			{
				Hwnd:      2001,
				Title:     title,
				ClassName: className,
				ProcessID: 10,
				Bounds:    RecognitionWindowBounds{Left: 1, Top: 2, Right: 101, Bottom: 202},
			},
			{
				Hwnd:      2002,
				Title:     title,
				ClassName: className,
				ProcessID: 11,
				Bounds:    RecognitionWindowBounds{Left: 5, Top: 6, Right: 105, Bottom: 206},
			},
		}, nil
	}

	hwnd, failureReason, err := resolveRecognitionWindowInstance(RecognitionTemplate{
		Hwnd:      1001,
		Title:     "target",
		ClassName: "game-window",
		Instance: RecognitionWindowInstanceMetadata{
			ProcessID: 11,
			Bounds:    RecognitionWindowBounds{Left: 5, Top: 6, Right: 105, Bottom: 206},
		},
	})
	if err != nil {
		t.Fatalf("resolveRecognitionWindowInstance returned error: %v", err)
	}
	if failureReason != publicRecruitmentRecognitionFailureNone {
		t.Fatalf("expected no failure reason, got %s", failureReason)
	}
	if hwnd != 2002 {
		t.Fatalf("expected fallback to select unique candidate, got %d", hwnd)
	}
}

func TestResolveRecognitionWindowInstanceRejectsAmbiguousCandidates(t *testing.T) {
	originalAlive := isWindowHandleAliveFunc
	originalCandidates := resolveRecognitionWindowCandidatesFunc
	defer func() {
		isWindowHandleAliveFunc = originalAlive
		resolveRecognitionWindowCandidatesFunc = originalCandidates
	}()

	isWindowHandleAliveFunc = func(hwnd uintptr) bool {
		return false
	}
	resolveRecognitionWindowCandidatesFunc = func(title string, className string) ([]RecognitionWindowCandidate, error) {
		return []RecognitionWindowCandidate{
			{
				Hwnd:      3001,
				Title:     title,
				ClassName: className,
				ProcessID: 99,
				Bounds:    RecognitionWindowBounds{Left: 10, Top: 10, Right: 210, Bottom: 210},
			},
			{
				Hwnd:      3002,
				Title:     title,
				ClassName: className,
				ProcessID: 99,
				Bounds:    RecognitionWindowBounds{Left: 10, Top: 10, Right: 210, Bottom: 210},
			},
		}, nil
	}

	_, failureReason, err := resolveRecognitionWindowInstance(RecognitionTemplate{
		Title:     "target",
		ClassName: "game-window",
		Instance: RecognitionWindowInstanceMetadata{
			ProcessID: 99,
			Bounds:    RecognitionWindowBounds{Left: 10, Top: 10, Right: 210, Bottom: 210},
		},
	})
	if err == nil {
		t.Fatal("expected ambiguous candidates to return an error")
	}
	if failureReason != publicRecruitmentRecognitionFailureAmbiguousWindow {
		t.Fatalf("expected ambiguous window failure, got %s", failureReason)
	}
}

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
	if catalog.Groups[0].Label != "\u804c\u4e1a" {
		t.Fatalf("expected first group to be profession group, got %s", catalog.Groups[0].Label)
	}
	if !isRecruitmentTag("\u8fd1\u536b") || !isRecruitmentTag("\u9ad8\u7ea7\u8d44\u6df1\u5e72\u5458") {
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
				Label:         "\u8fd1\u536b",
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
	if template.Regions[0].States[0].Tolerance != 0 {
		t.Fatalf("expected legacy state tolerance to default to 0, got %d", template.Regions[0].States[0].Tolerance)
	}
	if template.Regions[0].States[0].Tag != "\u8fd1\u536b" {
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
		{ID: "state-01", Tag: "\u8fd1\u536b", Tolerance: 0, ReferencePath: "regions/region-01/state-01.png"},
		{ID: "state-02", Tag: "\u72d9\u51fb", Tolerance: 0, ReferencePath: "regions/region-01/missing.png"},
	})
	if len(matched) != 1 {
		t.Fatalf("expected one matched state, got %#v", matched)
	}
	if matched[0].Tag != "\u8fd1\u536b" {
		t.Fatalf("expected matched guard tag, got %s", matched[0].Tag)
	}
}

func TestMatchRecognitionRegionStatesUsesPerStateTolerance(t *testing.T) {
	templateDir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(templateDir, "regions", "region-01"), 0o755); err != nil {
		t.Fatalf("failed to create region dir: %v", err)
	}

	current := image.NewNRGBA(image.Rect(0, 0, 1, 1))
	current.SetNRGBA(0, 0, color.NRGBA{R: 100, G: 120, B: 140, A: 255})

	withinTolerance := image.NewNRGBA(image.Rect(0, 0, 1, 1))
	withinTolerance.SetNRGBA(0, 0, color.NRGBA{R: 102, G: 121, B: 143, A: 255})
	if err := os.WriteFile(filepath.Join(templateDir, "regions", "region-01", "state-01.png"), mustEncodePNGForTest(t, withinTolerance), 0o644); err != nil {
		t.Fatalf("failed to write first state image: %v", err)
	}

	outsideTolerance := image.NewNRGBA(image.Rect(0, 0, 1, 1))
	outsideTolerance.SetNRGBA(0, 0, color.NRGBA{R: 106, G: 120, B: 140, A: 255})
	if err := os.WriteFile(filepath.Join(templateDir, "regions", "region-01", "state-02.png"), mustEncodePNGForTest(t, outsideTolerance), 0o644); err != nil {
		t.Fatalf("failed to write second state image: %v", err)
	}

	matched := matchRecognitionRegionStates(templateDir, current, []RecognitionRegionState{
		{ID: "state-01", Tag: "\u8fd1\u536b", Tolerance: 4, ReferencePath: "regions/region-01/state-01.png"},
		{ID: "state-02", Tag: "\u72d9\u51fb", Tolerance: 4, ReferencePath: "regions/region-01/state-02.png"},
	})
	if len(matched) != 1 {
		t.Fatalf("expected one matched state with per-state tolerance, got %#v", matched)
	}
	if matched[0].StateID != "state-01" {
		t.Fatalf("expected only first state to match, got %#v", matched)
	}
}

func TestCompareImagesWithTolerance(t *testing.T) {
	base := image.NewNRGBA(image.Rect(0, 0, 1, 1))
	base.SetNRGBA(0, 0, color.NRGBA{R: 100, G: 120, B: 140, A: 255})

	exact := image.NewNRGBA(image.Rect(0, 0, 1, 1))
	exact.SetNRGBA(0, 0, color.NRGBA{R: 100, G: 120, B: 140, A: 255})

	withinTolerance := image.NewNRGBA(image.Rect(0, 0, 1, 1))
	withinTolerance.SetNRGBA(0, 0, color.NRGBA{R: 102, G: 121, B: 143, A: 255})

	outsideTolerance := image.NewNRGBA(image.Rect(0, 0, 1, 1))
	outsideTolerance.SetNRGBA(0, 0, color.NRGBA{R: 106, G: 120, B: 140, A: 255})

	if !compareImages(base, exact, 0) {
		t.Fatal("expected zero tolerance to keep exact-match behavior")
	}
	if compareImages(base, withinTolerance, 0) {
		t.Fatal("expected zero tolerance to reject non-identical pixels")
	}
	if !compareImages(base, withinTolerance, 4) {
		t.Fatal("expected within-tolerance pixel difference to match")
	}
	if compareImages(base, outsideTolerance, 4) {
		t.Fatal("expected out-of-tolerance pixel difference to mismatch")
	}
}

func TestNormalizeRegionInputRejectsMissingTaggedState(t *testing.T) {
	_, err := normalizeRegionInput(RecognitionRegionInput{
		ID:     "region-01",
		Label:  "test-region",
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

func TestSaveRecognitionTemplateRejectsNegativeStateTolerance(t *testing.T) {
	app := NewApp()
	screenshot := base64.StdEncoding.EncodeToString(mustEncodePNGForTest(t, image.NewNRGBA(image.Rect(0, 0, 1, 1))))
	_, err := app.SaveRecognitionTemplate(RecognitionTemplateInput{
		Title:         "test-window",
		ClassName:     "test-class",
		ScreenshotPNG: screenshot,
		Width:         1,
		Height:        1,
		Regions: []RecognitionRegionInput{
			{
				ID:     "region-01",
				Label:  "region",
				X:      0,
				Y:      0,
				Width:  1,
				Height: 1,
				States: []RecognitionRegionStateInput{
					{
						ID:        "state-01",
						Tag:       "\u8fd1\u536b",
						Tolerance: -1,
						ImagePNG:  screenshot,
					},
				},
			},
		},
	})
	if err == nil {
		t.Fatal("expected negative state tolerance to be rejected")
	}
}

func TestSaveRecognitionTemplatePersistsPerStateTolerance(t *testing.T) {
	originalWorkingDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}

	tempDir := t.TempDir()
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("failed to change working directory: %v", err)
	}
	t.Cleanup(func() {
		if chdirErr := os.Chdir(originalWorkingDir); chdirErr != nil {
			t.Fatalf("failed to restore working directory: %v", chdirErr)
		}
	})

	app := NewApp()
	screenshot := base64.StdEncoding.EncodeToString(mustEncodePNGForTest(t, image.NewNRGBA(image.Rect(0, 0, 1, 1))))
	saved, err := app.SaveRecognitionTemplate(RecognitionTemplateInput{
		Title:         "test-window",
		ClassName:     "test-class",
		ScreenshotPNG: screenshot,
		Width:         1,
		Height:        1,
		Regions: []RecognitionRegionInput{
			{
				ID:     "region-01",
				Label:  "region",
				X:      0,
				Y:      0,
				Width:  1,
				Height: 1,
				States: []RecognitionRegionStateInput{
					{
						ID:        "state-01",
						Tag:       "\u8fd1\u536b",
						Tolerance: 3,
						ImagePNG:  screenshot,
					},
					{
						ID:        "state-02",
						Tag:       "\u72d9\u51fb",
						Tolerance: 0,
						ImagePNG:  screenshot,
					},
				},
			},
		},
	})
	if err != nil {
		t.Fatalf("SaveRecognitionTemplate returned error: %v", err)
	}
	if len(saved.Regions) != 1 || len(saved.Regions[0].States) != 2 {
		t.Fatalf("expected saved template to include two states, got %+v", saved.Regions)
	}
	if saved.Regions[0].States[0].Tolerance != 3 || saved.Regions[0].States[1].Tolerance != 0 {
		t.Fatalf("expected saved state tolerances to persist, got %+v", saved.Regions[0].States)
	}

	loaded, err := app.GetRecognitionTemplate(saved.ID)
	if err != nil {
		t.Fatalf("GetRecognitionTemplate returned error: %v", err)
	}
	if loaded.Regions[0].States[0].Tolerance != 3 || loaded.Regions[0].States[1].Tolerance != 0 {
		t.Fatalf("expected loaded state tolerances to persist, got %+v", loaded.Regions[0].States)
	}
}

func TestAggregateRecognizedRecruitmentTagsReturnsUniqueTags(t *testing.T) {
	tags, ok, message := aggregateRecognizedRecruitmentTags([]RecognitionRegionMatchResult{
		{
			RegionID: "region-01",
			MatchedStates: []RecognitionRegionStateMatchItem{
				{StateID: "state-01", Tag: "\u8fd1\u536b"},
			},
		},
		{
			RegionID: "region-02",
			MatchedStates: []RecognitionRegionStateMatchItem{
				{StateID: "state-02", Tag: "\u8f93\u51fa"},
			},
		},
		{
			RegionID: "region-03",
			MatchedStates: []RecognitionRegionStateMatchItem{
				{StateID: "state-03", Tag: "\u8fd1\u536b"},
			},
		},
	})
	if !ok {
		t.Fatalf("expected aggregateRecognizedRecruitmentTags to succeed, got %s", message)
	}
	if len(tags) != 2 || tags[0] != "\u8fd1\u536b" || tags[1] != "\u8f93\u51fa" {
		t.Fatalf("unexpected aggregated tags: %#v", tags)
	}
}

func TestAggregateRecognizedRecruitmentTagsRejectsZeroMatchRegion(t *testing.T) {
	_, ok, message := aggregateRecognizedRecruitmentTags([]RecognitionRegionMatchResult{
		{RegionID: "region-01", MatchedStates: []RecognitionRegionStateMatchItem{}},
	})
	if ok {
		t.Fatal("expected zero-match region to fail aggregation")
	}
	if message == "" {
		t.Fatal("expected failure message for zero-match region")
	}
}

func TestAggregateRecognizedRecruitmentTagsRejectsMultiMatchRegion(t *testing.T) {
	_, ok, message := aggregateRecognizedRecruitmentTags([]RecognitionRegionMatchResult{
		{
			RegionID: "region-01",
			MatchedStates: []RecognitionRegionStateMatchItem{
				{StateID: "state-01", Tag: "\u8fd1\u536b"},
				{StateID: "state-02", Tag: "\u72d9\u51fb"},
			},
		},
	})
	if ok {
		t.Fatal("expected multi-match region to fail aggregation")
	}
	if message == "" {
		t.Fatal("expected failure message for multi-match region")
	}
}

func TestRunPublicRecruitmentRecognitionReturnsNoTemplateFailure(t *testing.T) {
	app := NewApp()
	result, err := app.RunPublicRecruitmentRecognition(PublicRecruitmentRecognitionRequest{})
	if err != nil {
		t.Fatalf("RunPublicRecruitmentRecognition returned error: %v", err)
	}
	if result.Success {
		t.Fatal("expected empty template id to fail")
	}
	if result.FailureReason != publicRecruitmentRecognitionFailureNoTemplate {
		t.Fatalf("expected no_template failure, got %s", result.FailureReason)
	}
}

func TestRunPublicRecruitmentRecognitionReturnsWindowResolutionFailures(t *testing.T) {
	for _, tc := range []struct {
		name          string
		failureReason string
	}{
		{name: "no window", failureReason: publicRecruitmentRecognitionFailureNoWindow},
		{name: "ambiguous window", failureReason: publicRecruitmentRecognitionFailureAmbiguousWindow},
	} {
		t.Run(tc.name, func(t *testing.T) {
			templateID := prepareRecognitionTemplateFixture(t, recognitionTemplateFixture{
				regions: []recognitionTemplateFixtureRegion{
					{
						id:    "region-01",
						label: "slot-1",
						taggedStates: []recognitionTemplateFixtureState{
							{id: "state-01", tag: "\u8fd1\u536b", fill: color.NRGBA{R: 255, A: 255}},
						},
					},
				},
			})

			originalResolve := resolveWindowInstanceFunc
			defer func() {
				resolveWindowInstanceFunc = originalResolve
			}()
			resolveWindowInstanceFunc = func(template RecognitionTemplate) (uintptr, string, error) {
				return 0, tc.failureReason, os.ErrNotExist
			}

			app := NewApp()
			result, err := app.RunPublicRecruitmentRecognition(PublicRecruitmentRecognitionRequest{TemplateID: templateID})
			if err != nil {
				t.Fatalf("RunPublicRecruitmentRecognition returned error: %v", err)
			}
			if result.Success {
				t.Fatal("expected recognition to fail when window resolution fails")
			}
			if result.FailureReason != tc.failureReason {
				t.Fatalf("expected failure reason %s, got %s", tc.failureReason, result.FailureReason)
			}
		})
	}
}

func TestRunPublicRecruitmentRecognitionReturnsRecognizedTagsOnUniqueMatch(t *testing.T) {
	templateID := prepareRecognitionTemplateFixture(t, recognitionTemplateFixture{
		regions: []recognitionTemplateFixtureRegion{
			{
				id:    "region-01",
				label: "slot-1",
				taggedStates: []recognitionTemplateFixtureState{
					{id: "state-01", tag: "\u8fd1\u536b", fill: color.NRGBA{R: 255, A: 255}},
				},
			},
		},
	})

	originalResolve := resolveWindowInstanceFunc
	originalCapture := captureWindowPNGFunc
	defer func() {
		resolveWindowInstanceFunc = originalResolve
		captureWindowPNGFunc = originalCapture
	}()

	resolveWindowInstanceFunc = func(template RecognitionTemplate) (uintptr, string, error) {
		return 4321, publicRecruitmentRecognitionFailureNone, nil
	}
	captureWindowPNGFunc = func(hwnd uintptr) ([]byte, int, int, error) {
		img := image.NewNRGBA(image.Rect(0, 0, 2, 2))
		fillTestRect(img, img.Bounds(), color.NRGBA{R: 255, A: 255})
		return mustEncodePNGForTest(t, img), 2, 2, nil
	}

	app := NewApp()
	result, err := app.RunPublicRecruitmentRecognition(PublicRecruitmentRecognitionRequest{TemplateID: templateID})
	if err != nil {
		t.Fatalf("RunPublicRecruitmentRecognition returned error: %v", err)
	}
	if !result.Success {
		t.Fatalf("expected recognition success, got %+v", result)
	}
	if len(result.RecognizedTags) != 1 || result.RecognizedTags[0] != "\u8fd1\u536b" {
		t.Fatalf("unexpected recognized tags: %#v", result.RecognizedTags)
	}
}

func TestRunPublicRecruitmentRecognitionRejectsZeroAndMultiMatchRuns(t *testing.T) {
	for _, tc := range []struct {
		name           string
		template       recognitionTemplateFixture
		screenshotFill color.NRGBA
	}{
		{
			name: "zero match",
			template: recognitionTemplateFixture{
				regions: []recognitionTemplateFixtureRegion{
					{
						id:    "region-01",
						label: "slot-1",
						taggedStates: []recognitionTemplateFixtureState{
							{id: "state-01", tag: "\u8fd1\u536b", fill: color.NRGBA{R: 255, A: 255}},
						},
					},
				},
			},
			screenshotFill: color.NRGBA{B: 255, A: 255},
		},
		{
			name: "multi match",
			template: recognitionTemplateFixture{
				regions: []recognitionTemplateFixtureRegion{
					{
						id:    "region-01",
						label: "slot-1",
						taggedStates: []recognitionTemplateFixtureState{
							{id: "state-01", tag: "\u8fd1\u536b", fill: color.NRGBA{R: 255, A: 255}},
							{id: "state-02", tag: "\u8f93\u51fa", fill: color.NRGBA{R: 255, A: 255}},
						},
					},
				},
			},
			screenshotFill: color.NRGBA{R: 255, A: 255},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			templateID := prepareRecognitionTemplateFixture(t, tc.template)

			originalResolve := resolveWindowInstanceFunc
			originalCapture := captureWindowPNGFunc
			defer func() {
				resolveWindowInstanceFunc = originalResolve
				captureWindowPNGFunc = originalCapture
			}()

			resolveWindowInstanceFunc = func(template RecognitionTemplate) (uintptr, string, error) {
				return 8765, publicRecruitmentRecognitionFailureNone, nil
			}
			captureWindowPNGFunc = func(hwnd uintptr) ([]byte, int, int, error) {
				img := image.NewNRGBA(image.Rect(0, 0, 2, 2))
				fillTestRect(img, img.Bounds(), tc.screenshotFill)
				return mustEncodePNGForTest(t, img), 2, 2, nil
			}

			app := NewApp()
			result, err := app.RunPublicRecruitmentRecognition(PublicRecruitmentRecognitionRequest{TemplateID: templateID})
			if err != nil {
				t.Fatalf("RunPublicRecruitmentRecognition returned error: %v", err)
			}
			if result.Success {
				t.Fatalf("expected recognition failure for %s, got %+v", tc.name, result)
			}
			if result.FailureReason != publicRecruitmentRecognitionFailureIncompleteMatch {
				t.Fatalf("expected incomplete_match failure, got %s", result.FailureReason)
			}
			if len(result.RecognizedTags) != 0 {
				t.Fatalf("expected no recognized tags, got %#v", result.RecognizedTags)
			}
		})
	}
}

type recognitionTemplateFixture struct {
	regions []recognitionTemplateFixtureRegion
}

type recognitionTemplateFixtureRegion struct {
	id           string
	label        string
	taggedStates []recognitionTemplateFixtureState
}

type recognitionTemplateFixtureState struct {
	id        string
	tag       string
	tolerance int
	fill      color.NRGBA
}

func prepareRecognitionTemplateFixture(t *testing.T, fixture recognitionTemplateFixture) string {
	t.Helper()

	originalWorkingDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}

	tempDir := t.TempDir()
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("failed to change working directory: %v", err)
	}
	t.Cleanup(func() {
		if chdirErr := os.Chdir(originalWorkingDir); chdirErr != nil {
			t.Fatalf("failed to restore working directory: %v", chdirErr)
		}
	})

	templateID := "fixture-template"
	templateDir := filepath.Join(tempDir, recognitionTemplateDirName, templateID)
	if err := os.MkdirAll(filepath.Join(templateDir, "regions"), 0o755); err != nil {
		t.Fatalf("failed to create template dir: %v", err)
	}

	windowImage := image.NewNRGBA(image.Rect(0, 0, 2, 2))
	fillTestRect(windowImage, windowImage.Bounds(), color.NRGBA{R: 10, G: 20, B: 30, A: 255})
	if err := os.WriteFile(filepath.Join(templateDir, "window.png"), mustEncodePNGForTest(t, windowImage), 0o644); err != nil {
		t.Fatalf("failed to write window image: %v", err)
	}

	regions := make([]RecognitionRegion, 0, len(fixture.regions))
	for _, region := range fixture.regions {
		regionDir := filepath.Join(templateDir, "regions", region.id)
		if err := os.MkdirAll(regionDir, 0o755); err != nil {
			t.Fatalf("failed to create region dir: %v", err)
		}

		states := make([]RecognitionRegionState, 0, len(region.taggedStates))
		for _, state := range region.taggedStates {
			stateImage := image.NewNRGBA(image.Rect(0, 0, 2, 2))
			fillTestRect(stateImage, stateImage.Bounds(), state.fill)
			stateFilename := state.id + ".png"
			statePath := filepath.Join(regionDir, stateFilename)
			if err := os.WriteFile(statePath, mustEncodePNGForTest(t, stateImage), 0o644); err != nil {
				t.Fatalf("failed to write state image: %v", err)
			}
			states = append(states, RecognitionRegionState{
				ID:            state.id,
				Tag:           state.tag,
				Tolerance:     state.tolerance,
				ReferencePath: filepath.ToSlash(filepath.Join("regions", region.id, stateFilename)),
				CreatedAt:     time.Now().Format(time.RFC3339),
			})
		}

		regions = append(regions, RecognitionRegion{
			ID:     region.id,
			Label:  region.label,
			X:      0,
			Y:      0,
			Width:  1,
			Height: 1,
			States: states,
		})
	}

	template := RecognitionTemplate{
		ID:             templateID,
		Hwnd:           1234,
		Title:          "fixture-window",
		ClassName:      "fixture-class",
		Width:          2,
		Height:         2,
		ScreenshotPath: "window.png",
		CreatedAt:      time.Now().Format(time.RFC3339),
		Regions:        regions,
	}
	if err := saveRecognitionTemplateMetadata(templateDir, template); err != nil {
		t.Fatalf("failed to write template metadata: %v", err)
	}

	return templateID
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
