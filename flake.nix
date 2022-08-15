{
  inputs = {
    naersk.url = "github:nmattia/naersk/master";
    naersk.inputs.nixpkgs.follows = "nixpkgs";
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
          unique-monster = pkgs.stdenv.mkDerivation {
            src = self.packages."${system}".waifud;
            pname = "unique-monster";
            version = self.packages."${system}".waifud.version;
            phases = "installPhase";
            installPhase = ''
              mkdir -p $out/bin
              cp $src/bin/unique-monster $out/bin
            '';
          };

          waifud = naersk-lib.buildPackage {
            src = ./.;
            buildInputs = with pkgs; [
              pkg-config
              openssl
              sqlite-interactive
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
              mkdir -p $out/share/man/man1
              HOME=. $out/bin/waifuctl utils manpage $out/share/man/man1
              gzip -r $out/share/man/man1
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
            environment.defaultPackages =
              [ self.packages."${system}".waifuctl ];
          };

          waifud-common = { lib, ... }: {
            users.groups.waifud = lib.mkDefault { };

            users.users.waifud = {
              createHome = true;
              description = "waifud user";
              isSystemUser = true;
              group = "waifud";
              home = "/var/lib/waifud";
            };
          };

          waifud-host = { lib, pkgs, config, ... }:
            with lib;
            let cfg = config.xeserv.waifud;
            in {
              imports = [
                self.nixosModules."${system}".waifud-common
                self.nixosModules."${system}".waifuctl
              ];

              config = {
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

                  waifud-ssh-loadkey = {
                    wantedBy = [ "multi-user.target" ];
                    after = [ "waifud-ssh-agent" ];

                    environment.SSH_AUTH_SOCK = "/var/lib/waifud/agent.sock";
                    unitConfig.ConditionPathExists =
                      "/var/lib/waifud/id_ed25519";
                    serviceConfig = {
                      User = "waifud";
                      Group = "waifud";
                      Restart = "always";
                      WorkingDirectory = "/var/lib/waifud";
                      ExecStart =
                        "${pkgs.openssh}/bin/ssh-add /var/lib/waifud/id_ed25519";
                    };
                  };

                  waifud = {
                    wantedBy = [ "multi-user.target" ];
                    after = [ "waifud-ssh-agent" ];

                    environment = {
                      RUST_LOG = "tower_http=debug,waifud=debug,info";
                      SSH_AUTH_SOCK = "/var/lib/waifud/agent.sock";
                    };
                    unitConfig.ConditionPathExists =
                      "/var/lib/waifud/config.dhall";
                    serviceConfig = {
                      User = "waifud";
                      Group = "waifud";
                      Restart = "always";
                      WorkingDirectory = "/var/lib/waifud";
                      RestartSec = "30s";
                      ExecStartPre =
                        "ln --symbolic --force ${cfgDhall} ./config.dhall";
                      ExecStart =
                        "${self.packages."${system}".waifud}/bin/waifud";
                    };
                  };
                };
              };
            };

          waifud-runner = { pkgs, lib, config, ... }:
            with lib;
            let cfg = config.xeserv.waifud.runner;
            in {
              imports = [ self.nixosModules."${system}".waifud-common ];

              options.xeserv.waifud.runner = with lib; {
                parentDataset = mkOption {
                  type = types.str;
                  default = "rpool/local/vms";
                  description =
                    "the parent dataset to grant the waifud group zfs management access on";
                };

                sshKeys = mkOption {
                  type = with types; listOf str;
                  default = [ ];
                  description =
                    "the list of SSH public keys to allow waifud to ssh in as";
                };
              };

              config = {
                environment.defaultPackages = with pkgs; [ qemu zfs wget ];
                virtualisation.libvirtd.enable = lib.mkDefault true;

                systemd.services.waifud-runner-setup = {
                  wantedBy = [ "multi-user.target" ];
                  serviceConfig.Type = "oneshot";
                  script = ''
                    /run/current-system/sw/bin/zfs allow -g waifud create,destroy,mount,snapshot,rollback ${cfg.parentDataset}
                  '';
                };

                security.polkit.extraConfig = ''
                  /* Allow users in the waifud group to manage the libvirt daemon without authentication */
                  polkit.addRule(function(action, subject) {
                      if (action.id == "org.libvirt.unix.manage" && subject.isInGroup("waifud")) {
                              return polkit.Result.YES;
                      }
                  });
                '';

                users.users.waifud.openssh.authorizedKeys.keys = cfg.sshKeys;

                security.sudo.extraRules = [{
                  groups = [ "waifud" ];
                  users = [ "waifud" ];
                  runAs = "root:root";
                  commands = [{
                    command = "/run/current-system/sw/bin/qemu-img";
                    options = [ "NOPASSWD" ];
                  }];
                }];
              };
            };
        };

        devShell = with pkgs;
          mkShell {
            buildInputs = [
              cargo
              cargo-watch
              rustc
              rustfmt
              rust-analyzer
              pre-commit
              rustPackages.clippy
              openssl
              pkg-config
              sqlite-interactive
              libvirt
              dhall
              dhall-json
              jq
              jo
            ];
            DATABASE_URL = "./var/waifud.db";
            RUST_LOG = "tower_http=trace,debug";
            RUST_SRC_PATH = rustPlatform.rustLibSrc;
          };
      });
}
