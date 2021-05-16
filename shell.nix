{ pkgs ? import <nixpkgs> { } }:

pkgs.mkShell {
  buildInputs = with pkgs; [ dhall dhall-json go goimports gopls cdrkit ];
}
