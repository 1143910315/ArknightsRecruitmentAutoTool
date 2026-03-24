## 1. Navigation shell and theme refresh

- [x] 1.1 Replace the fixed left-side navigation with an Element Plus Drawer triggered from the main shell
- [x] 1.2 Keep page switching and current-page highlighting working correctly inside the drawer flow
- [x] 1.3 Update the global shell and operator data page styling to a light blue and white theme

## 2. Operator cache data model and storage

- [x] 2.1 Extend the operator data model with source order, remote image URL, local image path, and cache payload fields
- [x] 2.2 Choose and implement a writable local cache directory for operator data JSON and downloaded images
- [x] 2.3 Save parsed operator records to a local cache file after a successful fetch
- [x] 2.4 Download and store operator images locally without blocking cache writes when an individual image fails

## 3. Cache-first loading behavior

- [x] 3.1 Add backend logic to read cached operator data before any remote fetch is requested
- [x] 3.2 Initialize the operator data page by loading local cache first and falling back to the empty state when no valid cache exists
- [x] 3.3 Keep the remote fetch action available so the user can refresh data and overwrite the local cache explicitly

## 4. Source-order preservation

- [x] 4.1 Preserve `.contentDetail` parse order when generating operator records from the wiki HTML
- [x] 4.2 Persist and reload the preserved order through the local cache format
- [x] 4.3 Render the operator list in preserved source order without applying derived sorting in the UI

## 5. Verification

- [x] 5.1 Add or update backend tests for cache serialization, source-order preservation, and image/cache fallback cases
- [x] 5.2 Run `go test ./...` and fix any backend regressions introduced by the cache and persistence work
- [x] 5.3 Run `pnpm build` in `frontend/` and fix any frontend regressions from the Drawer and theme changes
- [x] 5.4 Verify the packaged Wails app still starts and shows cached operator data correctly when local cache is present
