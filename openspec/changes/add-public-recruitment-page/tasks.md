## 1. Navigation and Page Skeleton

- [x] 1.1 Add a public recruitment navigation entry and route it to a dedicated frontend page
- [x] 1.2 Create the public recruitment page with sections for tag selection, result groups, and missing-data guidance state
- [x] 1.3 Reuse existing local operator data loading on the new page so recruitment filtering works from cached data without changing fetch flow

## 2. Tag Selection and Filtering Logic

- [x] 2.1 Add a static recruitment tag configuration covering profession, deployment position, trait tags, and operator seniority
- [x] 2.2 Implement selected-tag state management with a hard limit of five simultaneous selections
- [x] 2.3 Show clear UI feedback when the user tries to select more than five tags and keep existing selections unchanged
- [x] 2.4 Normalize each operator record into a comparable recruitment-tag set using existing operator fields and only include publicly recruitable operators

## 3. Combination Results Rendering

- [x] 3.1 Generate all non-empty combinations from the selected tags in stable order, prioritizing longer combinations before shorter ones
- [x] 3.2 Filter out combinations that match no operators and keep only valid result groups
- [x] 3.3 Render each valid combination as its own result group and preserve operator source order within each group
- [x] 3.4 Ensure operators appear in every combination they satisfy, including both full-match and subset-match groups

## 4. Verification

- [x] 4.1 Verify empty-data, single-tag, multi-tag, and five-tag-limit scenarios on the public recruitment page
- [x] 4.2 Run `pnpm exec vite build` and fix any frontend regressions introduced by the new page
- [x] 4.3 Run `wails build` to confirm the app still packages successfully with the new recruitment page
