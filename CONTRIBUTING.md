# Contributing

Thank you for contributing.

## Development Setup

1. Install Go 1.22+.
2. Ensure cgo is available with a working C compiler.
3. Clone and initialize submodules:
```bash
git submodule update --init --recursive
```

4. Run tests:
```bash
go test ./...
```

## Code Style

- Keep changes focused and minimal.
- Follow existing package boundaries:
  - third_party/quickjs is vendored upstream source and should not be edited.
  - internal/cgo is the only C-Go interop boundary.
  - internal/* maps and manages ownership/lifecycle.
  - quickjs/* is public API only.
- Use gofmt and keep CI green.

## Pull Requests

- Open a focused PR with a clear description.
- Include tests for behavior changes.
- Mention any platform/compiler assumptions.

## Reporting Bugs

Use the GitHub bug report issue template and include:

- Go version
- OS and architecture
- Compiler details (for cgo)
- Reproduction steps
