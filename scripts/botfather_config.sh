#!/bin/bash
set -euo pipefail

cd "$(dirname "$0")/.."
go run ./internal/bot/commands/cmd