{ pkgs ? import <nixpkgs> { } }:

pkgs.mkShell {
  buildInputs = with pkgs; [
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
}
