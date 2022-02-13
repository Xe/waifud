{

  inputs = {
    naersk.url = "github:nmattia/naersk/master";
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, utils, naersk }:
    utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit system; };
        naersk-lib = pkgs.callPackage naersk { };
      in rec {
        packages = {
          waifud = naersk-lib.buildPackage {
            src = ./.;
            buildInputs = with pkgs; [
              pkg-config
              openssl
              sqliteInteractive
              libvirt
            ];
          };

          waifuctl = pkgs.stdenv.mkDerivation {
            src = self.packages."${system}".waifud;
            pname = "waifuctl";
            version = self.packages."${system}".waifud.version;
            phases = "installPhase";
            installPhase = ''
              mkdir -p $out/bin
              cp $src/bin/waifuctl $out/bin
            '';
          };
        };

        defaultPackage = self.packages."${system}".waifud;

        apps = {
          waifud = utils.lib.mkApp { drv = self.packages."${system}".waifud; };
          waifuctl =
            utils.lib.mkApp { drv = self.packages."${system}".waifuctl; };
        };

        defaultApp = self.apps."${system}".waifud;

        devShell = with pkgs;
          mkShell {
            buildInputs = [
              cargo
              rustc
              rustfmt
              pre-commit
              rustPackages.clippy
              openssl
              pkg-config
              sqliteInteractive
              libvirt
              dhall
              dhall-json
              go
              goimports
              gopls
              cdrkit
              jq
              jo
            ];
            RUST_SRC_PATH = rustPlatform.rustLibSrc;
          };
      });

}
