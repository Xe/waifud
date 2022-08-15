#!/usr/bin/env nix-shell
#! nix-shell -p jq -p curl -i bash

curl https://raw.githubusercontent.com/fanzeyi/pokemon.json/master/pokedex.json \
    | jq -r '.[].name.english' \
    | tr '[:upper:]' '[:lower:]' \
    | tr ' ' '-' \
    | sed 's/\.//g' \
    | sed 's/://g' \
    | jq --raw-input '.' \
    | jq -s > pokemon.json
