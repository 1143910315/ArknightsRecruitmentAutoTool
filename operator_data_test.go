package main

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestParseOperatorDataHTML(t *testing.T) {
	const sampleHTML = `
<div class="contentDetail" data-param1="狙击, 女, 6, , 高级资深干员, 拉特兰, 远程位, 输出, 中坚寻访, 公开招募, 否" data-param2="6">
  <p class="picture"><img src="https://example.com/images/exusiai.jpg" /><span class="picText">空弦</span></p>
  <p class="tags"><span class="btn btn-default tagText">远程位</span><span class="btn btn-default tagText">输出</span></p>
</div>
<div class="contentDetail" data-param1="近卫, 女, 5, 资深干员, , 哥伦比亚, 近战位, 输出, 防护, 公开招募, 中坚寻访, 是" data-param2="5">
  <p class="picture"><img src="https://example.com/images/indra.png" /><span class="picText">星极</span></p>
  <p class="tags"><span class="btn btn-default tagText">近战位</span><span class="btn btn-default tagText">输出</span><span class="btn btn-default tagText">防护</span></p>
</div>
<div class="contentDetail" data-param1="特种, 女, 4, , , 炎-龙门, 位移, 近战位, 关卡1-12首次通关掉落, 公开招募, 标准寻访, 中坚寻访, 主题曲获得, 否" data-param2="4">
  <p class="picture"><img src="https://example.com/images/shaw.webp" /><span class="picText">阿消</span></p>
  <p class="tags"><span class="btn btn-default tagText">位移</span><span class="btn btn-default tagText">近战位</span></p>
</div>`

	operators, err := parseOperatorDataHTML(strings.NewReader(sampleHTML))
	if err != nil {
		t.Fatalf("parseOperatorDataHTML returned error: %v", err)
	}

	if len(operators) != 3 {
		t.Fatalf("expected 3 operators, got %d", len(operators))
	}

	first := operators[0]
	if first.Order != 0 {
		t.Fatalf("expected first operator order to be 0, got %d", first.Order)
	}
	if first.Name != "空弦" {
		t.Fatalf("expected first operator to be 空弦, got %s", first.Name)
	}
	if first.Rarity != 6 {
		t.Fatalf("expected rarity 6, got %d", first.Rarity)
	}
	if first.RemoteImageURL != "https://example.com/images/exusiai.jpg" {
		t.Fatalf("unexpected remote image url: %s", first.RemoteImageURL)
	}
	if !first.IsPublicRecruitable {
		t.Fatal("expected first operator to be public recruitable")
	}
	if first.Metadata.Profession != "狙击" || first.Metadata.Gender != "女" {
		t.Fatalf("unexpected first metadata: %+v", first.Metadata)
	}
	if first.Metadata.Origin != "拉特兰" {
		t.Fatalf("expected origin 拉特兰, got %s", first.Metadata.Origin)
	}
	if !contains(first.DisplayTags, "远程位") || !contains(first.DisplayTags, "输出") {
		t.Fatalf("unexpected display tags: %#v", first.DisplayTags)
	}
	if !contains(first.Metadata.SeniorityTags, "高级资深干员") {
		t.Fatalf("expected 高级资深干员 seniority tag, got %#v", first.Metadata.SeniorityTags)
	}

	third := operators[2]
	if third.Order != 2 {
		t.Fatalf("expected third operator order to be 2, got %d", third.Order)
	}
	if !contains(third.Metadata.AcquisitionMethods, "关卡1-12首次通关掉落") {
		t.Fatalf("expected acquisition method to be preserved, got %#v", third.Metadata.AcquisitionMethods)
	}
	if !contains(third.Metadata.AcquisitionMethods, "主题曲获得") {
		t.Fatalf("expected theme unlock acquisition to be preserved, got %#v", third.Metadata.AcquisitionMethods)
	}
}

func TestEnsureOperatorCacheDirUsesProvidedRuntimeDirectory(t *testing.T) {
	runtimeDir := t.TempDir()
	explicitCacheDir, err := ensureOperatorCacheDir(filepath.Join(runtimeDir, "operator-data"))
	if err != nil {
		t.Fatalf("ensureOperatorCacheDir returned error: %v", err)
	}

	expected := filepath.Join(runtimeDir, "operator-data")
	if explicitCacheDir != expected {
		t.Fatalf("expected cache dir %s, got %s", expected, explicitCacheDir)
	}
}

func TestSaveAndLoadOperatorCachePreservesOrder(t *testing.T) {
	cacheDir := t.TempDir()
	cache := operatorCachePayload{
		SourceURL: operatorDataSourceURL,
		FetchedAt: time.Now().Format(time.RFC3339),
		Operators: []OperatorRecord{
			{Order: 2, Name: "阿消", LocalImagePath: filepath.Join("images", "002.jpg")},
			{Order: 0, Name: "空弦", LocalImagePath: filepath.Join("images", "000.jpg")},
		},
	}

	if err := saveOperatorCache(cacheDir, cache); err != nil {
		t.Fatalf("saveOperatorCache returned error: %v", err)
	}

	loaded, err := loadCachedOperatorDataFromDir(cacheDir)
	if err != nil {
		t.Fatalf("loadCachedOperatorDataFromDir returned error: %v", err)
	}
	if !loaded.CacheAvailable || !loaded.FromCache {
		t.Fatalf("expected cache result to indicate local cache usage: %+v", loaded)
	}
	if len(loaded.Operators) != 2 {
		t.Fatalf("expected 2 operators, got %d", len(loaded.Operators))
	}
	if loaded.Operators[0].Name != "空弦" || loaded.Operators[1].Name != "阿消" {
		t.Fatalf("expected order to be preserved after load, got %#v", loaded.Operators)
	}
	if loaded.Operators[0].LocalImagePath != "images/000.jpg" {
		t.Fatalf("expected local image path to be normalized, got %s", loaded.Operators[0].LocalImagePath)
	}
	if loaded.Operators[0].LocalImageURL != "" {
		t.Fatalf("expected local image url to stay empty, got %s", loaded.Operators[0].LocalImageURL)
	}
}

func TestResolveCachedOperatorImagePathUsesRelativeCachePath(t *testing.T) {
	cacheDir := t.TempDir()
	imagePath, err := resolveCachedOperatorImagePath(cacheDir, "images/007.png")
	if err != nil {
		t.Fatalf("resolveCachedOperatorImagePath returned error: %v", err)
	}

	expected := filepath.Join(cacheDir, "images", "007.png")
	if imagePath != expected {
		t.Fatalf("expected image path %s, got %s", expected, imagePath)
	}
}

func TestLoadCachedOperatorImageFromDirReturnsBase64AndMimeType(t *testing.T) {
	cacheDir, err := ensureOperatorCacheDir(t.TempDir())
	if err != nil {
		t.Fatalf("ensureOperatorCacheDir returned error: %v", err)
	}

	imageBytes := []byte("gif-image-data")
	imagePath := filepath.Join(cacheDir, "images", "000.gif")
	if err := os.WriteFile(imagePath, imageBytes, 0o644); err != nil {
		t.Fatalf("failed to write test image: %v", err)
	}

	result, err := loadCachedOperatorImageFromDir(cacheDir, "images/000.gif")
	if err != nil {
		t.Fatalf("loadCachedOperatorImageFromDir returned error: %v", err)
	}
	if !result.Found {
		t.Fatal("expected cached image to be found")
	}
	if result.MimeType != "image/gif" {
		t.Fatalf("expected mime type image/gif, got %s", result.MimeType)
	}
	if result.DataBase64 != base64.StdEncoding.EncodeToString(imageBytes) {
		t.Fatalf("unexpected base64 payload: %s", result.DataBase64)
	}
}

func TestLoadCachedOperatorImageFromDirReturnsNotFoundForMissingImage(t *testing.T) {
	cacheDir, err := ensureOperatorCacheDir(t.TempDir())
	if err != nil {
		t.Fatalf("ensureOperatorCacheDir returned error: %v", err)
	}

	result, err := loadCachedOperatorImageFromDir(cacheDir, "images/missing.jpg")
	if err != nil {
		t.Fatalf("expected missing image to avoid fatal error, got %v", err)
	}
	if result.Found {
		t.Fatalf("expected missing image result, got %+v", result)
	}
	if result.MimeType != "" || result.DataBase64 != "" {
		t.Fatalf("expected empty payload for missing image, got %+v", result)
	}
}

func TestLoadCachedOperatorImageFromDirRejectsEscapingCacheDirectory(t *testing.T) {
	cacheDir, err := ensureOperatorCacheDir(t.TempDir())
	if err != nil {
		t.Fatalf("ensureOperatorCacheDir returned error: %v", err)
	}

	_, err = loadCachedOperatorImageFromDir(cacheDir, "../outside.jpg")
	if err == nil {
		t.Fatal("expected escaping path to fail")
	}
}

func TestLoadCachedOperatorDataReturnsEmptyWhenCacheMissing(t *testing.T) {
	loaded, err := loadCachedOperatorDataFromDir(t.TempDir())
	if err != nil {
		t.Fatalf("expected missing cache to return nil error, got %v", err)
	}
	if loaded.CacheAvailable {
		t.Fatalf("expected no cache to be available, got %+v", loaded)
	}
}

func TestCacheOperatorImagesAllowsPartialFailures(t *testing.T) {
	imageServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/ok.jpg" {
			w.Header().Set("Content-Type", "image/jpeg")
			_, _ = w.Write([]byte("image-bytes"))
			return
		}
		http.NotFound(w, r)
	}))
	defer imageServer.Close()

	cacheDir, err := ensureOperatorCacheDir(t.TempDir())
	if err != nil {
		t.Fatalf("ensureOperatorCacheDir returned error: %v", err)
	}

	operators := []OperatorRecord{
		{Order: 0, Name: "空弦", RemoteImageURL: imageServer.URL + "/ok.jpg"},
		{Order: 1, Name: "阿消", RemoteImageURL: imageServer.URL + "/missing.jpg"},
	}

	cached := cacheOperatorImages(imageServer.Client(), cacheDir, operators)
	if cached[0].LocalImagePath != "images/000.jpg" {
		t.Fatalf("expected first operator image path to be relative, got %s", cached[0].LocalImagePath)
	}
	if _, err := os.Stat(filepath.Join(cacheDir, filepath.FromSlash(cached[0].LocalImagePath))); err != nil {
		t.Fatalf("expected cached image file to exist: %v", err)
	}
	if cached[0].LocalImageURL != "" {
		t.Fatalf("expected first operator image url to remain empty, got %s", cached[0].LocalImageURL)
	}
	if cached[1].LocalImagePath != "" {
		t.Fatalf("expected missing image to leave local image path empty, got %s", cached[1].LocalImagePath)
	}
}

func TestParseOperatorDataHTMLRejectsMalformedInput(t *testing.T) {
	_, err := parseOperatorDataHTML(strings.NewReader(`<div class="contentDetail" data-param1="狙击" data-param2="not-a-number"><span class="picText">空弦</span></div>`))
	if err == nil {
		t.Fatal("expected malformed rarity to fail")
	}
}

func TestParseOperatorDataHTMLRejectsMissingEntries(t *testing.T) {
	_, err := parseOperatorDataHTML(strings.NewReader(`<html><body><p>empty</p></body></html>`))
	if err == nil {
		t.Fatal("expected missing contentDetail entries to fail")
	}
}
