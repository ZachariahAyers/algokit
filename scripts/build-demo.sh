#!/usr/bin/env bash
set -euo pipefail

repo_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
go_command="${GO_COMMAND:-go}"
go_root="$($go_command env GOROOT)"
wasm_support="$go_root/lib/wasm/wasm_exec.js"
if [[ ! -f "$wasm_support" ]]; then
  wasm_support="$go_root/misc/wasm/wasm_exec.js"
fi
if [[ ! -f "$wasm_support" ]]; then
  echo "wasm_exec.js was not found under $go_root" >&2
  exit 1
fi

output_dir="$repo_dir/demo/dist"
mkdir -p "$output_dir"
cp "$repo_dir/demo/index.html" "$repo_dir/demo/styles.css" "$repo_dir/demo/runtime.js" \
  "$repo_dir/demo/app.js" "$output_dir/"
cp "$wasm_support" "$output_dir/wasm_exec.js"
GOOS=js GOARCH=wasm "$go_command" build -trimpath -o "$output_dir/algokit.wasm" "$repo_dir/cmd/wasm"

echo "Browser demo built at $output_dir"
