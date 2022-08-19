{ lib, ... }:

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

  systemd.services.cloud-init.requires = lib.mkForce [ "network.target" ];

  services.tailscale.enable = true;
  services.openssh.enable = true;

  services.cloud-init = {
    enable = true;
    ext4.enable = true;
  };
}
