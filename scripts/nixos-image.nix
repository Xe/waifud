{ lib, pkgs, ... }:

{
  boot.initrd.availableKernelModules =
    [ "ata_piix" "uhci_hcd" "virtio_pci" "sr_mod" "virtio_blk" ];
  boot.initrd.kernelModules = [ ];
  boot.kernelModules = [ ];
  boot.extraModulePackages = [ ];
  boot.growPartition = true;
  boot.kernelParams = [ "console=ttyS0" ];
  boot.loader.grub.device = "/dev/vda";
  boot.loader.timeout = 0;

  fileSystems."/" = {
    device = "/dev/disk/by-label/nixos";
    fsType = "ext4";
    autoResize = true;
  };

  nix = {
    package = pkgs.nixVersions.stable;
    extraOptions = ''
      experimental-features = nix-command flakes
    '';

    settings = {
      auto-optimise-store = true;
      sandbox = true;
      substituters = [
        "https://xe.cachix.org"
        "https://nix-community.cachix.org"
        "https://cuda-maintainers.cachix.org"
        "https://cache.floxdev.com?trusted=1"
        "https://cache.garnix.io"
      ];
      trusted-users = [ "root" "cadey" ];
      trusted-public-keys = [
        "xe.cachix.org-1:kT/2G09KzMvQf64WrPBDcNWTKsA79h7+y2Fn2N7Xk2Y="
        "nix-community.cachix.org-1:mB9FSh9qf2dCimDSUo8Zy7bkq5CX+/rkCWyvRCYg3Fs="
        "cuda-maintainers.cachix.org-1:0dq3bujKpuEPMCX6U4WylrUDZ9JyUG0VpVZa7CNfq5E="
        "flox-store-public-0:8c/B+kjIaQ+BloCmNkRUKwaVPFWkriSAd0JJvuDu4F0="
        "cache.garnix.io:CTFPyKSLcx5RMJKfLo5EEPUObbA78b0YQ2DTCJXqr9g="
      ];
    };
  };

  systemd.services."within.website-first-run" = {
    description = "bootstrap the first run of a NixOS machine on waifud";
    wantedBy = [ "multi-user.target" ];
    after = [ "network.target" "polkit.service" ];
    path = [ "/run/current-system/sw/" ];

    script = with pkgs; ''
      if ! [ -f /etc/nixos/configuration.nix ]; then
        install -D ${./nixos-image.nix} /mnt/etc/nixos/configuration.nix
      fi
    '';
  };

  systemd.services.cloud-init.requires = lib.mkForce [ "network.target" ];

  services.tailscale.enable = true;
  services.openssh.enable = true;

  services.cloud-init = {
    enable = true;
    ext4.enable = true;
  };

  users.motd = "Welcome to waifud <3";
}
