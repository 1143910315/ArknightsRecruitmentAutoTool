package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
	"unsafe"

	"github.com/lxn/win"
)

const recognitionTemplateDirName = "recognition-templates"

const (
	publicRecruitmentRecognitionFailureNone            = ""
	publicRecruitmentRecognitionFailureNoTemplate      = "no_template"
	publicRecruitmentRecognitionFailureNoWindow        = "no_window"
	publicRecruitmentRecognitionFailureAmbiguousWindow = "ambiguous_window"
	publicRecruitmentRecognitionFailureCaptureFailed   = "capture_failed"
	publicRecruitmentRecognitionFailureIncompleteMatch = "incomplete_match"
)

var (
	captureWindowPNGFunc                   = captureWindowPNG
	resolveRecognitionWindowCandidatesFunc = resolveRecognitionWindowCandidates
	isWindowHandleAliveFunc                = isWindowHandleAlive
	resolveWindowInstanceFunc              = resolveRecognitionWindowInstance
	isWindowProc                           = user32.NewProc("IsWindow")
)

type RecognitionWindowCaptureResult struct {
	Hwnd        uintptr `json:"hwnd"`
	Title       string  `json:"title"`
	ClassName   string  `json:"className"`
	Width       int     `json:"width"`
	Height      int     `json:"height"`
	ImageBase64 string  `json:"imageBase64"`
}

type RecognitionRegionInput struct {
	ID     string                        `json:"id"`
	Label  string                        `json:"label"`
	X      float64                       `json:"x"`
	Y      float64                       `json:"y"`
	Width  float64                       `json:"width"`
	Height float64                       `json:"height"`
	States []RecognitionRegionStateInput `json:"states"`
}

type RecognitionRegionStateInput struct {
	ID        string `json:"id"`
	Tag       string `json:"tag"`
	Tolerance int    `json:"tolerance"`
	ImagePNG  string `json:"imagePng"`
}

type RecognitionTemplateInput struct {
	ID            string                   `json:"id"`
	Hwnd          uintptr                  `json:"hwnd"`
	Title         string                   `json:"title"`
	ClassName     string                   `json:"className"`
	ScreenshotPNG string                   `json:"screenshotPng"`
	Width         int                      `json:"width"`
	Height        int                      `json:"height"`
	Regions       []RecognitionRegionInput `json:"regions"`
}

type RecognitionRegion struct {
	ID     string                   `json:"id"`
	Label  string                   `json:"label"`
	X      float64                  `json:"x"`
	Y      float64                  `json:"y"`
	Width  float64                  `json:"width"`
	Height float64                  `json:"height"`
	States []RecognitionRegionState `json:"states"`
}

type RecognitionRegionState struct {
	ID            string `json:"id"`
	Tag           string `json:"tag"`
	Tolerance     int    `json:"tolerance"`
	ReferencePath string `json:"referencePath"`
	ReferencePNG  string `json:"referencePng,omitempty"`
	CreatedAt     string `json:"createdAt,omitempty"`
}

type RecognitionTemplateSummary struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	ClassName   string `json:"className"`
	CreatedAt   string `json:"createdAt"`
	RegionCount int    `json:"regionCount"`
}

type RecognitionTemplate struct {
	ID             string                            `json:"id"`
	Hwnd           uintptr                           `json:"hwnd"`
	Title          string                            `json:"title"`
	ClassName      string                            `json:"className"`
	Instance       RecognitionWindowInstanceMetadata `json:"instance"`
	Width          int                               `json:"width"`
	Height         int                               `json:"height"`
	ScreenshotPath string                            `json:"screenshotPath"`
	ScreenshotPNG  string                            `json:"screenshotPng,omitempty"`
	CreatedAt      string                            `json:"createdAt"`
	Regions        []RecognitionRegion               `json:"regions"`
}

type RecognitionWindowInstanceMetadata struct {
	ProcessID uint32                  `json:"processId"`
	Bounds    RecognitionWindowBounds `json:"bounds"`
}

type RecognitionWindowBounds struct {
	Left   int `json:"left"`
	Top    int `json:"top"`
	Right  int `json:"right"`
	Bottom int `json:"bottom"`
}

type RecognitionWindowCandidate struct {
	Hwnd      uintptr                 `json:"hwnd"`
	Title     string                  `json:"title"`
	ClassName string                  `json:"className"`
	ProcessID uint32                  `json:"processId"`
	Bounds    RecognitionWindowBounds `json:"bounds"`
}

type RecognitionRegionMatchResult struct {
	RegionID      string                            `json:"regionId"`
	Label         string                            `json:"label"`
	Match         bool                              `json:"match"`
	MatchedStates []RecognitionRegionStateMatchItem `json:"matchedStates"`
}

type RecognitionRegionStateMatchItem struct {
	StateID string `json:"stateId"`
	Tag     string `json:"tag"`
}

type RecognitionMatchRequest struct {
	TemplateID string  `json:"templateId"`
	Hwnd       uintptr `json:"hwnd"`
}

type RecognitionMatchResult struct {
	TemplateID string                         `json:"templateId"`
	Hwnd       uintptr                        `json:"hwnd"`
	Results    []RecognitionRegionMatchResult `json:"results"`
}

type PublicRecruitmentRecognitionRequest struct {
	TemplateID string `json:"templateId"`
}

type PublicRecruitmentRecognitionResult struct {
	TemplateID     string                         `json:"templateId"`
	Hwnd           uintptr                        `json:"hwnd,omitempty"`
	Success        bool                           `json:"success"`
	FailureReason  string                         `json:"failureReason,omitempty"`
	FailureMessage string                         `json:"failureMessage,omitempty"`
	RecognizedTags []string                       `json:"recognizedTags"`
	Results        []RecognitionRegionMatchResult `json:"results"`
}

func (a *App) CaptureWindowForRecognition(hwnd uintptr) (RecognitionWindowCaptureResult, error) {
	if hwnd == 0 {
		return RecognitionWindowCaptureResult{}, errors.New("invalid window handle")
	}

	title, className := getWindowTextAndClass(hwnd)
	imageBytes, width, height, err := captureWindowPNG(hwnd)
	if err != nil {
		return RecognitionWindowCaptureResult{}, err
	}

	return RecognitionWindowCaptureResult{
		Hwnd:        hwnd,
		Title:       title,
		ClassName:   className,
		Width:       width,
		Height:      height,
		ImageBase64: base64.StdEncoding.EncodeToString(imageBytes),
	}, nil
}

func (a *App) SaveRecognitionTemplate(input RecognitionTemplateInput) (RecognitionTemplate, error) {
	if strings.TrimSpace(input.Title) == "" && strings.TrimSpace(input.ClassName) == "" {
		return RecognitionTemplate{}, errors.New("window identity is required")
	}
	if len(input.Regions) == 0 {
		return RecognitionTemplate{}, errors.New("at least one region is required")
	}

	screenshotBytes, err := base64.StdEncoding.DecodeString(strings.TrimSpace(input.ScreenshotPNG))
	if err != nil {
		return RecognitionTemplate{}, fmt.Errorf("failed to decode screenshot: %w", err)
	}

	templateID := strings.TrimSpace(input.ID)
	if templateID == "" {
		templateID = buildRecognitionTemplateID(input.Title, input.ClassName)
	}

	rootDir, err := ensureRecognitionTemplateRoot("")
	if err != nil {
		return RecognitionTemplate{}, err
	}

	templateDir := filepath.Join(rootDir, templateID)
	regionDir := filepath.Join(templateDir, "regions")
	if err := os.MkdirAll(regionDir, 0o755); err != nil {
		return RecognitionTemplate{}, err
	}

	screenshotPath := filepath.Join(templateDir, "window.png")
	if err := os.WriteFile(screenshotPath, screenshotBytes, 0o644); err != nil {
		return RecognitionTemplate{}, err
	}

	regions := make([]RecognitionRegion, 0, len(input.Regions))
	for index, regionInput := range input.Regions {
		normalized, err := normalizeRegionInput(regionInput)
		if err != nil {
			return RecognitionTemplate{}, fmt.Errorf("invalid region %d: %w", index+1, err)
		}

		regionID := normalized.ID
		if regionID == "" {
			regionID = fmt.Sprintf("region-%02d", index+1)
		}
		regionStateDir := filepath.Join(regionDir, regionID)
		if err := os.MkdirAll(regionStateDir, 0o755); err != nil {
			return RecognitionTemplate{}, fmt.Errorf("failed to prepare region state directory: %w", err)
		}

		states, err := saveRecognitionRegionStates(regionStateDir, regionID, normalized)
		if err != nil {
			return RecognitionTemplate{}, fmt.Errorf("failed to save region %d states: %w", index+1, err)
		}

		regions = append(regions, RecognitionRegion{
			ID:     regionID,
			Label:  normalized.Label,
			X:      normalized.X,
			Y:      normalized.Y,
			Width:  normalized.Width,
			Height: normalized.Height,
			States: states,
		})
	}

	template := RecognitionTemplate{
		ID:             templateID,
		Hwnd:           input.Hwnd,
		Title:          input.Title,
		ClassName:      input.ClassName,
		Instance:       captureRecognitionWindowInstanceMetadata(input.Hwnd),
		Width:          input.Width,
		Height:         input.Height,
		ScreenshotPath: filepath.ToSlash("window.png"),
		CreatedAt:      time.Now().Format(time.RFC3339),
		Regions:        regions,
	}

	if err := saveRecognitionTemplateMetadata(templateDir, template); err != nil {
		return RecognitionTemplate{}, err
	}

	template.ScreenshotPNG = base64.StdEncoding.EncodeToString(screenshotBytes)
	return template, nil
}

func (a *App) LoadRecognitionTemplates() ([]RecognitionTemplateSummary, error) {
	rootDir, err := ensureRecognitionTemplateRoot("")
	if err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(rootDir)
	if err != nil {
		return nil, err
	}

	summaries := make([]RecognitionTemplateSummary, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		template, err := readRecognitionTemplate(filepath.Join(rootDir, entry.Name()))
		if err != nil {
			continue
		}
		summaries = append(summaries, RecognitionTemplateSummary{
			ID:          template.ID,
			Title:       template.Title,
			ClassName:   template.ClassName,
			CreatedAt:   template.CreatedAt,
			RegionCount: len(template.Regions),
		})
	}

	sort.SliceStable(summaries, func(i, j int) bool {
		return summaries[i].CreatedAt > summaries[j].CreatedAt
	})

	return summaries, nil
}

func (a *App) GetRecognitionTemplate(id string) (RecognitionTemplate, error) {
	rootDir, err := ensureRecognitionTemplateRoot("")
	if err != nil {
		return RecognitionTemplate{}, err
	}

	templateDir := filepath.Join(rootDir, sanitizeRecognitionID(id))
	template, err := readRecognitionTemplate(templateDir)
	if err != nil {
		return RecognitionTemplate{}, err
	}

	screenshotBytes, err := os.ReadFile(filepath.Join(templateDir, filepath.FromSlash(template.ScreenshotPath)))
	if err != nil {
		return RecognitionTemplate{}, err
	}
	template.ScreenshotPNG = base64.StdEncoding.EncodeToString(screenshotBytes)
	if err := attachRecognitionStateImages(&template, templateDir); err != nil {
		return RecognitionTemplate{}, err
	}
	return template, nil
}

func (a *App) MatchRecognitionTemplate(input RecognitionMatchRequest) (RecognitionMatchResult, error) {
	if strings.TrimSpace(input.TemplateID) == "" {
		return RecognitionMatchResult{}, errors.New("template id is required")
	}
	if input.Hwnd == 0 {
		return RecognitionMatchResult{}, errors.New("window handle is required")
	}

	template, err := a.GetRecognitionTemplate(input.TemplateID)
	if err != nil {
		return RecognitionMatchResult{}, err
	}

	screenshotBytes, _, _, err := captureWindowPNGFunc(input.Hwnd)
	if err != nil {
		return RecognitionMatchResult{}, err
	}
	screenshotImage, err := png.Decode(bytes.NewReader(screenshotBytes))
	if err != nil {
		return RecognitionMatchResult{}, fmt.Errorf("failed to decode current screenshot: %w", err)
	}

	results := make([]RecognitionRegionMatchResult, 0, len(template.Regions))
	rootDir, err := ensureRecognitionTemplateRoot("")
	if err != nil {
		return RecognitionMatchResult{}, err
	}
	templateDir := filepath.Join(rootDir, sanitizeRecognitionID(template.ID))

	for _, region := range template.Regions {
		currentRegion, err := cropNormalizedRegion(screenshotImage, region.X, region.Y, region.Width, region.Height)
		if err != nil {
			results = append(results, RecognitionRegionMatchResult{
				RegionID:      region.ID,
				Label:         region.Label,
				Match:         false,
				MatchedStates: []RecognitionRegionStateMatchItem{},
			})
			continue
		}

		matchedStates := matchRecognitionRegionStates(templateDir, currentRegion, region.States)
		results = append(results, RecognitionRegionMatchResult{
			RegionID:      region.ID,
			Label:         region.Label,
			Match:         len(matchedStates) > 0,
			MatchedStates: matchedStates,
		})
	}

	return RecognitionMatchResult{TemplateID: template.ID, Hwnd: input.Hwnd, Results: results}, nil
}

func (a *App) RunPublicRecruitmentRecognition(input PublicRecruitmentRecognitionRequest) (PublicRecruitmentRecognitionResult, error) {
	templateID := strings.TrimSpace(input.TemplateID)
	if templateID == "" {
		return PublicRecruitmentRecognitionResult{
			Success:        false,
			FailureReason:  publicRecruitmentRecognitionFailureNoTemplate,
			FailureMessage: "template id is required",
			RecognizedTags: []string{},
			Results:        []RecognitionRegionMatchResult{},
		}, nil
	}

	template, err := a.GetRecognitionTemplate(templateID)
	if err != nil {
		return PublicRecruitmentRecognitionResult{}, err
	}

	resolvedHwnd, failureReason, err := resolveWindowInstanceFunc(template)
	if err != nil {
		return PublicRecruitmentRecognitionResult{
			TemplateID:     template.ID,
			Success:        false,
			FailureReason:  failureReason,
			FailureMessage: err.Error(),
			RecognizedTags: []string{},
			Results:        []RecognitionRegionMatchResult{},
		}, nil
	}

	matchResult, err := a.MatchRecognitionTemplate(RecognitionMatchRequest{
		TemplateID: template.ID,
		Hwnd:       resolvedHwnd,
	})
	if err != nil {
		return PublicRecruitmentRecognitionResult{
			TemplateID:     template.ID,
			Hwnd:           resolvedHwnd,
			Success:        false,
			FailureReason:  publicRecruitmentRecognitionFailureCaptureFailed,
			FailureMessage: err.Error(),
			RecognizedTags: []string{},
			Results:        []RecognitionRegionMatchResult{},
		}, nil
	}

	recognizedTags, success, failureMessage := aggregateRecognizedRecruitmentTags(matchResult.Results)
	if !success {
		return PublicRecruitmentRecognitionResult{
			TemplateID:     template.ID,
			Hwnd:           resolvedHwnd,
			Success:        false,
			FailureReason:  publicRecruitmentRecognitionFailureIncompleteMatch,
			FailureMessage: failureMessage,
			RecognizedTags: []string{},
			Results:        matchResult.Results,
		}, nil
	}

	return PublicRecruitmentRecognitionResult{
		TemplateID:     template.ID,
		Hwnd:           resolvedHwnd,
		Success:        true,
		RecognizedTags: recognizedTags,
		Results:        matchResult.Results,
	}, nil
}

func ensureRecognitionTemplateRoot(baseDir string) (string, error) {
	root := baseDir
	if root == "" {
		runtimeDir, err := resolveAppRuntimeDir()
		if err != nil {
			return "", err
		}
		root = filepath.Join(runtimeDir, recognitionTemplateDirName)
	}
	if err := os.MkdirAll(root, 0o755); err != nil {
		return "", err
	}
	return root, nil
}

func buildRecognitionTemplateID(title string, className string) string {
	seed := fmt.Sprintf("%s|%s|%d", strings.TrimSpace(title), strings.TrimSpace(className), time.Now().UnixNano())
	sum := sha1.Sum([]byte(seed))
	return hex.EncodeToString(sum[:8])
}

func sanitizeRecognitionID(id string) string {
	cleaned := strings.TrimSpace(id)
	cleaned = strings.ReplaceAll(cleaned, "..", "")
	cleaned = strings.ReplaceAll(cleaned, "/", "")
	cleaned = strings.ReplaceAll(cleaned, "\\", "")
	if cleaned == "" {
		return "template"
	}
	return cleaned
}

func normalizeRegionInput(input RecognitionRegionInput) (RecognitionRegionInput, error) {
	if input.Width <= 0 || input.Height <= 0 {
		return RecognitionRegionInput{}, errors.New("region width and height must be positive")
	}
	if len(input.States) == 0 {
		return RecognitionRegionInput{}, errors.New("at least one region state is required")
	}

	states := make([]RecognitionRegionStateInput, 0, len(input.States))
	for index, state := range input.States {
		normalizedState, err := normalizeRegionStateInput(state)
		if err != nil {
			return RecognitionRegionInput{}, fmt.Errorf("invalid state %d: %w", index+1, err)
		}
		states = append(states, normalizedState)
	}

	normalized := RecognitionRegionInput{
		ID:     strings.TrimSpace(input.ID),
		Label:  strings.TrimSpace(input.Label),
		X:      clamp01(input.X),
		Y:      clamp01(input.Y),
		Width:  clamp01(input.Width),
		Height: clamp01(input.Height),
		States: states,
	}
	if normalized.X >= 1 || normalized.Y >= 1 {
		return RecognitionRegionInput{}, errors.New("region origin is out of bounds")
	}
	if normalized.X+normalized.Width > 1 {
		normalized.Width = 1 - normalized.X
	}
	if normalized.Y+normalized.Height > 1 {
		normalized.Height = 1 - normalized.Y
	}
	if normalized.Width <= 0 || normalized.Height <= 0 {
		return RecognitionRegionInput{}, errors.New("region area is empty after normalization")
	}
	return normalized, nil
}

func normalizeRegionStateInput(input RecognitionRegionStateInput) (RecognitionRegionStateInput, error) {
	normalized := RecognitionRegionStateInput{
		ID:        strings.TrimSpace(input.ID),
		Tag:       strings.TrimSpace(input.Tag),
		Tolerance: input.Tolerance,
		ImagePNG:  strings.TrimSpace(input.ImagePNG),
	}
	if normalized.Tag == "" {
		return RecognitionRegionStateInput{}, errors.New("state tag is required")
	}
	if !isRecruitmentTag(normalized.Tag) {
		return RecognitionRegionStateInput{}, fmt.Errorf("unsupported recruitment tag: %s", normalized.Tag)
	}
	if normalized.Tolerance < 0 {
		return RecognitionRegionStateInput{}, errors.New("state tolerance must be a valid non-negative number")
	}
	if normalized.ImagePNG == "" {
		return RecognitionRegionStateInput{}, errors.New("state image is required")
	}
	return normalized, nil
}

func clamp01(value float64) float64 {
	switch {
	case value < 0:
		return 0
	case value > 1:
		return 1
	default:
		return value
	}
}

func saveRecognitionTemplateMetadata(templateDir string, template RecognitionTemplate) error {
	payload, err := json.MarshalIndent(template, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(templateDir, "template.json"), payload, 0o644)
}

func readRecognitionTemplate(templateDir string) (RecognitionTemplate, error) {
	raw, err := os.ReadFile(filepath.Join(templateDir, "template.json"))
	if err != nil {
		return RecognitionTemplate{}, err
	}
	return decodeRecognitionTemplate(raw)
}

type legacyRecognitionTemplate struct {
	ID             string                            `json:"id"`
	Hwnd           uintptr                           `json:"hwnd"`
	Title          string                            `json:"title"`
	ClassName      string                            `json:"className"`
	Instance       RecognitionWindowInstanceMetadata `json:"instance"`
	Width          int                               `json:"width"`
	Height         int                               `json:"height"`
	ScreenshotPath string                            `json:"screenshotPath"`
	ScreenshotPNG  string                            `json:"screenshotPng,omitempty"`
	CreatedAt      string                            `json:"createdAt"`
	Regions        []legacyRecognitionRegion         `json:"regions"`
}

type legacyRecognitionRegion struct {
	ID            string  `json:"id"`
	Label         string  `json:"label"`
	X             float64 `json:"x"`
	Y             float64 `json:"y"`
	Width         float64 `json:"width"`
	Height        float64 `json:"height"`
	ReferencePath string  `json:"referencePath"`
}

func decodeRecognitionTemplate(raw []byte) (RecognitionTemplate, error) {
	var template RecognitionTemplate
	if err := json.Unmarshal(raw, &template); err == nil && recognitionTemplateHasStates(template) {
		return template, nil
	}

	var legacy legacyRecognitionTemplate
	if err := json.Unmarshal(raw, &legacy); err != nil {
		return RecognitionTemplate{}, err
	}
	return normalizeLegacyRecognitionTemplate(legacy), nil
}

func recognitionTemplateHasStates(template RecognitionTemplate) bool {
	if len(template.Regions) == 0 {
		return true
	}
	for _, region := range template.Regions {
		if len(region.States) > 0 {
			return true
		}
	}
	return false
}

func normalizeLegacyRecognitionTemplate(legacy legacyRecognitionTemplate) RecognitionTemplate {
	regions := make([]RecognitionRegion, 0, len(legacy.Regions))
	for _, region := range legacy.Regions {
		tag := ""
		if isRecruitmentTag(strings.TrimSpace(region.Label)) {
			tag = strings.TrimSpace(region.Label)
		}
		states := []RecognitionRegionState{}
		if strings.TrimSpace(region.ReferencePath) != "" {
			states = append(states, RecognitionRegionState{
				ID:            defaultRegionStateID(region.ID, 1),
				Tag:           tag,
				Tolerance:     0,
				ReferencePath: filepath.ToSlash(region.ReferencePath),
			})
		}
		regions = append(regions, RecognitionRegion{
			ID:     region.ID,
			Label:  region.Label,
			X:      region.X,
			Y:      region.Y,
			Width:  region.Width,
			Height: region.Height,
			States: states,
		})
	}
	return RecognitionTemplate{
		ID:             legacy.ID,
		Hwnd:           legacy.Hwnd,
		Title:          legacy.Title,
		ClassName:      legacy.ClassName,
		Instance:       legacy.Instance,
		Width:          legacy.Width,
		Height:         legacy.Height,
		ScreenshotPath: legacy.ScreenshotPath,
		ScreenshotPNG:  legacy.ScreenshotPNG,
		CreatedAt:      legacy.CreatedAt,
		Regions:        regions,
	}
}

func captureRecognitionWindowInstanceMetadata(hwnd uintptr) RecognitionWindowInstanceMetadata {
	if hwnd == 0 {
		return RecognitionWindowInstanceMetadata{}
	}
	candidate := buildRecognitionWindowCandidate(hwnd)
	return RecognitionWindowInstanceMetadata{
		ProcessID: candidate.ProcessID,
		Bounds:    candidate.Bounds,
	}
}

func resolveRecognitionWindowInstance(template RecognitionTemplate) (uintptr, string, error) {
	if template.Hwnd != 0 && isWindowHandleAliveFunc(template.Hwnd) {
		return template.Hwnd, publicRecruitmentRecognitionFailureNone, nil
	}

	candidates, err := resolveRecognitionWindowCandidatesFunc(template.Title, template.ClassName)
	if err != nil {
		return 0, publicRecruitmentRecognitionFailureNoWindow, err
	}
	if len(candidates) == 0 {
		return 0, publicRecruitmentRecognitionFailureNoWindow, errors.New("no matching window instance found")
	}

	matchingCandidates := filterRecognitionWindowCandidates(template.Instance, candidates)
	if len(matchingCandidates) == 1 {
		return matchingCandidates[0].Hwnd, publicRecruitmentRecognitionFailureNone, nil
	}
	if len(matchingCandidates) == 0 {
		return 0, publicRecruitmentRecognitionFailureNoWindow, errors.New("no matching window instance found")
	}
	return 0, publicRecruitmentRecognitionFailureAmbiguousWindow, errors.New("multiple matching window instances found")
}

func resolveRecognitionWindowCandidates(title string, className string) ([]RecognitionWindowCandidate, error) {
	windows, err := (&App{}).GetTopWindows()
	if err != nil {
		return nil, err
	}

	filtered := make([]RecognitionWindowCandidate, 0)
	normalizedTitle := strings.TrimSpace(title)
	normalizedClassName := strings.TrimSpace(className)
	for _, windowInfo := range windows {
		if normalizedTitle != "" && windowInfo.Title != normalizedTitle {
			continue
		}
		if normalizedClassName != "" && windowInfo.ClassName != normalizedClassName {
			continue
		}
		if !isWindowHandleAliveFunc(windowInfo.Hwnd) {
			continue
		}
		filtered = append(filtered, buildRecognitionWindowCandidate(windowInfo.Hwnd))
	}
	return filtered, nil
}

func buildRecognitionWindowCandidate(hwnd uintptr) RecognitionWindowCandidate {
	title, className := getWindowTextAndClass(hwnd)
	return RecognitionWindowCandidate{
		Hwnd:      hwnd,
		Title:     title,
		ClassName: className,
		ProcessID: getWindowProcessID(hwnd),
		Bounds:    getRecognitionWindowBounds(hwnd),
	}
}

func filterRecognitionWindowCandidates(instance RecognitionWindowInstanceMetadata, candidates []RecognitionWindowCandidate) []RecognitionWindowCandidate {
	if len(candidates) == 0 {
		return nil
	}

	hasProcessID := instance.ProcessID != 0
	hasBounds := instance.Bounds != (RecognitionWindowBounds{})
	filtered := append([]RecognitionWindowCandidate(nil), candidates...)

	if hasProcessID {
		next := make([]RecognitionWindowCandidate, 0, len(filtered))
		for _, candidate := range filtered {
			if candidate.ProcessID == instance.ProcessID {
				next = append(next, candidate)
			}
		}
		if len(next) > 0 {
			filtered = next
		}
	}

	if hasBounds {
		next := make([]RecognitionWindowCandidate, 0, len(filtered))
		for _, candidate := range filtered {
			if candidate.Bounds == instance.Bounds {
				next = append(next, candidate)
			}
		}
		if len(next) > 0 {
			filtered = next
		}
	}

	return filtered
}

func aggregateRecognizedRecruitmentTags(results []RecognitionRegionMatchResult) ([]string, bool, string) {
	if len(results) == 0 {
		return nil, false, "template has no recognition regions"
	}

	seen := make(map[string]struct{}, len(results))
	tags := make([]string, 0, len(results))
	for _, result := range results {
		if len(result.MatchedStates) != 1 {
			if len(result.MatchedStates) == 0 {
				return nil, false, fmt.Sprintf("region %s matched no state", result.RegionID)
			}
			return nil, false, fmt.Sprintf("region %s matched multiple states", result.RegionID)
		}

		tag := strings.TrimSpace(result.MatchedStates[0].Tag)
		if !isRecruitmentTag(tag) {
			return nil, false, fmt.Sprintf("region %s matched unsupported recruitment tag", result.RegionID)
		}
		if _, ok := seen[tag]; ok {
			continue
		}
		seen[tag] = struct{}{}
		tags = append(tags, tag)
	}

	return tags, true, ""
}

func isWindowHandleAlive(hwnd uintptr) bool {
	if hwnd == 0 {
		return false
	}
	ret, _, _ := isWindowProc.Call(hwnd)
	return ret != 0
}

func getWindowProcessID(hwnd uintptr) uint32 {
	if hwnd == 0 {
		return 0
	}
	var processID uint32
	win.GetWindowThreadProcessId(win.HWND(hwnd), &processID)
	return processID
}

func getRecognitionWindowBounds(hwnd uintptr) RecognitionWindowBounds {
	rect, err := getWindowRect(win.HWND(hwnd))
	if err != nil {
		return RecognitionWindowBounds{}
	}
	return RecognitionWindowBounds{
		Left:   rect.Min.X,
		Top:    rect.Min.Y,
		Right:  rect.Max.X,
		Bottom: rect.Max.Y,
	}
}

func saveRecognitionRegionStates(regionStateDir string, regionID string, region RecognitionRegionInput) ([]RecognitionRegionState, error) {
	states := make([]RecognitionRegionState, 0, len(region.States))
	for index, state := range region.States {
		stateID := state.ID
		if stateID == "" {
			stateID = defaultRegionStateID(region.ID, index+1)
		}

		stateBytes, err := base64.StdEncoding.DecodeString(state.ImagePNG)
		if err != nil {
			return nil, fmt.Errorf("failed to decode state image: %w", err)
		}
		stateImage, err := png.Decode(bytes.NewReader(stateBytes))
		if err != nil {
			return nil, fmt.Errorf("failed to decode state png: %w", err)
		}

		statePath := filepath.Join(regionStateDir, stateID+".png")
		if err := writePNG(statePath, stateImage); err != nil {
			return nil, fmt.Errorf("failed to write state image: %w", err)
		}

		states = append(states, RecognitionRegionState{
			ID:            stateID,
			Tag:           state.Tag,
			Tolerance:     state.Tolerance,
			ReferencePath: filepath.ToSlash(filepath.Join("regions", regionID, stateID+".png")),
			CreatedAt:     time.Now().Format(time.RFC3339),
		})
	}
	return states, nil
}

func defaultRegionStateID(regionID string, index int) string {
	base := strings.TrimSpace(regionID)
	if base == "" {
		base = "region"
	}
	return fmt.Sprintf("%s-state-%02d", base, index)
}

func attachRecognitionStateImages(template *RecognitionTemplate, templateDir string) error {
	for regionIndex, region := range template.Regions {
		for stateIndex, state := range region.States {
			if strings.TrimSpace(state.ReferencePath) == "" {
				continue
			}
			stateBytes, err := os.ReadFile(filepath.Join(templateDir, filepath.FromSlash(state.ReferencePath)))
			if err != nil {
				return err
			}
			template.Regions[regionIndex].States[stateIndex].ReferencePNG = base64.StdEncoding.EncodeToString(stateBytes)
		}
	}
	return nil
}

func matchRecognitionRegionStates(templateDir string, currentRegion image.Image, states []RecognitionRegionState) []RecognitionRegionStateMatchItem {
	matchedStates := make([]RecognitionRegionStateMatchItem, 0)
	for _, state := range states {
		if strings.TrimSpace(state.ReferencePath) == "" {
			continue
		}
		referenceBytes, err := os.ReadFile(filepath.Join(templateDir, filepath.FromSlash(state.ReferencePath)))
		if err != nil {
			continue
		}
		referenceImage, err := png.Decode(bytes.NewReader(referenceBytes))
		if err != nil {
			continue
		}
		if compareImages(currentRegion, referenceImage, state.Tolerance) {
			matchedStates = append(matchedStates, RecognitionRegionStateMatchItem{
				StateID: state.ID,
				Tag:     state.Tag,
			})
		}
	}
	return matchedStates
}

func cropNormalizedRegion(img image.Image, x, y, width, height float64) (*image.NRGBA, error) {
	bounds := img.Bounds()
	imgWidth := bounds.Dx()
	imgHeight := bounds.Dy()
	if imgWidth <= 0 || imgHeight <= 0 {
		return nil, errors.New("image bounds are invalid")
	}

	minX := int(math.Round(x * float64(imgWidth)))
	minY := int(math.Round(y * float64(imgHeight)))
	maxX := int(math.Round((x + width) * float64(imgWidth)))
	maxY := int(math.Round((y + height) * float64(imgHeight)))

	if maxX <= minX {
		maxX = minX + 1
	}
	if maxY <= minY {
		maxY = minY + 1
	}
	if maxX > imgWidth {
		maxX = imgWidth
	}
	if maxY > imgHeight {
		maxY = imgHeight
	}

	rect := image.Rect(0, 0, maxX-minX, maxY-minY)
	dst := image.NewNRGBA(rect)
	draw.Draw(dst, rect, img, image.Point{X: bounds.Min.X + minX, Y: bounds.Min.Y + minY}, draw.Src)
	return dst, nil
}

func compareImages(left image.Image, right image.Image, tolerance int) bool {
	if left.Bounds().Dx() != right.Bounds().Dx() || left.Bounds().Dy() != right.Bounds().Dy() {
		return false
	}
	bounds := left.Bounds()
	for y := 0; y < bounds.Dy(); y++ {
		for x := 0; x < bounds.Dx(); x++ {
			if !sameColor(left.At(bounds.Min.X+x, bounds.Min.Y+y), right.At(right.Bounds().Min.X+x, right.Bounds().Min.Y+y), tolerance) {
				return false
			}
		}
	}
	return true
}

func sameColor(left color.Color, right color.Color, tolerance int) bool {
	if tolerance <= 0 {
		return color.NRGBAModel.Convert(left) == color.NRGBAModel.Convert(right)
	}

	leftColor := color.NRGBAModel.Convert(left).(color.NRGBA)
	rightColor := color.NRGBAModel.Convert(right).(color.NRGBA)

	return absInt(int(leftColor.R)-int(rightColor.R)) < tolerance &&
		absInt(int(leftColor.G)-int(rightColor.G)) < tolerance &&
		absInt(int(leftColor.B)-int(rightColor.B)) < tolerance &&
		absInt(int(leftColor.A)-int(rightColor.A)) < tolerance
}

func absInt(value int) int {
	if value < 0 {
		return -value
	}
	return value
}

func writePNG(targetPath string, img image.Image) error {
	imageBytes, err := encodePNG(img)
	if err != nil {
		return err
	}
	return os.WriteFile(targetPath, imageBytes, 0o644)
}

func captureWindowPNG(hwnd uintptr) ([]byte, int, int, error) {
	window := win.HWND(hwnd)

	windowRect, err := getWindowRect(window)
	if err != nil {
		return nil, 0, 0, err
	}
	if windowRect.Dx() <= 0 || windowRect.Dy() <= 0 {
		return nil, 0, 0, errors.New("window size is invalid")
	}

	screenSnapshot, err := captureScreenSnapshot()
	if err != nil {
		return nil, 0, 0, err
	}

	windowImage, err := cropImageRect(screenSnapshot, windowRect)
	if err != nil {
		return nil, 0, 0, err
	}

	imageBytes, err := encodePNG(windowImage)
	if err != nil {
		return nil, 0, 0, err
	}

	return imageBytes, windowImage.Bounds().Dx(), windowImage.Bounds().Dy(), nil
}

func getWindowRect(window win.HWND) (image.Rectangle, error) {
	var rect win.RECT
	if !win.GetWindowRect(window, &rect) {
		return image.Rectangle{}, errors.New("failed to get window bounds")
	}
	return image.Rect(int(rect.Left), int(rect.Top), int(rect.Right), int(rect.Bottom)), nil
}

func captureScreenSnapshot() (*image.NRGBA, error) {
	screenRect := getVirtualScreenRect()
	width := screenRect.Dx()
	height := screenRect.Dy()
	if width <= 0 || height <= 0 {
		return nil, errors.New("screen size is invalid")
	}

	screenDC := win.GetDC(0)
	if screenDC == 0 {
		return nil, errors.New("failed to get screen dc")
	}
	defer win.ReleaseDC(0, screenDC)

	memoryDC := win.CreateCompatibleDC(screenDC)
	if memoryDC == 0 {
		return nil, errors.New("failed to create memory dc")
	}
	defer win.DeleteDC(memoryDC)

	bitmap := win.CreateCompatibleBitmap(screenDC, int32(width), int32(height))
	if bitmap == 0 {
		return nil, errors.New("failed to create bitmap")
	}
	defer win.DeleteObject(win.HGDIOBJ(bitmap))

	oldObject := win.SelectObject(memoryDC, win.HGDIOBJ(bitmap))
	defer win.SelectObject(memoryDC, oldObject)

	if !win.BitBlt(memoryDC, 0, 0, int32(width), int32(height), screenDC, int32(screenRect.Min.X), int32(screenRect.Min.Y), win.SRCCOPY) {
		return nil, errors.New("failed to capture screen image")
	}

	pixelBytes, err := readBitmapPixels(memoryDC, bitmap, width, height)
	if err != nil {
		return nil, err
	}

	imageBuffer := image.NewNRGBA(screenRect)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			sourceOffset := (y*width + x) * 4
			destOffset := imageBuffer.PixOffset(screenRect.Min.X+x, screenRect.Min.Y+y)
			imageBuffer.Pix[destOffset+0] = pixelBytes[sourceOffset+2]
			imageBuffer.Pix[destOffset+1] = pixelBytes[sourceOffset+1]
			imageBuffer.Pix[destOffset+2] = pixelBytes[sourceOffset+0]
			imageBuffer.Pix[destOffset+3] = pixelBytes[sourceOffset+3]
		}
	}

	return imageBuffer, nil
}

func getVirtualScreenRect() image.Rectangle {
	left := int(win.GetSystemMetrics(win.SM_XVIRTUALSCREEN))
	top := int(win.GetSystemMetrics(win.SM_YVIRTUALSCREEN))
	width := int(win.GetSystemMetrics(win.SM_CXVIRTUALSCREEN))
	height := int(win.GetSystemMetrics(win.SM_CYVIRTUALSCREEN))
	return image.Rect(left, top, left+width, top+height)
}

func readBitmapPixels(memoryDC win.HDC, bitmap win.HBITMAP, width, height int) ([]byte, error) {
	if width <= 0 || height <= 0 {
		return nil, errors.New("bitmap size is invalid")
	}

	var bitmapInfo win.BITMAPINFO
	bitmapInfo.BmiHeader.BiSize = uint32(unsafe.Sizeof(bitmapInfo.BmiHeader))
	bitmapInfo.BmiHeader.BiWidth = int32(width)
	bitmapInfo.BmiHeader.BiHeight = -int32(height)
	bitmapInfo.BmiHeader.BiPlanes = 1
	bitmapInfo.BmiHeader.BiBitCount = 32
	bitmapInfo.BmiHeader.BiCompression = win.BI_RGB

	pixelBytes := make([]byte, width*height*4)
	if win.GetDIBits(memoryDC, bitmap, 0, uint32(height), &pixelBytes[0], &bitmapInfo, win.DIB_RGB_COLORS) == 0 {
		return nil, errors.New("failed to read bitmap pixels")
	}

	return pixelBytes, nil
}

func cropImageRect(src image.Image, target image.Rectangle) (*image.NRGBA, error) {
	visibleRect := src.Bounds().Intersect(target)
	if visibleRect.Empty() {
		return nil, errors.New("window has no visible intersection with screen snapshot")
	}

	dst := image.NewNRGBA(image.Rect(0, 0, visibleRect.Dx(), visibleRect.Dy()))
	draw.Draw(dst, dst.Bounds(), src, visibleRect.Min, draw.Src)
	return dst, nil
}

func encodePNG(img image.Image) ([]byte, error) {
	var buffer bytes.Buffer
	if err := png.Encode(&buffer, img); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}
