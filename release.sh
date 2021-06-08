#!/usr/bin/env bash

set -e

_version=$(grep -A 1 version info.plist | tail -n 1 | perl -pe 's/.*(\d+\.\d+\.\d+).*/v\1/')

if [ -z "${_version}" ]; then
  echo "Could not detect version"
  exit 1
fi

_gooses=(darwin)
_goarches=(amd64 arm64)

echo "Prep to build for OSes ${_gooses[*]} with arches ${_goarches[*]}"

for _goos in "${_gooses[@]}"; do
  for _goarch in "${_goarches[@]}"; do
    echo -n "Building OS ${_goos} arch ${_goarch}... "
    GOOS=${_goos} GOARCH=${_goarch} go build -o run ./app
    echo "built"

    _zipname="Safari_History_${_version}_${_goos}_${_goarch}.alfredworkflow"

    echo -n "Packing files into ${_zipname}... "
    zip -rq "${_zipname}" icon.png info.plist preflight.sh run
    echo "packed"
  done
done
