## 1. Runtime directory cache migration

- [ ] 1.1 Replace the current user-directory cache root logic with a cache directory under the program runtime directory
- [ ] 1.2 Update cache file and image path handling to store relative paths appropriate for the runtime directory layout
- [ ] 1.3 Ensure cache reads and writes continue to preserve local-first behavior and source order after the directory move

## 2. Wails-served local asset handling

- [ ] 2.1 Add a Wails asset handling path or equivalent resource-serving hook for cached operator images under the runtime directory
- [ ] 2.2 Return Wails-resolvable local asset paths to the frontend instead of direct local file URLs
- [ ] 2.3 Keep operator records renderable when no local image can be served

## 3. Frontend integration

- [ ] 3.1 Update the operator data page to consume Wails-served image paths for cached operator images
- [ ] 3.2 Verify that cached operator data still loads first when the page opens and that the list order remains unchanged

## 4. Verification

- [ ] 4.1 Add or update backend tests for runtime-directory cache path generation, cache reload behavior, and Wails-served asset path generation
- [ ] 4.2 Run `go test ./...` and fix any regressions introduced by the cache location and asset-serving changes
- [ ] 4.3 Run `pnpm build` in `frontend/` and fix any frontend regressions from the new local asset path flow
- [ ] 4.4 Verify the packaged Wails app starts and can display cached operator data and images from the program runtime directory
