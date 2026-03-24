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

type RecognitionWindowCaptureResult struct {
	Hwnd        uintptr `json:"hwnd"`
	Title       string  `json:"title"`
	ClassName   string  `json:"className"`
	Width       int     `json:"width"`
	Height      int     `json:"height"`
	ImageBase64 string  `json:"imageBase64"`
}

type RecognitionRegionInput struct {
	ID     string  `json:"id"`
	Label  string  `json:"label"`
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
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
	ID            string  `json:"id"`
	Label         string  `json:"label"`
	X             float64 `json:"x"`
	Y             float64 `json:"y"`
	Width         float64 `json:"width"`
	Height        float64 `json:"height"`
	ReferencePath string  `json:"referencePath"`
}

type RecognitionTemplateSummary struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	ClassName   string `json:"className"`
	CreatedAt   string `json:"createdAt"`
	RegionCount int    `json:"regionCount"`
}

type RecognitionTemplate struct {
	ID             string              `json:"id"`
	Hwnd           uintptr             `json:"hwnd"`
	Title          string              `json:"title"`
	ClassName      string              `json:"className"`
	Width          int                 `json:"width"`
	Height         int                 `json:"height"`
	ScreenshotPath string              `json:"screenshotPath"`
	ScreenshotPNG  string              `json:"screenshotPng,omitempty"`
	CreatedAt      string              `json:"createdAt"`
	Regions        []RecognitionRegion `json:"regions"`
}

type RecognitionRegionMatchResult struct {
	RegionID string `json:"regionId"`
	Label    string `json:"label"`
	Match    bool   `json:"match"`
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

var printWindow = user32.NewProc("PrintWindow")

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

	screenshotImage, err := png.Decode(bytes.NewReader(screenshotBytes))
	if err != nil {
		return RecognitionTemplate{}, fmt.Errorf("failed to decode screenshot png: %w", err)
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

		regionImage, err := cropNormalizedRegion(screenshotImage, normalized.X, normalized.Y, normalized.Width, normalized.Height)
		if err != nil {
			return RecognitionTemplate{}, fmt.Errorf("failed to crop region %d: %w", index+1, err)
		}

		regionID := normalized.ID
		if regionID == "" {
			regionID = fmt.Sprintf("region-%02d", index+1)
		}
		referencePath := filepath.Join(regionDir, regionID+".png")
		if err := writePNG(referencePath, regionImage); err != nil {
			return RecognitionTemplate{}, fmt.Errorf("failed to save region image: %w", err)
		}

		regions = append(regions, RecognitionRegion{
			ID:            regionID,
			Label:         normalized.Label,
			X:             normalized.X,
			Y:             normalized.Y,
			Width:         normalized.Width,
			Height:        normalized.Height,
			ReferencePath: filepath.ToSlash(filepath.Join("regions", regionID+".png")),
		})
	}

	template := RecognitionTemplate{
		ID:             templateID,
		Hwnd:           input.Hwnd,
		Title:          input.Title,
		ClassName:      input.ClassName,
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

	screenshotBytes, _, _, err := captureWindowPNG(input.Hwnd)
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
			results = append(results, RecognitionRegionMatchResult{RegionID: region.ID, Label: region.Label, Match: false})
			continue
		}

		referenceBytes, err := os.ReadFile(filepath.Join(templateDir, filepath.FromSlash(region.ReferencePath)))
		if err != nil {
			results = append(results, RecognitionRegionMatchResult{RegionID: region.ID, Label: region.Label, Match: false})
			continue
		}
		referenceImage, err := png.Decode(bytes.NewReader(referenceBytes))
		if err != nil {
			results = append(results, RecognitionRegionMatchResult{RegionID: region.ID, Label: region.Label, Match: false})
			continue
		}

		results = append(results, RecognitionRegionMatchResult{
			RegionID: region.ID,
			Label:    region.Label,
			Match:    compareImages(currentRegion, referenceImage),
		})
	}

	return RecognitionMatchResult{TemplateID: template.ID, Hwnd: input.Hwnd, Results: results}, nil
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
	if strings.TrimSpace(input.Label) == "" {
		return RecognitionRegionInput{}, errors.New("label is required")
	}
	if input.Width <= 0 || input.Height <= 0 {
		return RecognitionRegionInput{}, errors.New("region width and height must be positive")
	}
	normalized := RecognitionRegionInput{
		ID:     strings.TrimSpace(input.ID),
		Label:  strings.TrimSpace(input.Label),
		X:      clamp01(input.X),
		Y:      clamp01(input.Y),
		Width:  clamp01(input.Width),
		Height: clamp01(input.Height),
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
	var template RecognitionTemplate
	if err := json.Unmarshal(raw, &template); err != nil {
		return RecognitionTemplate{}, err
	}
	return template, nil
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

func compareImages(left image.Image, right image.Image) bool {
	if left.Bounds().Dx() != right.Bounds().Dx() || left.Bounds().Dy() != right.Bounds().Dy() {
		return false
	}
	bounds := left.Bounds()
	for y := 0; y < bounds.Dy(); y++ {
		for x := 0; x < bounds.Dx(); x++ {
			if !sameColor(left.At(bounds.Min.X+x, bounds.Min.Y+y), right.At(right.Bounds().Min.X+x, right.Bounds().Min.Y+y)) {
				return false
			}
		}
	}
	return true
}

func sameColor(left color.Color, right color.Color) bool {
	lr, lg, lb, la := left.RGBA()
	rr, rg, rb, ra := right.RGBA()
	return lr == rr && lg == rg && lb == rb && la == ra
}

func writePNG(targetPath string, img image.Image) error {
	var buffer bytes.Buffer
	if err := png.Encode(&buffer, img); err != nil {
		return err
	}
	return os.WriteFile(targetPath, buffer.Bytes(), 0o644)
}

func captureWindowPNG(hwnd uintptr) ([]byte, int, int, error) {
	window := win.HWND(hwnd)

	var rect win.RECT
	if !win.GetWindowRect(window, &rect) {
		return nil, 0, 0, errors.New("failed to get window bounds")
	}

	width := int(rect.Right - rect.Left)
	height := int(rect.Bottom - rect.Top)
	if width <= 0 || height <= 0 {
		return nil, 0, 0, errors.New("window size is invalid")
	}

	windowDC := win.GetDC(window)
	if windowDC == 0 {
		return nil, 0, 0, errors.New("failed to get window dc")
	}
	defer win.ReleaseDC(window, windowDC)

	memoryDC := win.CreateCompatibleDC(windowDC)
	if memoryDC == 0 {
		return nil, 0, 0, errors.New("failed to create memory dc")
	}
	defer win.DeleteDC(memoryDC)

	bitmap := win.CreateCompatibleBitmap(windowDC, int32(width), int32(height))
	if bitmap == 0 {
		return nil, 0, 0, errors.New("failed to create bitmap")
	}
	defer win.DeleteObject(win.HGDIOBJ(bitmap))

	oldObject := win.SelectObject(memoryDC, win.HGDIOBJ(bitmap))
	defer win.SelectObject(memoryDC, oldObject)

	success, _, _ := printWindow.Call(uintptr(window), uintptr(memoryDC), uintptr(0))
	if success == 0 {
		if !win.BitBlt(memoryDC, 0, 0, int32(width), int32(height), windowDC, 0, 0, win.SRCCOPY) {
			return nil, 0, 0, errors.New("failed to capture window image")
		}
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
		return nil, 0, 0, errors.New("failed to read bitmap pixels")
	}

	imageBuffer := image.NewNRGBA(image.Rect(0, 0, width, height))
	for index := 0; index < len(pixelBytes); index += 4 {
		imageBuffer.Pix[index+0] = pixelBytes[index+2]
		imageBuffer.Pix[index+1] = pixelBytes[index+1]
		imageBuffer.Pix[index+2] = pixelBytes[index+0]
		imageBuffer.Pix[index+3] = pixelBytes[index+3]
	}

	var buffer bytes.Buffer
	if err := png.Encode(&buffer, imageBuffer); err != nil {
		return nil, 0, 0, err
	}

	return buffer.Bytes(), width, height, nil
}
