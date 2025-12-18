# Techstack of dendrite-pulse

- Language: Golang Version 1.25
- Linting: `golangci-lint` version 2.5.0
- VCS: `git` and [GitHub](https://github.com/thorstenkramm/fillfs)
- CI/CD: GitHub Actions

## Mandatory Go modules and packages

- Use [viper](https://github.com/spf13/viper) to handle command line arguments and flags.

## Testing and linting Guidelines

- Target at least 75% coverage (see `rules/testing.md`); favor table-driven tests for handlers and services.
- Use `github.com/stretchr/testify/assert`/`require` in `_test.go` files; prefer `require.NoError` for setup and
  `assert` for behavioral checks.
- `go test -race -v ./...` to run the full suite with race detection (required after each task).
- `golangci-lint run ./...` (v2.5.0) for linting; add linters to `.golangci.yml` if configuration is introduced
  (required after each task).
- Search for code duplication with using [JSCPD](https://github.com/kucherenko/jscpd) and the command
  `npx jscpd --pattern "**/*.go" --ignore "**/*_test.go" --threshold 0 --exitCode 1`
- `go fmt ./...` and `go vet ./...` to keep code idiomatic before committing.
- Always run all tests after a task is completed.
- Make sure the test coverage doesn't fall below the required threshold of 75% after a task is completed.

## Target OS

fillfs is intented to be used on Linux and macOS. Windows is out of scope.

## Other rules

- Use Go's built-in concurrency features when beneficial for performance.

## Docker & k8s

As go creates self-contained and easy to deploy single-file binaries, containerized deployments are not planned.
Also, instructions on how to develop having the go compiler inside a container are not provided.
All documentation assumes you have the `go` command line utility installed directly to your OS.
