## 1. Template Data Model

- [ ] 1.1 Extend the recognition template model and persistence format to store a template-level image matching tolerance value
- [ ] 1.2 Normalize legacy templates so missing tolerance values default to 0 during load and remain writable in the new format
- [ ] 1.3 Add backend validation to reject invalid tolerance values before saving a template

## 2. Tolerance-Based Matching Logic

- [ ] 2.1 Replace exact pixel equality matching with per-pixel absolute difference comparison against the configured tolerance value
- [ ] 2.2 Keep zero tolerance behavior equivalent to the current exact-match logic and ensure region-level matching stays deterministic
- [ ] 2.3 Add backend tests covering exact match, within-tolerance match, out-of-tolerance mismatch, and legacy default tolerance behavior

## 3. Recognition Settings UI

- [ ] 3.1 Add an image matching tolerance input and explanatory copy to the recognition settings page
- [ ] 3.2 Load and save the configured tolerance value together with the rest of the recognition template fields
- [ ] 3.3 Prevent invalid tolerance input in the UI and show clear feedback before save

## 4. Verification

- [ ] 4.1 Verify the recognition settings flow for new templates, legacy templates, zero tolerance, and non-zero tolerance cases
- [ ] 4.2 Run `go test ./...` and fix any regressions introduced by tolerance-based matching
- [ ] 4.3 Run `pnpm exec vite build` and `wails build` to confirm the desktop app still builds successfully
