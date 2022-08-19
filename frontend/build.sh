#!/usr/bin/env nix-shell
#! nix-shell -p deno -i bash

set -e
cd $(dirname $0)

DENO_FLAGS=''

if [ "$1" == "--dev" ]; then
	DENO_FLAGS='--watch'
fi

export RUST_LOG=info

mkdir -p ./static/js
deno bundle $DENO_FLAGS ./instance_detail.tsx ./static/js/instance_detail.js &
deno bundle $DENO_FLAGS ./instance_create.tsx ./static/js/instance_create.js &

wait
