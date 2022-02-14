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

        defaultPackage = self.packages."${system}".waifuctl;

        apps = {
          waifud = utils.lib.mkApp { drv = self.packages."${system}".waifud; };
          waifuctl =
            utils.lib.mkApp { drv = self.packages."${system}".waifuctl; };
        };

        defaultApp = self.apps."${system}".waifuctl;

        nixosModules = {
          waifuctl = { ... }: {
            environment.systemPackages = [ self.packages."${system}".waifuctl ];
          };

          waifud-common = { lib, ... }: {
            users.groups.waifud = lib.mkDefault { };
          };

          waifud-host = { lib, pkgs, config, ... }:
            with lib;
            let
              cfg = config.xeserv.waifud;
              cfgJSON = pkgs.writeTextFile {
                name = "waifud.json";
                text = toJSON cfg;
              };
              cfgDhall = pkgs.stdenv.mkDerivation {
                name = "waifud-config-latest";
                src = cfgJSON;
                inputs = [ pkgs.dhall-json ];
                phases = "installPhase";
                installPhase = ''
                  cat $src | json-to-dhall --output $out
                '';
              };
            in {
              imports = [
                self.nixosModules."${system}".waifud-common
                self.nixosModules."${system}".waifuctl
              ];

              options = {
                baseURL = mkOption {
                  type = types.str;
                  default = "http://192.168.122.1:23818";
                  description =
                    "the base URL for VMs to poke waifud for config";
                };

                hosts = mkOption {
                  type = with types; listOf str;
                  default = [ "kos-mos" "logos" "ontos" "pneuma" ];
                  description =
                    "the list of hosts that waifud can manage (ip address/domain name)";
                };

                bindHost = mkOption {
                  type = types.str;
                  default = "::";
                  description = "the host to bind waifud on";
                };

                port = mkOption {
                  type = types.port;
                  default = 23818;
                  description = "the port that waifud should bind on";
                };
              };

              config = {
                users.users.waifud = {
                  createHome = true;
                  description = "waifud user";
                  isSystemUser = true;
                  group = "waifud";
                  home = "/var/lib/waifud";
                };

                systemd.services = {
                  waifud-ssh-agent = {
                    wantedBy = [ "multi-user.target" ];
                    serviceConfig = {
                      User = "waifud";
                      Group = "waifud";
                      Restart = "always";
                      WorkingDirectory = "/var/lib/waifud";
                      ExecStart =
                        "${pkgs.openssh}/bin/ssh-agent -D -a /var/lib/waifud/agent.sock";
                    };
                  };

                  waifud = {
                    wantedBy = [ "multi-user.target" ];

                    environment = {
                      RUST_LOG = "tower_http=debug,waifud=debug,info";
                      SSH_AUTH_SOCK = "/var/lib/waifud/agent.sock";
                    };
                    serviceConfig = {
                      User = "waifud";
                      Group = "waifud";
                      Restart = "always";
                      WorkingDirectory = "/var/lib/waifud";
                      RestartSec = "30s";
                      ExecStartPre = [
                        "rm -f ./config.dhall"
                        "cp ${cfgDhall} ./config.dhall"
                      ];
                      ExecStart =
                        "${self.packages."${system}".waifud}/bin/waifud";
                    };
                  };
                };
              };
            };

          waifud-runner = { pkgs, lib, ... }: {
            imports = [ self.nixosModules."${system}".waifud-common ];

            environment.systemPackages = with pkgs; [ qemu zfs wget ];
            services.libvirt.enable = lib.mkDefault true;

            security.polkit.extraConfig = ''
              /* Allow users in the waifud group to manage the libvirt daemon without authentication */
              polkit.addRule(function(action, subject) {
                  if (action.id == "org.libvirt.unix.manage" && subject.isInGroup("waifud")) {
                          return polkit.Result.YES;
                  }
              });
            '';
          };
        };

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
              jq
              jo
            ];
            RUST_SRC_PATH = rustPlatform.rustLibSrc;
          };
      });

}
