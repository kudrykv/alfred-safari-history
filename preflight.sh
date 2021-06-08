#!/usr/bin/env bash

set -e

function file_last_modified() {
  stat -f '%m' "$1"
}

function store_last_modified() {
  file_last_modified "$1" >"$2"
}

_wf_path=$(pwd)
_copy_file=$(pwd)/History.db
_copy_stat_file=${_wf_path}/History.db.stat

# yep, for some reason `cd ~/Library/Safari` doesn't work ¯\_(ツ)_/¯
# actually, any full-path operations don't work
cd ~/Library
cd Safari

_orig_file=./History.db
_orig_stat=$(file_last_modified "${_orig_file}")

if [ ! -e "${_copy_file}" ]; then
  cp "${_orig_file}" "${_copy_file}"
  store_last_modified "${_orig_file}" "${_copy_stat_file}"
else
  _copy_stat=$(cat "${_copy_stat_file}")

  if [ "${_orig_stat}" != "${_copy_stat}" ]; then
    cp "${_orig_file}" "${_copy_file}"
  fi
fi

cd "${_wf_path}"

./run "$1"