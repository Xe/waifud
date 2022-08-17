#!/usr/bin/env nix-shell
#! nix-shell -p deno -i bash

set -e
cd $(dirname $0)

export RUST_LOG=info

mkdir -p ./static/js
deno bundle ./instance_details.tsx ./static/js/instance_detail.js
deno bundle ./instance_create.tsx ./static/js/instance_create.js
