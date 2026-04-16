#!/usr/bin/env bash

echo "==> Syncing Go modules"
go mod tidy

echo "==> Starting API server"
go run ./cmd/api
