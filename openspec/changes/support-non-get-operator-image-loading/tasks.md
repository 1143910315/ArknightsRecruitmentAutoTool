## 1. Backend image retrieval flow

- [ ] 1.1 Remove the current dependence on GET-style cached image asset URLs for operator images
- [ ] 1.2 Add a Wails-compatible backend method that reads a cached operator image from the runtime directory and returns structured image content data
- [ ] 1.3 Ensure missing cached images return a graceful not-found style result rather than a fatal error

## 2. Frontend image rendering integration

- [ ] 2.1 Update the operator data page to request cached image content through the new backend-supported flow
- [ ] 2.2 Convert the returned image content into a renderable frontend image source and display it in the operator list
- [ ] 2.3 Preserve the existing missing-image fallback UI when no cached image content is available

## 3. Behavior preservation

- [ ] 3.1 Keep runtime-directory cache storage unchanged while replacing only the image access mechanism
- [ ] 3.2 Verify that cached operator data still loads first when the page opens and that operator order remains unchanged

## 4. Verification

- [ ] 4.1 Add or update backend tests for cached image content retrieval, missing-image handling, and runtime-directory path resolution
- [ ] 4.2 Run `go test ./...` and fix any regressions introduced by the new image retrieval flow
- [ ] 4.3 Run `pnpm build` in `frontend/` and fix any frontend regressions from the new image rendering path
- [ ] 4.4 Verify the packaged Wails app starts and can display cached operator images without relying on GET asset requests
