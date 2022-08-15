#!/usr/bin/env nix-shell
#! nix-shell -p jq -p curl -i bash

curl 'https://api.sibr.dev/chronicler/v2/entities?type=player&at=2020-11-01T00:00:00Z' \
   | jq '.items[].data.name' -r \
   | grep -v -- "-" \
   | tr '[:upper:]' '[:lower:]' \
   | tr ' ' '-' \
   | sed 's/\.//g' \
   | sed 's/://g' \
   | jq --raw-input '.' \
   | jq -s > blaseball.json
