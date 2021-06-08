#!/usr/bin/env bash

_version=$(grep -A 1 version info.plist | tail -n 1 | perl -pe 's/.*(\d+\.\d+\.\d+).*/v\1/')

_gooses=(darwin)
_goarches=(amd64)

for _goos in "${_gooses[@]}"; do
  for _goarch in "${_goarches[@]}"; do
    GOOS=${_goos} GOARCH=${_goarch} go build -o run ./app

    _zipname="Safari_History_${_version}_${_goos}_${_goarch}.alfredworkflow"
    zip -r "${_zipname}" icon.png info.plist preflight.sh run
  done
done
