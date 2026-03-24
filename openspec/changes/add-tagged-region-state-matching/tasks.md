## 1. Shared Tag Catalog

- [ ] 1.1 Extract the public recruitment tag groups and values into a shared backend/frontend-friendly catalog structure.
- [ ] 1.2 Expose a data access path so the recognition settings page can load the shared recruitment tag catalog instead of maintaining its own local list.

## 2. Template Data Model

- [ ] 2.1 Replace the single-reference region model with a region-state collection model that stores multiple tagged screenshots under one region.
- [ ] 2.2 Update template save and load logic so new templates persist nested region states and legacy single-state templates remain readable.
- [ ] 2.3 Update region image storage layout so each saved state keeps its own reference image path and metadata.

## 3. Matching Results

- [ ] 3.1 Update the recognition matching pipeline to compare each region against all of its saved state screenshots.
- [ ] 3.2 Return region-level match results that include the matched state identifiers and recruitment tags, including the no-match case.

## 4. Recognition Settings UI

- [ ] 4.1 Update the recognition settings page to manage multiple state screenshots under one region, including add, review, and delete actions.
- [ ] 4.2 Add recruitment tag selection for each region state by using the shared recruitment tag catalog.
- [ ] 4.3 Update the match result display so users can see which tag state matched successfully for each region.

## 5. Verification

- [ ] 5.1 Add or update tests for shared tag catalog data, template backward compatibility, and multi-state region matching results.
- [ ] 5.2 Verify the recognition settings workflow for loading tags, saving multi-state regions, reopening templates, and viewing matched tag states.
