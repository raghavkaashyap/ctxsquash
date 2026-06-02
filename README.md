# ctxsquash

ctxsquash is a local-first CLI that converts a repository or folder into a clean Markdown context file for AI assistants, debugging, documentation, and code review.

It does not call external APIs, upload data, or require a network connection at runtime.

## Install

Install with Go:

```bash
go install github.com/raghavkaashyap/ctxsquash/cmd/ctxsquash@latest
```

Make sure your Go binary directory is on `PATH`. For many setups, that means adding this to your shell profile:

```bash
export PATH="$PATH:$(go env GOPATH)/bin"
```

Or build from a local checkout:

```bash
go build -o bin/ctxsquash ./cmd/ctxsquash
```

Or run directly:

```bash
go run ./cmd/ctxsquash . --stdout
```

## Usage

```bash
ctxsquash .
ctxsquash . --output context.md
ctxsquash . --include go,java,py,md,yml,json
ctxsquash . --exclude node_modules,target,dist,build
ctxsquash . --tree-only
ctxsquash . --stdout
```

When `--output` is omitted, ctxsquash prints to stdout. When `--output` is provided, it writes the Markdown file to that path unless `--stdout` is also set.

## Examples

Create a context file for the current repository:

```bash
ctxsquash . --output context.md
```

Print only Go, Markdown, YAML, and JSON files:

```bash
ctxsquash . --include go,md,yml,json --stdout
```

Print only the project tree:

```bash
ctxsquash . --tree-only --stdout
```

## Output

The generated Markdown includes:

- A deterministic project tree.
- File path headings.
- Text file contents in Markdown code fences.
- Language identifiers based on file extensions.

Binary files are skipped. Common generated directories such as `.git`, `node_modules`, `target`, `dist`, `build`, and `vendor` are skipped by default.

## Limitations

- `.gitignore` parsing is not implemented in the MVP.
- File size limits are not implemented yet.
- Secret detection and redaction are not implemented yet.
- The tree format is intentionally simple rather than a full graphical tree.

## Test

```bash
go test ./...
```
