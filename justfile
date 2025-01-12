# CLI tools

# Help
help:
    @echo "Command line helpers for this project.\n"
    @just -l

# Run all pre-commit checks
all-checks:
	pre-commit run --all-files

# Fix imports
imports:
  goimports -w ./..

# Go fmt
fmt:
  go fmt ./...

# Snapshot
snapshot:
  goreleaser build --snapshot --single-target --clean -f .goreleaser.yml

# Snapshot (verboise)
snapshot-verbose:
  goreleaser build --verbose --snapshot --single-target --clean -f .goreleaser.yml

# Release
release:
  goreleaser release --skip=publish --clean -f .goreleaser.yml

# Single-target release
target:
  goreleaser build --single-target --clean -f .goreleaser.yml
