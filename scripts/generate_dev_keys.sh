#!/bin/bash
set -euo pipefail

ROOT="$(dirname "$0")/.."
KEY_DIR="$ROOT/dev/keys"
mkdir -p "$KEY_DIR"

"$ROOT/bin/standupctl" secrets gen-keys --dir "$KEY_DIR" "$@"