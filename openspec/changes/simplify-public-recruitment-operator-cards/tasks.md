## 1. Data Preparation

- [x] 1.1 Add operator image loading state to `frontend/src/components/PublicRecruitmentPage.vue` and resolve unique `localImagePath` values into base64 data URLs.
- [x] 1.2 Add helper logic in `frontend/src/components/PublicRecruitmentPage.vue` to return each operator's image source and rarity border color from the required mapping.

## 2. Result Card Simplification

- [x] 2.1 Update the Combination result card template in `frontend/src/components/PublicRecruitmentPage.vue` so each card only renders the operator image area and operator name.
- [x] 2.2 Remove rarity text, profession/origin text, and recruitment tag chips from the Combination result cards while preserving no-image fallback rendering.

## 3. Visual Styling And Verification

- [x] 3.1 Adjust the result card CSS in `frontend/src/components/PublicRecruitmentPage.vue` to support a compact image-first layout and the specified rarity border colors.
- [x] 3.2 Run the frontend build and verify the public recruitment page still compiles with the simplified operator cards.

