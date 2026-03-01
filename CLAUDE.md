# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project

Pingu is a Go CLI tool. Module path: `github.com/512uk/pingu`.

## Commands

```bash
make build    # Build binary to bin/pingu
make run      # Run without building
make test     # Run all tests
make lint     # Run go vet
make clean    # Remove build artifacts
```

Run a single test:
```bash
go test ./internal/... -run TestName
```

## Architecture

- `cmd/pingu/main.go` — Entrypoint. Calls `run()` which returns an error for clean exit handling.
- `internal/` — Private application logic. Add packages here as the project grows.
