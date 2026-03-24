## 1. Navigation shell and theme refresh

- [ ] 1.1 Replace the fixed left-side navigation with an Element Plus Drawer triggered from the main shell
- [ ] 1.2 Keep page switching and current-page highlighting working correctly inside the drawer flow
- [ ] 1.3 Update the global shell and operator data page styling to a light blue and white theme

## 2. Operator cache data model and storage

- [ ] 2.1 Extend the operator data model with source order, remote image URL, local image path, and cache payload fields
- [ ] 2.2 Choose and implement a writable local cache directory for operator data JSON and downloaded images
- [ ] 2.3 Save parsed operator records to a local cache file after a successful fetch
- [ ] 2.4 Download and store operator images locally without blocking cache writes when an individual image fails

## 3. Cache-first loading behavior

- [ ] 3.1 Add backend logic to read cached operator data before any remote fetch is requested
- [ ] 3.2 Initialize the operator data page by loading local cache first and falling back to the empty state when no valid cache exists
- [ ] 3.3 Keep the remote fetch action available so the user can refresh data and overwrite the local cache explicitly

## 4. Source-order preservation

- [ ] 4.1 Preserve `.contentDetail` parse order when generating operator records from the wiki HTML
- [ ] 4.2 Persist and reload the preserved order through the local cache format
- [ ] 4.3 Render the operator list in preserved source order without applying derived sorting in the UI

## 5. Verification

- [ ] 5.1 Add or update backend tests for cache serialization, source-order preservation, and image/cache fallback cases
- [ ] 5.2 Run `go test ./...` and fix any backend regressions introduced by the cache and persistence work
- [ ] 5.3 Run `pnpm build` in `frontend/` and fix any frontend regressions from the Drawer and theme changes
- [ ] 5.4 Verify the packaged Wails app still starts and shows cached operator data correctly when local cache is present
