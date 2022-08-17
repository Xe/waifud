#!/usr/bin/env nix-shell
#! nix-shell -p deno -i bash

deno bundle ./instance_details.tsx ../static/instance_detail.js
