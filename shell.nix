{ pkgs ? import <nixpkgs> {} }:

pkgs.mkShell {
  buildInputs = with pkgs; [
    dhall go goimports gopls cdrkit
  ];
}
