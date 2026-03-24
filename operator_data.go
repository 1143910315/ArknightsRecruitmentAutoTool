package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html"
)

const operatorDataSourceURL = "https://wiki.biligame.com/arknights/公开招募工具"

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
	Name                string           `json:"name"`
	Rarity              int              `json:"rarity"`
	DisplayTags         []string         `json:"displayTags"`
	IsPublicRecruitable bool             `json:"isPublicRecruitable"`
	Metadata            OperatorMetadata `json:"metadata"`
}

type FetchOperatorDataResult struct {
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

	return FetchOperatorDataResult{
		SourceURL: operatorDataSourceURL,
		FetchedAt: time.Now().Format(time.RFC3339),
		Operators: operators,
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
	for _, node := range nodes {
		record, err := parseOperatorRecord(node)
		if err != nil {
			return nil, err
		}
		operators = append(operators, record)
	}

	return operators, nil
}

func parseOperatorRecord(node *html.Node) (OperatorRecord, error) {
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
		Name:                name,
		Rarity:              rarity,
		DisplayTags:         displayTags,
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
	var matches []*html.Node
	var walk func(*html.Node)
	walk = func(node *html.Node) {
		if node.Type == html.ElementNode && hasClass(node, className) {
			matches = append(matches, node)
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			walk(child)
		}
	}
	walk(root)
	return matches
}

func findFirstNodeByClass(root *html.Node, className string) *html.Node {
	if root == nil {
		return nil
	}
	if root.Type == html.ElementNode && hasClass(root, className) {
		return root
	}
	for child := root.FirstChild; child != nil; child = child.NextSibling {
		if result := findFirstNodeByClass(child, className); result != nil {
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
