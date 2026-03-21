# claude-television (ctv)

A read-only CLI tool that explores the local configuration state of Claude Code with a TUI dashboard.

## Build and Run

- `go build -o ctv .` — Build
- `go run .` — Run
- `go test ./...` — Run all tests
- `go test ./internal/claude/... -v` — Test claude package
- `go test ./internal/scanner/... -v` — Test scanner package

## Project Structure

- `cmd/` — Cobra CLI commands
- `internal/claude/` — Claude Code configuration file parsing (settings, CLAUDE.md, plugins, skills, hooks, projects)
- `internal/scanner/` — Project directory scanning
- `internal/config/` — Viper-based ctv self-configuration management
- `internal/tui/` — Bubble Tea TUI components

## Coding Conventions

- Go standard project layout (`cmd/`, `internal/`)
- Error handling: Use `fmt.Errorf("context: %w", err)` pattern
- Testing: Place fixture files in `testdata/` directory, prefer table-driven tests
- Minimize external dependencies, add only when necessary
