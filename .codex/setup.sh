#!/bin/bash
set -euo pipefail

# Navigate to repository root
REPO_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$REPO_ROOT"

# Install frontend dependencies
if [ -d "ui" ]; then
  echo "Installing npm dependencies for UI..."
  (cd ui && npm install)
fi

# Download Go modules
if [ -d "api" ]; then
  echo "Downloading Go modules..."
  (cd api && go mod download)
fi

# Install Go tools used by the build
echo "Installing go-bindata and swag..."
go install github.com/go-bindata/go-bindata/go-bindata@latest
go install github.com/swaggo/swag/cmd/swag@latest
