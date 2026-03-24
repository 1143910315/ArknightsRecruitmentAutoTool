package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html"
)

const (
	operatorDataSourceURL = "https://wiki.biligame.com/arknights/公开招募工具"
	operatorCacheFileName = "operators.json"
)

var errOperatorCacheNotFound = errors.New("operator cache not found")

type OperatorMetadata struct {
	Raw                []string `json:"raw"`
	Profession         string   `json:"profession"`
	Gender             string   `json:"gender"`
	RawRarity          string   `json:"rawRarity"`
	Origin             string   `json:"origin"`
	SeniorityTags      []string `json:"seniorityTags"`
	RecruitmentTags    []string `json:"recruitmentTags"`
	AcquisitionMethods []string `json:"acquisitionMethods"`
	Extra              []string `json:"extra"`
}

type OperatorRecord struct {
	Order               int              `json:"order"`
	Name                string           `json:"name"`
	Rarity              int              `json:"rarity"`
	DisplayTags         []string         `json:"displayTags"`
	RemoteImageURL      string           `json:"remoteImageUrl"`
	LocalImagePath      string           `json:"localImagePath"`
	LocalImageURL       string           `json:"localImageUrl"`
	IsPublicRecruitable bool             `json:"isPublicRecruitable"`
	Metadata            OperatorMetadata `json:"metadata"`
}

type FetchOperatorDataResult struct {
	SourceURL      string           `json:"sourceUrl"`
	FetchedAt      string           `json:"fetchedAt"`
	Operators      []OperatorRecord `json:"operators"`
	FromCache      bool             `json:"fromCache"`
	CacheAvailable bool             `json:"cacheAvailable"`
}

type operatorCachePayload struct {
	SourceURL string           `json:"sourceUrl"`
	FetchedAt string           `json:"fetchedAt"`
	Operators []OperatorRecord `json:"operators"`
}

func (a *App) FetchOperatorData() (FetchOperatorDataResult, error) {
	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Get(operatorDataSourceURL)
	if err != nil {
		return FetchOperatorDataResult{}, fmt.Errorf("failed to fetch operator data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return FetchOperatorDataResult{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 8<<20))
	if err != nil {
		return FetchOperatorDataResult{}, fmt.Errorf("failed to read operator data response: %w", err)
	}

	if len(strings.TrimSpace(string(body))) == 0 {
		return FetchOperatorDataResult{}, errors.New("operator data response is empty")
	}

	operators, err := parseOperatorDataHTML(strings.NewReader(string(body)))
	if err != nil {
		return FetchOperatorDataResult{}, err
	}

	cacheDir, err := ensureOperatorCacheDir("")
	if err != nil {
		return FetchOperatorDataResult{}, fmt.Errorf("failed to prepare operator cache directory: %w", err)
	}

	operators = cacheOperatorImages(client, cacheDir, operators)
	cache := operatorCachePayload{
		SourceURL: operatorDataSourceURL,
		FetchedAt: time.Now().Format(time.RFC3339),
		Operators: operators,
	}
	if err := saveOperatorCache(cacheDir, cache); err != nil {
		return FetchOperatorDataResult{}, fmt.Errorf("failed to save operator cache: %w", err)
	}

	return loadCachedOperatorDataFromDir(cacheDir)
}

func (a *App) LoadCachedOperatorData() (FetchOperatorDataResult, error) {
	return loadCachedOperatorDataFromDir("")
}

func loadCachedOperatorDataFromDir(baseDir string) (FetchOperatorDataResult, error) {
	cacheDir, err := ensureOperatorCacheDir(baseDir)
	if err != nil {
		return FetchOperatorDataResult{}, fmt.Errorf("failed to resolve operator cache directory: %w", err)
	}

	cache, err := readOperatorCache(cacheDir)
	if err != nil {
		if errors.Is(err, errOperatorCacheNotFound) {
			return FetchOperatorDataResult{CacheAvailable: false}, nil
		}
		return FetchOperatorDataResult{}, err
	}

	sortOperatorRecords(cache.Operators)
	applyLocalImageURLs(cacheDir, cache.Operators)

	return FetchOperatorDataResult{
		SourceURL:      cache.SourceURL,
		FetchedAt:      cache.FetchedAt,
		Operators:      cache.Operators,
		FromCache:      true,
		CacheAvailable: true,
	}, nil
}

func parseOperatorDataHTML(r io.Reader) ([]OperatorRecord, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, fmt.Errorf("failed to parse operator data HTML: %w", err)
	}

	nodes := findNodesByClass(doc, "contentDetail")
	if len(nodes) == 0 {
		return nil, errors.New("no operator entries found in source HTML")
	}

	operators := make([]OperatorRecord, 0, len(nodes))
	for index, node := range nodes {
		record, err := parseOperatorRecord(node, index)
		if err != nil {
			return nil, err
		}
		operators = append(operators, record)
	}

	return operators, nil
}

func parseOperatorRecord(node *html.Node, order int) (OperatorRecord, error) {
	nameNode := findFirstNodeByClass(node, "picText")
	name := strings.TrimSpace(extractText(nameNode))
	if name == "" {
		return OperatorRecord{}, errors.New("operator entry is missing .picText content")
	}

	rarityValue := strings.TrimSpace(getAttr(node, "data-param2"))
	if rarityValue == "" {
		return OperatorRecord{}, fmt.Errorf("operator %s is missing data-param2", name)
	}

	rarity, err := strconv.Atoi(rarityValue)
	if err != nil {
		return OperatorRecord{}, fmt.Errorf("operator %s has invalid rarity %q", name, rarityValue)
	}

	tagNodes := findNodesByClass(node, "tagText")
	displayTags := make([]string, 0, len(tagNodes))
	for _, tagNode := range tagNodes {
		tag := strings.TrimSpace(extractText(tagNode))
		if tag == "" {
			continue
		}
		displayTags = append(displayTags, tag)
	}

	metadata := parseOperatorMetadata(strings.TrimSpace(getAttr(node, "data-param1")), displayTags)

	return OperatorRecord{
		Order:               order,
		Name:                name,
		Rarity:              rarity,
		DisplayTags:         displayTags,
		RemoteImageURL:      extractRemoteImageURL(node),
		IsPublicRecruitable: contains(metadata.AcquisitionMethods, "公开招募"),
		Metadata:            metadata,
	}, nil
}

func parseOperatorMetadata(raw string, displayTags []string) OperatorMetadata {
	parts := splitMetadata(raw)
	metadata := OperatorMetadata{
		Raw:             parts,
		RecruitmentTags: append([]string(nil), displayTags...),
	}

	if len(parts) > 0 {
		metadata.Profession = parts[0]
	}
	if len(parts) > 1 {
		metadata.Gender = parts[1]
	}
	if len(parts) > 2 {
		metadata.RawRarity = parts[2]
	}

	tagSet := make(map[string]struct{}, len(displayTags))
	for _, tag := range displayTags {
		tagSet[tag] = struct{}{}
	}

	for index, part := range parts {
		if index < 3 {
			continue
		}

		switch {
		case strings.Contains(part, "资深干员"):
			metadata.SeniorityTags = append(metadata.SeniorityTags, part)
		case isAcquisitionField(part):
			metadata.AcquisitionMethods = append(metadata.AcquisitionMethods, part)
		case part == "是" || part == "否":
			metadata.Extra = append(metadata.Extra, part)
		case metadata.Origin == "" && !isDisplayedTag(part, tagSet):
			metadata.Origin = part
		case isDisplayedTag(part, tagSet):
			metadata.RecruitmentTags = appendIfMissing(metadata.RecruitmentTags, part)
		default:
			metadata.Extra = append(metadata.Extra, part)
		}
	}

	return metadata
}

func ensureOperatorCacheDir(baseDir string) (string, error) {
	root := baseDir
	if root == "" {
		userCacheDir, err := os.UserCacheDir()
		if err != nil {
			return "", err
		}
		root = filepath.Join(userCacheDir, "ArknightsRecruitmentAutoTool", "operator-data")
	}

	imageDir := filepath.Join(root, "images")
	if err := os.MkdirAll(imageDir, 0o755); err != nil {
		return "", err
	}
	return root, nil
}

func operatorCacheFilePath(baseDir string) string {
	return filepath.Join(baseDir, operatorCacheFileName)
}

func operatorImageDir(baseDir string) string {
	return filepath.Join(baseDir, "images")
}

func saveOperatorCache(baseDir string, cache operatorCachePayload) error {
	payload, err := json.MarshalIndent(cache, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(operatorCacheFilePath(baseDir), payload, 0o644)
}

func readOperatorCache(baseDir string) (operatorCachePayload, error) {
	raw, err := os.ReadFile(operatorCacheFilePath(baseDir))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return operatorCachePayload{}, errOperatorCacheNotFound
		}
		return operatorCachePayload{}, err
	}

	var cache operatorCachePayload
	if err := json.Unmarshal(raw, &cache); err != nil {
		return operatorCachePayload{}, fmt.Errorf("failed to parse operator cache file: %w", err)
	}
	return cache, nil
}

func cacheOperatorImages(client *http.Client, baseDir string, operators []OperatorRecord) []OperatorRecord {
	result := make([]OperatorRecord, len(operators))
	copy(result, operators)

	for index := range result {
		if result[index].RemoteImageURL == "" {
			continue
		}

		localPath, err := downloadOperatorImage(client, operatorImageDir(baseDir), result[index])
		if err != nil {
			continue
		}
		result[index].LocalImagePath = localPath
		result[index].LocalImageURL = localFileURL(localPath)
	}

	return result
}

func downloadOperatorImage(client *http.Client, imageDir string, operator OperatorRecord) (string, error) {
	response, err := client.Get(operator.RemoteImageURL)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected image status code: %d", response.StatusCode)
	}

	imageBytes, err := io.ReadAll(io.LimitReader(response.Body, 4<<20))
	if err != nil {
		return "", err
	}

	extension := imageExtension(operator.RemoteImageURL, response.Header.Get("Content-Type"))
	filename := fmt.Sprintf("%03d%s", operator.Order, extension)
	fullPath := filepath.Join(imageDir, filename)
	if err := os.WriteFile(fullPath, imageBytes, 0o644); err != nil {
		return "", err
	}
	return fullPath, nil
}

func imageExtension(rawURL, contentType string) string {
	parsedURL, err := url.Parse(rawURL)
	if err == nil {
		ext := strings.ToLower(path.Ext(parsedURL.Path))
		if ext != "" && len(ext) <= 5 {
			return ext
		}
	}

	extensions, err := mime.ExtensionsByType(strings.Split(contentType, ";")[0])
	if err == nil && len(extensions) > 0 {
		return extensions[0]
	}

	return ".jpg"
}

func applyLocalImageURLs(baseDir string, operators []OperatorRecord) {
	for index := range operators {
		if operators[index].LocalImagePath == "" {
			continue
		}

		if !filepath.IsAbs(operators[index].LocalImagePath) {
			operators[index].LocalImagePath = filepath.Join(baseDir, operators[index].LocalImagePath)
		}
		operators[index].LocalImageURL = localFileURL(operators[index].LocalImagePath)
	}
}

func localFileURL(fullPath string) string {
	normalized := filepath.ToSlash(fullPath)
	if !strings.HasPrefix(normalized, "/") {
		normalized = "/" + normalized
	}
	return (&url.URL{Scheme: "file", Path: normalized}).String()
}

func sortOperatorRecords(records []OperatorRecord) {
	sort.SliceStable(records, func(i, j int) bool {
		return records[i].Order < records[j].Order
	})
}

func extractRemoteImageURL(node *html.Node) string {
	imgNode := findFirstNode(node, func(current *html.Node) bool {
		return current.Type == html.ElementNode && current.Data == "img"
	})
	if imgNode == nil {
		return ""
	}

	rawURL := strings.TrimSpace(getAttr(imgNode, "src"))
	if rawURL == "" {
		return ""
	}
	return resolveSourceURL(rawURL)
}

func resolveSourceURL(rawURL string) string {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return rawURL
	}
	if parsedURL.IsAbs() {
		return parsedURL.String()
	}

	sourceURL, err := url.Parse(operatorDataSourceURL)
	if err != nil {
		return rawURL
	}
	return sourceURL.ResolveReference(parsedURL).String()
}

func splitMetadata(raw string) []string {
	if raw == "" {
		return nil
	}

	pieces := strings.Split(raw, ",")
	result := make([]string, 0, len(pieces))
	for _, piece := range pieces {
		trimmed := strings.TrimSpace(piece)
		if trimmed == "" {
			continue
		}
		result = append(result, trimmed)
	}
	return result
}

func isAcquisitionField(value string) bool {
	return value == "公开招募" ||
		strings.Contains(value, "寻访") ||
		strings.Contains(value, "获得") ||
		strings.Contains(value, "掉落") ||
		strings.Contains(value, "兑换")
}

func isDisplayedTag(value string, tagSet map[string]struct{}) bool {
	_, ok := tagSet[value]
	return ok
}

func appendIfMissing(values []string, target string) []string {
	if contains(values, target) {
		return values
	}
	return append(values, target)
}

func contains(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}

func findNodesByClass(root *html.Node, className string) []*html.Node {
	return findNodes(root, func(node *html.Node) bool {
		return node.Type == html.ElementNode && hasClass(node, className)
	})
}

func findFirstNodeByClass(root *html.Node, className string) *html.Node {
	return findFirstNode(root, func(node *html.Node) bool {
		return node.Type == html.ElementNode && hasClass(node, className)
	})
}

func findNodes(root *html.Node, match func(*html.Node) bool) []*html.Node {
	var matches []*html.Node
	var walk func(*html.Node)
	walk = func(node *html.Node) {
		if match(node) {
			matches = append(matches, node)
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			walk(child)
		}
	}
	walk(root)
	return matches
}

func findFirstNode(root *html.Node, match func(*html.Node) bool) *html.Node {
	if root == nil {
		return nil
	}
	if match(root) {
		return root
	}
	for child := root.FirstChild; child != nil; child = child.NextSibling {
		if result := findFirstNode(child, match); result != nil {
			return result
		}
	}
	return nil
}

func hasClass(node *html.Node, className string) bool {
	for _, attr := range node.Attr {
		if attr.Key != "class" {
			continue
		}
		for _, class := range strings.Fields(attr.Val) {
			if class == className {
				return true
			}
		}
	}
	return false
}

func getAttr(node *html.Node, key string) string {
	for _, attr := range node.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}

func extractText(node *html.Node) string {
	if node == nil {
		return ""
	}

	var builder strings.Builder
	var walk func(*html.Node)
	walk = func(current *html.Node) {
		if current.Type == html.TextNode {
			builder.WriteString(current.Data)
		}
		for child := current.FirstChild; child != nil; child = child.NextSibling {
			walk(child)
		}
	}
	walk(node)
	return builder.String()
}
