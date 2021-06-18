{ pkgs ? import <nixpkgs> { } }:

let
  sources = import ./nix/sources.nix;
  gcss = pkgs.callPackage sources.gruvbox-css { };
in pkgs.mkShell {
  buildInputs = with pkgs; [
    # rust
    rustc
    rustfmt
    rust-analyzer
    cargo
    cargo-watch
    pkg-config
    sqliteInteractive
    diesel-cli

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

  # shellHook = ''
  #   rm ./public/static/gruvbox.css
  #   ln -s ${gcss}/gruvbox.css ./public/static/gruvbox.css
  #   rm ./public/static/alpine.js
  #   ln -s ${sources.alpinejs} ./public/static/alpine.js
  # '';
}
