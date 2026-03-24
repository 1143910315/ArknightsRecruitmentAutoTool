package main

type RecruitmentTagGroup struct {
	Key   string   `json:"key"`
	Label string   `json:"label"`
	Tags  []string `json:"tags"`
}

type RecruitmentTagCatalog struct {
	Groups []RecruitmentTagGroup `json:"groups"`
}

var recruitmentTagCatalog = RecruitmentTagCatalog{
	Groups: []RecruitmentTagGroup{
		{Key: "profession", Label: "\u804c\u4e1a", Tags: []string{"\u8fd1\u536b", "\u72d9\u51fb", "\u91cd\u88c5", "\u533b\u7597", "\u8f85\u52a9", "\u672f\u5e08", "\u7279\u79cd", "\u5148\u950b"}},
		{Key: "position", Label: "\u90e8\u7f72\u4f4d\u7f6e", Tags: []string{"\u8fd1\u6218\u4f4d", "\u8fdc\u7a0b\u4f4d"}},
		{Key: "traits", Label: "\u7279\u6027\u6807\u7b7e", Tags: []string{"\u63a7\u573a", "\u7206\u53d1", "\u6cbb\u7597", "\u652f\u63f4", "\u8d39\u7528\u56de\u590d", "\u8f93\u51fa", "\u751f\u5b58", "\u7fa4\u653b", "\u9632\u62a4", "\u51cf\u901f", "\u524a\u5f31", "\u5feb\u901f\u590d\u6d3b", "\u4f4d\u79fb", "\u53ec\u5524", "\u652f\u63f4\u673a\u68b0", "\u5143\u7d20"}},
		{Key: "seniority", Label: "\u5e72\u5458\u8d44\u5386", Tags: []string{"\u65b0\u624b", "\u8d44\u6df1\u5e72\u5458", "\u9ad8\u7ea7\u8d44\u6df1\u5e72\u5458"}},
	},
}

var recruitmentTagSet = buildRecruitmentTagSet(recruitmentTagCatalog)

func (a *App) GetRecruitmentTagCatalog() RecruitmentTagCatalog {
	return cloneRecruitmentTagCatalog(recruitmentTagCatalog)
}

func cloneRecruitmentTagCatalog(catalog RecruitmentTagCatalog) RecruitmentTagCatalog {
	groups := make([]RecruitmentTagGroup, 0, len(catalog.Groups))
	for _, group := range catalog.Groups {
		groups = append(groups, RecruitmentTagGroup{
			Key:   group.Key,
			Label: group.Label,
			Tags:  append([]string(nil), group.Tags...),
		})
	}
	return RecruitmentTagCatalog{Groups: groups}
}

func buildRecruitmentTagSet(catalog RecruitmentTagCatalog) map[string]struct{} {
	values := make(map[string]struct{})
	for _, group := range catalog.Groups {
		for _, tag := range group.Tags {
			values[tag] = struct{}{}
		}
	}
	return values
}

func isRecruitmentTag(tag string) bool {
	_, ok := recruitmentTagSet[tag]
	return ok
}
