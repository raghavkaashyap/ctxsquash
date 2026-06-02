.PHONY: test run build

test:
	go test ./...

run:
	go run ./cmd/ctxsquash . --stdout

build:
	go build -o bin/ctxsquash ./cmd/ctxsquash
