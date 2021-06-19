#!/usr/bin/env nix-shell
#! nix-shell -p jq -p curl -i bash

curl https://raw.githubusercontent.com/fanzeyi/pokemon.json/master/pokedex.json \
    | jq -r '.[].name.english' \
    | tr '[:upper:]' '[:lower:]' \
    | jq --raw-input '.' \
    | jq -s > pokemon.json
