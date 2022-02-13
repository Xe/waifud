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
      in {
        packages = {
          waifud-full = naersk-lib.buildPackage {
            src = ./.;
            buildInputs = with pkgs; [
              pkg-config
              openssl
              sqliteInteractive
              libvirt
            ];
          };

          waifuctl = pkgs.stdenv.mkDerivation {
            pname = "waifuctl";
            version = self.packages."${system}".waifud-full.version;

            installPhase = ''
              mkdir -p $out/bin
              cp ${self.packages."${system}".waifud-full}/bin/waifuctl $out/bin
            '';
          };
        };

        defaultPackage = self.packages.waifud-full;
        defaultApp = utils.lib.mkApp { drv = self.defaultPackage."${system}"; };

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
