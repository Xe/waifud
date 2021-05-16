# mkvm

Makes a virtual machine from a bunch of templates using
[cloud-init](https://cloudinit.readthedocs.io/en/latest/) userdata to customize
them.

This is an experimental tool I made for testing [this Tailscale
PR](https://github.com/tailscale/tailscale/pull/1934). Normally I assume that
this kind of stuff is very unsupported, however in this case I want to make an
exception and mention this explicitly: this tool is an unsupported thing I made
only to scratch an itch I've been having. If you want this tool to be supported,
please [contact me](https://christine.website) and arrange for a license and
payment.

This has only been tested on and assumes that it is running on NixOS with a zfs
root filesystem. You may need to create the parent zfs volume for VM disks with
a command like this:

```console
$ sudo zfs create -o mountpoint=none rpool/mkvm-test
```

Do not be surprised if this misbehaves. You have been warned.
