## 1. Data and domain foundation

- [ ] 1.1 Define the local data format for recruitment tags, operators, rarity, and recruitable tag relationships
- [ ] 1.2 Add an initial offline recruitment data set to the repository and document how it is maintained
- [ ] 1.3 Implement Go domain models for screenshot analysis input, recognized tags, and recommendation results

## 2. Screenshot input pipeline

- [ ] 2.1 Add a backend API that accepts a screenshot file path or image payload for recruitment analysis
- [ ] 2.2 Implement image decoding and basic screenshot validation for recruitment screen readiness
- [ ] 2.3 Return clear validation errors when the provided image cannot be analyzed

## 3. Tag recognition pipeline

- [ ] 3.1 Integrate an OCR or text-recognition adapter behind a Go interface for recruitment tag extraction
- [ ] 3.2 Normalize OCR output into standard recruitment tags with confidence and unresolved states
- [ ] 3.3 Add test coverage for tag normalization, low-confidence handling, and unmatched OCR text

## 4. Recommendation engine

- [ ] 4.1 Implement combination generation for confirmed single-tag, two-tag, and three-tag inputs
- [ ] 4.2 Implement operator matching against the local recruitment data set
- [ ] 4.3 Implement recommendation sorting and explanation generation based on rarity and certainty
- [ ] 4.4 Add backend tests for high-value, empty-result, and tie-breaking recommendation cases

## 5. Desktop UI flow

- [ ] 5.1 Replace or extend the current frontend with a recruitment analysis workflow centered on screenshot input
- [ ] 5.2 Add a recognition review UI that shows detected tags, confidence state, and manual correction controls
- [ ] 5.3 Add a recommendation results UI that lists tag combinations, candidate operators, and recommendation reasons
- [ ] 5.4 Handle loading, invalid screenshot, unresolved recognition, and empty recommendation states in the UI

## 6. Verification and cleanup

- [ ] 6.1 Run `go test ./...` and fix compile or logic regressions introduced by the new flow
- [ ] 6.2 Run `pnpm build` in `frontend/` and verify the Wails app can build with the new UI flow
- [ ] 6.3 Review whether legacy window-tool entry points should be hidden, removed, or left isolated after the recruitment flow is working
