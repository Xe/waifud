{ pkgs ? import <nixpkgs> { } }:

let
  sources = import ./nix/sources.nix;
  gcss = pkgs.callPackage sources.gruvbox-css { };
in pkgs.mkShell {
  buildInputs = with pkgs; [
    # rust
    rustc
    cargo
    cargo-watch
    pkg-config
    sqliteInteractive

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

  shellHook = ''
    ln -s ${gcss}/gruvbox.css ./public/static/gruvbox.css
    ln -s ${sources.alpinejs} ./public/static/alpine.js
  '';
}
