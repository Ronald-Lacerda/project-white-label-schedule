#!/bin/bash
set -e

echo "==> Running database migrations..."
migrate \
  -path ./migrations \
  -database "${DATABASE_URL}" \
  up

echo "==> Migrations done. Starting API..."
exec ./schedule-api
