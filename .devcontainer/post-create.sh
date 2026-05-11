#!/usr/bin/env bash
set -euo pipefail

git config --global --add safe.directory /workspace || true
git config core.filemode false || true

go mod download

air -v
swag --version
migrate -version

echo "Dev container is ready."
