## 1. State Data Model

- [x] 1.1 Extend recognition state image input, persistence, and loaded template models to store a per-state tolerance value
- [x] 1.2 Normalize legacy templates so state images without a stored tolerance default to 0 during load and save in the new format
- [x] 1.3 Add backend validation to reject invalid tolerance values on individual state images before saving a template

## 2. State-Level Matching Logic

- [x] 2.1 Update region state matching so each saved state image is compared using its own tolerance value
- [x] 2.2 Keep zero tolerance behavior equivalent to exact matching for that specific state image
- [x] 2.3 Add backend tests covering different tolerance values for different states in one region, within-tolerance matches, out-of-tolerance mismatches, and legacy default tolerance behavior

## 3. Recognition Settings UI

- [x] 3.1 Add a tolerance input and explanatory copy to each recognition state editor in the recognition settings page
- [x] 3.2 Load and save each state's tolerance value together with its tag and screenshot data
- [x] 3.3 Prevent invalid per-state tolerance input in the UI and show clear feedback before save

## 4. Verification

- [x] 4.1 Verify the recognition settings flow for new templates, legacy templates, zero tolerance, and mixed tolerance values across multiple states
- [x] 4.2 Run `go test ./...` and fix any regressions introduced by per-state tolerance matching
- [x] 4.3 Run `pnpm exec vite build` and `wails build` to confirm the desktop app still builds successfully
