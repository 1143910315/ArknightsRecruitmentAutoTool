package main

import (
	"strings"
	"testing"
)

func TestParseOperatorDataHTML(t *testing.T) {
	const sampleHTML = `
<div class="contentDetail" data-param1="狙击, 女, 6, , 高级资深干员, 拉特兰, 远程位, 输出, 中坚寻访, 公开招募, 否" data-param2="6">
<p class="picture"><a href="/arknights/%E7%A9%BA%E5%BC%A6" title="空弦" class="pic-star-6"><img alt="空弦hd.jpg" src="https://patchwiki.biligame.com/images/arknights/thumb/f/fe/rzrnsqjorinzm92f3da769vwsghca32.jpg/60px-%E7%A9%BA%E5%BC%A6hd.jpg" decoding="async" loading="lazy" width="60" height="60" srcset="https://patchwiki.biligame.com/images/arknights/thumb/f/fe/rzrnsqjorinzm92f3da769vwsghca32.jpg/90px-%E7%A9%BA%E5%BC%A6hd.jpg 1.5x, https://patchwiki.biligame.com/images/arknights/thumb/f/fe/rzrnsqjorinzm92f3da769vwsghca32.jpg/120px-%E7%A9%BA%E5%BC%A6hd.jpg 2x" data-file-width="180" data-file-height="180"></a><span class="picText">空弦</span></p>
<p class="tags">
<span class="btn btn-default tagText"> 远程位</span><span class="btn btn-default tagText">输出</span>
</p>
</div>
<div class="contentDetail" data-param1="近卫, 女, 5, 资深干员, , 哥伦比亚, 近战位, 输出, 防护, 公开招募, 中坚寻访, 是" data-param2="5">
<p class="picture"><a href="/arknights/%E6%98%9F%E6%9E%81" title="星极" class="pic-star-5"><img alt="星极hd.jpg" src="https://patchwiki.biligame.com/images/arknights/thumb/0/03/cedol3f4khu0nmue3wozfs6z2gcp0xg.jpg/60px-%E6%98%9F%E6%9E%81hd.jpg" decoding="async" loading="lazy" width="60" height="60" srcset="https://patchwiki.biligame.com/images/arknights/thumb/0/03/cedol3f4khu0nmue3wozfs6z2gcp0xg.jpg/90px-%E6%98%9F%E6%9E%81hd.jpg 1.5x, https://patchwiki.biligame.com/images/arknights/thumb/0/03/cedol3f4khu0nmue3wozfs6z2gcp0xg.jpg/120px-%E6%98%9F%E6%9E%81hd.jpg 2x" data-file-width="180" data-file-height="180"></a><span class="picText">星极</span></p>
<p class="tags">
<span class="btn btn-default tagText"> 近战位</span><span class="btn btn-default tagText">输出</span><span class="btn btn-default tagText">防护</span>
</p>
</div>
<div class="contentDetail" data-param1="特种, 女, 4, , , 炎-龙门, 位移, 近战位, 关卡1-12首次通关掉落, 公开招募, 标准寻访, 中坚寻访, 主题曲获得, 否" data-param2="4">
<p class="picture"><a href="/arknights/%E9%98%BF%E6%B6%88" title="阿消" class="pic-star-4"><img alt="阿消hd.jpg" src="https://patchwiki.biligame.com/images/arknights/thumb/1/1a/jsuvpeoduxn9ne07oc8c7s7vwngz4sl.jpg/60px-%E9%98%BF%E6%B6%88hd.jpg" decoding="async" loading="lazy" width="60" height="60" srcset="https://patchwiki.biligame.com/images/arknights/thumb/1/1a/jsuvpeoduxn9ne07oc8c7s7vwngz4sl.jpg/90px-%E9%98%BF%E6%B6%88hd.jpg 1.5x, https://patchwiki.biligame.com/images/arknights/thumb/1/1a/jsuvpeoduxn9ne07oc8c7s7vwngz4sl.jpg/120px-%E9%98%BF%E6%B6%88hd.jpg 2x" data-file-width="180" data-file-height="180"></a><span class="picText">阿消</span></p>
<p class="tags">
<span class="btn btn-default tagText"> 位移</span><span class="btn btn-default tagText">近战位</span>
</p>
</div>`

	operators, err := parseOperatorDataHTML(strings.NewReader(sampleHTML))
	if err != nil {
		t.Fatalf("parseOperatorDataHTML returned error: %v", err)
	}

	if len(operators) != 3 {
		t.Fatalf("expected 3 operators, got %d", len(operators))
	}

	first := operators[0]
	if first.Name != "空弦" {
		t.Fatalf("expected first operator to be 空弦, got %s", first.Name)
	}
	if first.Rarity != 6 {
		t.Fatalf("expected rarity 6, got %d", first.Rarity)
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
	if !contains(third.Metadata.AcquisitionMethods, "关卡1-12首次通关掉落") {
		t.Fatalf("expected acquisition method to be preserved, got %#v", third.Metadata.AcquisitionMethods)
	}
	if !contains(third.Metadata.AcquisitionMethods, "主题曲获得") {
		t.Fatalf("expected theme unlock acquisition to be preserved, got %#v", third.Metadata.AcquisitionMethods)
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
