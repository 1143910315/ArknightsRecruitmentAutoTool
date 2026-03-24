## 1. Navigation and Page Skeleton

- [ ] 1.1 Add a public recruitment navigation entry and route it to a dedicated frontend page
- [ ] 1.2 Create the public recruitment page with sections for tag selection, result groups, and missing-data guidance state
- [ ] 1.3 Reuse existing local operator data loading on the new page so recruitment filtering works from cached data without changing fetch flow

## 2. Tag Selection and Filtering Logic

- [ ] 2.1 Add a static recruitment tag configuration covering profession, deployment position, trait tags, and operator seniority
- [ ] 2.2 Implement selected-tag state management with a hard limit of five simultaneous selections
- [ ] 2.3 Show clear UI feedback when the user tries to select more than five tags and keep existing selections unchanged
- [ ] 2.4 Normalize each operator record into a comparable recruitment-tag set using existing operator fields and only include publicly recruitable operators

## 3. Combination Results Rendering

- [ ] 3.1 Generate all non-empty combinations from the selected tags in stable order, prioritizing longer combinations before shorter ones
- [ ] 3.2 Filter out combinations that match no operators and keep only valid result groups
- [ ] 3.3 Render each valid combination as its own result group and preserve operator source order within each group
- [ ] 3.4 Ensure operators appear in every combination they satisfy, including both full-match and subset-match groups

## 4. Verification

- [ ] 4.1 Verify empty-data, single-tag, multi-tag, and five-tag-limit scenarios on the public recruitment page
- [ ] 4.2 Run `pnpm exec vite build` and fix any frontend regressions introduced by the new page
- [ ] 4.3 Run `wails build` to confirm the app still packages successfully with the new recruitment page
