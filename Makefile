.PHONY: build install dev test clean release-dry

# Build the binary locally
build:
	go build -ldflags "-s -w -X github.com/ecoker/launchpad/internal/cli.version=dev" -o bin/launchpad ./cmd/launchpad

# Install to GOPATH/bin
install:
	go install -ldflags "-s -w -X github.com/ecoker/launchpad/internal/cli.version=dev" ./cmd/launchpad

# Quick dev run — pass args after --
# Usage: make dev ARGS="init ./my-app"
dev:
	go run ./cmd/launchpad $(ARGS)

# Run tests
test:
	go test ./...

# Clean build artifacts
clean:
	rm -rf bin/ dist/

# Dry-run GoReleaser (no publish)
release-dry:
	goreleaser release --snapshot --clean
