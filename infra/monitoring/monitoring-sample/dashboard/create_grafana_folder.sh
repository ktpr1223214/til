#!/usr/bin/env bash

#
# <uid> <title> の
# TODO: uid == title に決め打ちで問題あるか？

set -euo pipefail
IFS=$'\n\t'

source "params.sh"

usage() {
  cat <<-EOF
Usage $0 <uid> <title>

Description
指定した <uid> <title> で Grafana にフォルダを作成

GRAFANA_API_TOKEN を環境変数に設定しておく

Reference
https://grafana.com/docs/reference/dashboard_folders/
https://grafana.com/docs/http_api/folder/
EOF
}

function main() {
  local uid="$1"
  local title="$2"

  # 初期実行時には jq で失敗するようなので一時的に pipefail 無効化
  set +o pipefail
  # uid/title はいずれも一意である必要があるはずなので、uid のみで判定しておく
  # なので、現状 title の変更をしたいというのはこのスクリプトでは未対応
  current_title=$(curl --silent --fail \
    -H "Authorization: Bearer ${GRAFANA_API_TOKEN}" \
    -H "Accept: application/json" \
    -H "Content-Type: application/json" \
    "${GRAFANA_URL}/api/folders/${uid}" | jq .title)
  set -o pipefail

  if [[ -n ${current_title} ]]; then
    echo "uid: ${uid} already exists"
    exit 0
  fi

  curl -v \
    -H "Authorization: Bearer ${GRAFANA_API_TOKEN}" \
    -H "Accept: application/json" \
    -H "Content-Type: application/json" \
    "${GRAFANA_URL}/api/folders/" \
    -d @- <<EOF
{
  "uid": "${uid}",
  "title": "${title}"
}
EOF
}

# エントリー処理
# https://www.m3tech.blog/entry/2018/08/21/bash-scripting
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
  if [[ $# -ne 2 ]]; then
    usage
    exit 1
  fi

  main "$1" "$2"
fi