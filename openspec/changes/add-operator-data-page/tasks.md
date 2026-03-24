## 1. Application shell and navigation

- [ ] 1.1 Refactor the frontend app layout into a left-side navigation shell that can host multiple top-level pages
- [ ] 1.2 Move the existing window tools view into its own page panel under the new navigation structure
- [ ] 1.3 Add a `干员数据` navigation item and keep the active page state stable during page actions

## 2. Operator data domain and backend fetch flow

- [ ] 2.1 Define Go models for operator records, parsed metadata fields, page state payloads, and fetch results returned to the frontend
- [ ] 2.2 Implement a Wails backend method that performs HTTP GET to `https://wiki.biligame.com/arknights/公开招募工具`
- [ ] 2.3 Add error handling for network failures, non-success responses, and empty or invalid HTML payloads

## 3. Wiki HTML parsing

- [ ] 3.1 Implement parser logic that iterates over `.contentDetail` entries and extracts operator names from `.picText`
- [ ] 3.2 Parse rarity from `data-param2`, visible tags from `.tags .tagText`, and split `data-param1` into structured metadata fields
- [ ] 3.3 Preserve public recruitment availability and any remaining source metadata needed for later recruitment logic reuse
- [ ] 3.4 Add parser-focused tests or fixture-based checks using representative HTML snippets for success and malformed input cases

## 4. Operator data page UI

- [ ] 4.1 Build the `干员数据` page with explicit empty, loading, success, and error states
- [ ] 4.2 Add the `获取干员数据` button and wire it to the new Wails backend method
- [ ] 4.3 Render the fetched operator list with at least name, rarity, and parsed tags for each record
- [ ] 4.4 Ensure fetch failures do not clear previously loaded data and surface a clear user-facing error message

## 5. Verification

- [ ] 5.1 Run `go test ./...` and fix any backend compile or parsing regressions introduced by the change
- [ ] 5.2 Run `pnpm build` in `frontend/` and fix any frontend build issues introduced by the new layout and page
- [ ] 5.3 Smoke test the Wails app to verify navigation switching and operator data fetching behavior work together
