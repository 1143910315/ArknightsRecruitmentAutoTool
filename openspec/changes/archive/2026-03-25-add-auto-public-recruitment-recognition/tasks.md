## 1. Backend Window Resolution

- [x] 1.1 Extend recognition template persistence to save window-instance metadata needed for duplicate title/class disambiguation
- [x] 1.2 Implement target-window resolution that prefers a live recorded handle, falls back to stored instance metadata, and returns explicit missing/ambiguous failure states
- [x] 1.3 Add backend tests covering exact-handle reuse, unique fallback resolution, and ambiguous-window rejection

## 2. Backend Recruitment Recognition Aggregation

- [x] 2.1 Add a public recruitment recognition backend API that resolves the target window, captures it, and evaluates all configured regions
- [x] 2.2 Aggregate region match results into a strict success/failure contract that only returns tags when every region uniquely matches one state
- [x] 2.3 Enforce recruitment-tag-only output, include failure reasons for no-op cases, and add tests for zero-match and multi-match region failures

## 3. Public Recruitment Page Integration

- [x] 3.1 Add an automatic recognition toggle and status area to the public recruitment page
- [x] 3.2 Implement a serial recognition loop that triggers only while the page toggle is enabled and waits 500ms after each completed run before retrying
- [x] 3.3 Route successful recognition tags through the existing `selectedTags` state so combination generation and result rendering reuse current logic
- [x] 3.4 Keep current page selection unchanged on unusable recognition runs and stop scheduling when the page unmounts or the toggle is turned off

## 4. Verification

- [x] 4.1 Verify no-template, missing-window, ambiguous-window, unique-success, zero-match, and multi-match scenarios end-to-end
- [x] 4.2 Run `go test ./...` and fix backend regressions introduced by auto recognition changes
- [x] 4.3 Run `pnpm exec vite build` and `wails build` to confirm the desktop app still builds with the new public recruitment flow
