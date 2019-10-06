#!/usr/bin/env bash

function call_grafana_api() {
  local response

  # 失敗時に response って何か入ることあるか？
  response=$(curl -H 'Expect:' --http1.1 --compressed --silent --fail \
    -H "Authorization: Bearer ${GRAFANA_API_TOKEN}" \
    -H "Accept: application/json" \
    -H "Content-Type: application/json" \
    "$@") || {
    echo >&2 "API call to $1 failed: ${response}: exit code $?"
    return 1
  }

  echo "${response}"
}

function resolve_folder_id() {
  call_grafana_api "${GRAFANA_URL}/api/folders/$1" | jq '.id'
}

function prepare() {
  if [[ ! -d "vendor" ]]; then
    echo >&2 "vendor directory not found, running bundler.sh to install dependencies..."
    "bundler.sh"
  fi
}

function jsonnet_compile() {
  jsonnet -J . -J vendor "$@"
}
