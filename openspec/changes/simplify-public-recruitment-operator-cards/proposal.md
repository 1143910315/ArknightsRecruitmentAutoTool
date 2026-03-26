## Why

公开招募页面的 Combination 结果当前展示了星级、职业/阵营和标签等大量文本信息，用户在快速浏览组合命中干员时需要额外筛选无关内容。现在需要把结果卡片收敛为更直观的头像加名称展示，并用统一的星级边框颜色直接传达稀有度，提升扫描效率。

## What Changes

- 简化公开招募页面中 Combination 结果里的干员卡片，只展示干员 base64 图片和干员名字。
- 为 Combination 结果中的干员图片增加按星级区分的边框颜色：1 星和 2 星使用 `#dedede`，3 星使用 `#618bf5`，4 星使用 `#8960ce`，5 星使用 `#f0a94d`，6 星使用 `#f0e028`。
- 保持现有组合筛选、排序和命中逻辑不变，仅调整结果展示信息密度与视觉编码方式。

## Capabilities

### New Capabilities
- None.

### Modified Capabilities
- `public-recruitment-page`: 调整公开招募页面的 Combination 结果展示形式，使干员卡片仅保留图片和名称。
- `recruitment-tag-combination-results`: 为组合结果中的干员展示补充基于星级的图片边框颜色要求，并约束结果卡片的最小信息集合。

## Impact

- Affected code: `frontend/src/components/PublicRecruitmentPage.vue`
- Affected behavior: 公开招募页面 Combination 结果区的干员卡片布局、文案密度与稀有度视觉表达
- No backend or API changes.
