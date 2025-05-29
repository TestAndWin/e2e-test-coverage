#!/bin/bash
set -euo pipefail

# Navigate to repository root
REPO_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$REPO_ROOT"

(cd ui && npm install)
cd ..
(cd api && go mod download)

echo "Installing go-bindata and swag..."
go install github.com/go-bindata/go-bindata/go-bindata@latest
go install github.com/swaggo/swag/cmd/swag@latest
