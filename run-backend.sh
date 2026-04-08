#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
GO_BIN="C:/Program Files/Go/bin/go"

cd "$ROOT_DIR/backend"
exec "$GO_BIN" run ./cmd/api
