# Repository Guidelines

## Project Structure & Module Organization
- `cmd/`: Cobra CLI entrypoints and subcommands.
- `internal/app/`: Application logic and DTOs; parsing lives in `internal/app/parse/`.
- `internal/domain/`: Domain types like patch intents, values, and resource refs.
- `main.go`: CLI bootstrap and wiring.
- Tests are colocated with code as `*_test.go` (e.g., `internal/app/parse/value_parser_test.go`).

## Architecture & DDD Boundaries
- Treat `internal/domain/` as the pure domain: no Cobra, CLI parsing, or infrastructure concerns.
- Keep invariants and validation in domain constructors (e.g., `NewPatch`, `NewResourceRef`).
- `internal/app/` orchestrates use cases and maps input DTOs to domain objects.
- `cmd/` is the delivery layer; it should only assemble DTOs and call app handlers.
- Domain depends on nothing outside `internal/domain/`; app can depend on domain; cmd depends on app.

## Build, Test, and Development Commands
- `go run main.go <command> <subcommand>`: Run the CLI directly for local dev.
- `go build -o kubefuse`: Build the binary.
- `./kubefuse <command> <subcommand>`: Run the built binary.
- `go test ./...`: Run the full test suite.
- `go test -tags=integration ./internal/integration`: Run envtest-based integration tests (requires `KUBEBUILDER_ASSETS`).
- `scripts/kind-up.sh` and `scripts/kind-down.sh`: Start/stop a local kind cluster for end-to-end testing.
 - Releases are created by tagging (e.g., `v0.1.0`) and pushing the tag.

## Coding Style & Naming Conventions
- Go standard formatting; run `gofmt` on all `.go` files.
- Exported identifiers use `CamelCase`; unexported use `camelCase`.
- Filenames are lowercase with underscores for tests (`value_parser_test.go`).
- Keep parsing logic in `internal/app/parse/` and domain types in `internal/domain/`.
- Avoid leaking CLI flags or raw strings into the domain; map to domain types first.

## Testing Guidelines
- Go’s `testing` package is used; tests live beside source.
- Name tests as `TestXxx` and files `*_test.go`.
- Prefer focused unit tests for parsers and domain logic.
- Run all tests with `go test ./...` before PRs.
- For Kubernetes API integration, use `envtest` with `scripts/setup-envtest.sh` to set `KUBEBUILDER_ASSETS`.

## Commit & Pull Request Guidelines
- Recent commits are short, lowercase, and descriptive (e.g., “moved parsers to app layer”); follow that style unless the team adopts a standard.
- PRs should include: concise summary, test command(s) run, and any behavior changes. Link related issues if applicable.

## Configuration & Usage Notes
- Typical usage: `kubefuse set <kind/name> <path=value>... [--ttl 10m] [--reason "..."] [--dry-run]`.
- Namespace flags follow Kubernetes conventions (e.g., `-n prod`).
