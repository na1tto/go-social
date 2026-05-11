#!/usr/bin/env bash
set -euo pipefail

sudo chown -R app:app /home/app/go /home/app/.cache/go-build || true

git config --global --add safe.directory /workspace

go mod download

air -v
swag --version
migrate -version

echo "Dev container ready."
