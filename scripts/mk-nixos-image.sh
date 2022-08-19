#!/usr/bin/env nix-shell
#! nix-shell -i bash -p nixos-generators -p qemu -p rsync

set -ex

DATE="$(date +%Y-%m-%d-%H-%M)"

NIX_PATH=nixpkgs=channel:nixos-unstable-small nixos-generate -f qcow -c ./nixos-image.nix -o ./nixos-unstable-within-${DATE}
NIX_PATH=nixpkgs=channel:nixos-22.05-small nixos-generate -f qcow -c ./nixos-image.nix -o ./nixos-22.05-within-${DATE}

qemu-img convert -c -O qcow2 ./nixos-unstable-within-${DATE}/nixos.qcow2 nixos-unstable-within-${DATE}.qcow2
qemu-img convert -c -O qcow2 ./nixos-22.05-within-${DATE}/nixos.qcow2 nixos-22.05-within-${DATE}.qcow2

sha256sum nixos-unstable-within-${DATE}.qcow2 > nixos-unstable-within-${DATE}.qcow2.sha256
sha256sum nixos-22.05-within-${DATE}.qcow2 > nixos-22.05-within-${DATE}.qcow2.sha256

rsync -avz --progress *.qcow2* lufta:/srv/http/xena.greedo.xeserv.us/pkg/nixos/

rm ./nixos-unstable-within-${DATE} ./nixos-22.05-within-${DATE}
