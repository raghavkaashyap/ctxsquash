# ctxsquash

ctxsquash is a local-first CLI that converts a repository or folder into a clean Markdown context file for AI assistants, debugging, documentation, and code review.

It does not call external APIs, upload data, or require a network connection at runtime.

## Install

Download a prebuilt binary from the [latest GitHub release](https://github.com/raghavkaashyap/ctxsquash/releases/latest). This is the best option if you do not have Go installed. Prebuilt binaries are published when a version tag such as `v0.1.0` is pushed.

macOS Apple Silicon:

```bash
curl -LO https://github.com/raghavkaashyap/ctxsquash/releases/latest/download/ctxsquash_darwin_arm64.tar.gz
tar -xzf ctxsquash_darwin_arm64.tar.gz
chmod +x ctxsquash
sudo mv ctxsquash /usr/local/bin/
```

macOS Intel:

```bash
curl -LO https://github.com/raghavkaashyap/ctxsquash/releases/latest/download/ctxsquash_darwin_amd64.tar.gz
tar -xzf ctxsquash_darwin_amd64.tar.gz
chmod +x ctxsquash
sudo mv ctxsquash /usr/local/bin/
```

Linux x86_64:

```bash
curl -LO https://github.com/raghavkaashyap/ctxsquash/releases/latest/download/ctxsquash_linux_amd64.tar.gz
tar -xzf ctxsquash_linux_amd64.tar.gz
chmod +x ctxsquash
sudo mv ctxsquash /usr/local/bin/
```

Windows users can download `ctxsquash_windows_amd64.zip` from the latest release and place `ctxsquash.exe` somewhere on `PATH`.

If you have Go installed, you can also install from source:

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
ctxsquash . --max-file-size 262144
ctxsquash . --format json
ctxsquash . --tree-only
ctxsquash . --stdout
```

When `--output` is omitted, ctxsquash prints to stdout. When `--output` is provided, it writes the generated context file to that path unless `--stdout` is also set.

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

Skip files larger than 128 KiB:

```bash
ctxsquash . --max-file-size 131072 --stdout
```

Print machine-readable JSON:

```bash
ctxsquash . --format json --stdout
```

## Output

The default Markdown output includes:

- A deterministic project tree.
- File path headings.
- Text file contents in Markdown code fences.
- Language identifiers based on file extensions.

JSON output includes a `tree` array and a `files` array with each file path, language identifier, and content. With `--tree-only`, JSON output omits file contents.

Binary files and files larger than `--max-file-size` are skipped. The default max file size is 262144 bytes. Common generated directories such as `.git`, `node_modules`, `target`, `dist`, `build`, and `vendor` are skipped by default. Rules in the root `.gitignore` are also respected.

## Secret Warnings

ctxsquash prints warnings to stderr when included file contents look like private key headers or sensitive assignments such as API keys, tokens, secrets, or passwords. Warning messages include the file path, line number, and pattern type, but they do not print the matched value.

## Limitations

- Nested `.gitignore` files are not loaded yet.
- Secret warnings are pattern-based and may miss secrets or flag harmless test fixtures.
- The Markdown tree format is intentionally simple rather than a full graphical tree.

## Test

```bash
go test ./...
```
