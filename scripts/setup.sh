#!/usr/bin/env bash
set -euo pipefail
echo "Running go mod tidy..."
go mod tidy
echo "Running unit tests..."
go test ./... -v
echo "Done."
