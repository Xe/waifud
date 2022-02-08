#!/usr/bin/env nix-shell
#! nix-shell -p jq -p curl -i bash

cat pokedex-hisui.json \
    | jq -r '.[].name' \
    | tr '[:upper:]' '[:lower:]' \
    | tr ' ' '-' \
    | sed 's/\.//g' \
    | sed 's/://g' \
    | jq --raw-input '.' \
    | jq -s > pokemon-hisui.json
