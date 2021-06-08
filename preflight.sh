#!/usr/bin/env bash

set -e

_wf_path=$(pwd)

# yep, for some reason `cd ~/Library/Safari` doesn't work ¯\_(ツ)_/¯
# actually, any full-path operations don't work
# anyway, once we're here, we can access `History.db`
cd ~/Library
cd Safari

"${_wf_path}"/run "$1"
