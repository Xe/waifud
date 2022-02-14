{ pkgs ? import <nixpkgs> { } }:

pkgs.mkShell rec {
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
    libvirt

    # dhall
    dhall
    dhall-json

    # other tools
    cdrkit
    jq
    jo
  ];

  DATABASE_URL = "./var/waifud.db";
}
