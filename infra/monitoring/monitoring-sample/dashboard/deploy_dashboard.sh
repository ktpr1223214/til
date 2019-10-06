#!/usr/bin/env bash

set -euo pipefail

source "params.sh"
source "functions.sh"

# usage
function usage() {
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

# global
DASHBOARD_FILE=
DRY_RUN=

# 引数のパース
function parse_args() {
  while getopts ":Dh" o; do
    case ${o} in
      D)
        DRY_RUN="true"
        ;;
      h)
        usage
        exit 0
        ;;
      ?)
        usage
        ;;
    esac
  done

  # 処理したオプションの個数 shift
  shift $((OPTIND - 1))

  DASHBOARD_FILE=${1:-}
  # https://qiita.com/bsdhack/items/597eb7daee4a8b3276ba
  if [[ "${DASHBOARD_FILE}" == "" ]]; then
    usage
  fi
}

function create_dashboard() {


  body=$(echo "${dashboard}" | jq -c --arg uid "${uid}" --arg folder "${folder}" --arg folderId "${folderId}" --arg titleFolderId "${dashboard_folder}" --arg baseDashboardName "${base_dashboard_name}" --arg username "${username}" '
  {
    dashboard: .,
    folderId: $folderId | tonumber,
    overwrite: true
  } * {
    dashboard: {
      uid: $uid,
      editable: true,
      title: "TESTING \($username) \($titleFolderId) \($baseDashboardName): \(.title)",
      tags: ["playground"]
    }
  }
  ')
}

function main() {
  local dry_run="$1"
  local dashboard_file="$2"

  relative=${dashboard_file#"./"}
  folder="playground-for-testing"
  dashboard_folder="$(dirname "${relative}")"
  base_dashboard_name="$(basename "${dashboard_file}" | sed -e 's/\..*//')"

  username="${USER}"
  uid="testing-$(dirname "${relative}")-${username}-$(basename "${dashboard_file}" | sed -e 's/\..*//')"
  extension="${relative##*.}"

  folder_id=$(resolve_folder_id "${folder}")
  echo ${folder_id}

  if [[ -n ${dry_run} ]]; then
    echo "Running in dry run mode, would create ${dashboard_file} in folder ${folder} with uid ${uid}"
    exit 0
  fi
}

# エントリー処理
# https://www.m3tech.blog/entry/2018/08/21/bash-scripting
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
  parse_args "$@"
  main "${DRY_RUN}" "${DASHBOARD_FILE}"
fi