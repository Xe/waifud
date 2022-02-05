{ pkgs ? import <nixpkgs> { } }:

let
  sources = import ./nix/sources.nix;
  gcss = pkgs.callPackage sources.gruvbox-css { };
in pkgs.mkShell rec {
  buildInputs = with pkgs; [
    # rust
    rustc
    rustfmt
    rust-analyzer
    cargo
    cargo-watch
    openssl
    pkg-config
    sqliteInteractive
    diesel-cli
    libvirt

    # dhall
    dhall
    dhall-json

    # go
    go
    goimports
    gopls

    # other tools
    cdrkit
    jq
    redis
  ];

  DATABASE_URL = "./var/waifud.db";
  ROCKET_DATABASES = ''{ main = { url = "${DATABASE_URL}" } }'';

  shellHook = ''
    ln -s ${gcss}/gruvbox.css ./public/static/gruvbox.css ||:
    ln -s ${sources.alpinejs} ./public/static/alpine.js ||:
  '';
}
