# Repository Guidelines

## Project Structure & Module Organization
`main.go` contains the Wails desktop entrypoint and Windows integration logic. The frontend lives in `frontend/`: `src/` for Vue code, `assets/` for fonts and images, and `dist/` for the built web bundle embedded by Wails. Native packaging assets are in `build/`, including Windows icons and manifests. Product planning artifacts live in `openspec/`; use `openspec/changes/` for active work and `openspec/specs/` for accepted behavior.

## Build, Test, and Development Commands
Run backend and frontend together with `wails dev` from the repository root. Build a desktop binary with `wails build`. In `frontend/`, use `pnpm install` to sync dependencies, `pnpm dev` for the Vite dev server, and `pnpm build` to produce `frontend/dist`. For Go dependency checks, use `go test ./...` even when no tests exist yet; it still catches compile regressions.

## Coding Style & Naming Conventions
Use Go defaults: tabs for indentation, `gofmt` formatting, exported identifiers in `PascalCase`, and internal helpers in `camelCase`. Keep Windows API wrappers small and explicit. In Vue, prefer single-file components in `PascalCase` when components are introduced, and use `camelCase` for JS variables and functions. Keep imports grouped and minimal. Do not commit generated directories listed in `.gitignore`: `build/`, `frontend/dist/`, `frontend/node_modules/`, and `frontend/wailsjs/`.

## Testing Guidelines
There is no committed test suite yet. Add backend tests as `*_test.go` beside the code they cover, and frontend tests only if a test runner is added intentionally. Before opening a PR, run `go test ./...`, `pnpm build`, and a quick `wails dev` smoke test on Windows because this app depends on Win32 behavior.

## Commit & Pull Request Guidelines
Recent history uses short Chinese subjects such as `初始化项目` and `创建openspec`. Keep commit messages brief, imperative, and focused on one change. Pull requests should include: a short summary, impacted areas (`main.go`, `frontend/src`, `openspec/`), linked issues or change proposals, and screenshots or recordings for UI changes. Call out any Windows-specific assumptions or manual verification steps.

## Configuration Notes
Project config is centered in `wails.json` and `frontend/package.json`. If you change build commands or asset paths, update both the config and this guide when the contributor workflow changes.
