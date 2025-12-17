# Repository Guidelines

## Project Structure & Module Organization

- Root will host Go CLI sources alongside config and docs; see `README.md` for CLI usage examples.
- `samples/` holds fixture files used when generating realistic file trees.
- `rules/` documents tech stack expectations (`rules/techstack.md`) and Markdown standards (`rules/markdown.md`).
- `.golangci.yml` and `.markdownlint.json` define linting rules; `.github/` houses CI workflows if present.

## Build, Test, and Development Commands

- `go build ./...` — compile all modules to ensure the tree is healthy.
- `go test -race -v ./...` — full test suite with race detection; run before each PR.
- `golangci-lint run ./...` — lint using version 2.5.0 config in the repo.
- `go fmt ./...` and `go vet ./...` — format and vet prior to linting.
- Run the CLI locally with
  `./fillfs --dest ./fakefs --cache-dir /tmp/fillfs --folders 2 --files-per-folder 10 --depths 1`
  after building the binary; adjust flags as needed.

## Coding Style & Naming Conventions

- Go code must be `gofmt`-clean; use tabs (default Go style) and idiomatic naming (exported
  identifiers in PascalCase, locals in camelCase, package names lower_snake).
- Prefer table-driven tests and `require/assert` from `testify`.
- Keep CLI flag handling centralized with `viper` as noted in `rules/techstack.md`.
- Document non-obvious logic with short comments; keep line length ≤120 chars for Markdown
  per `rules/markdown.md`.

## Testing Guidelines

- Target ≥75% coverage; add or update tests alongside code changes.
- Place tests in `_test.go` files near the code; name functions `TestXYZ` with clear case coverage.
- For behavior that spawns files, use temp dirs and clean up; avoid touching `samples/`.

## Commit & Pull Request Guidelines

- Use concise, imperative commit messages (e.g., `Add cache warmup step`); group related changes.
- PRs should describe intent, summarize implementation, and list test commands run
  (`go test -race -v ./...`, `golangci-lint run ./...`).
- Link related issues or tickets and include screenshots only if UX-visible output changes.

## Safety & Environment Notes

- Filling the filesystem can grow quickly; verify disk space before running deep/wide configurations.
- The cache directory must be created by `fillfs`; pass `--clean-cache true` when you need a fresh start.
- When editing docs, run `npx markdownlint --fix AGENTS.md` to respect repo Markdown standards.
