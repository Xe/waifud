#!/usr/bin/env nix-shell
#! nix-shell -i bash -p nixos-generators

set -ex

NIX_PATH=nixpkgs=channel:nixos-unstable nixos-generate -f qcow -c ./nixos-image.nix -o ./nixos-unstable-within-$(date +%Y-%m-%d-%H-%M).qcow2
NIX_PATH=nixpkgs=channel:nixos-21.11 nixos-generate -f qcow -c ./nixos-image.nix -o ./nixos-21.11-within-$(date +%Y-%m-%d-%H-%M).qcow2
