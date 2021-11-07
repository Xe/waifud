let Distro =
      { Type =
          { name : Text
          , downloadURL : Text
          , sha256Sum : Text
          , minSize : Natural
          }
      , default = { name = "", downloadURL = "", sha256Sum = "", minSize = 5 }
      }

in  [ Distro::{
      , name = "alpine-3.13"
      , downloadURL =
          "https://xena.greedo.xeserv.us/pkg/alpine/img/alpine-3.13.5-cloud-init-within.qcow2"
      , sha256Sum =
          "a2665c16724e75899723e81d81126bd0254a876e5de286b0b21553734baec287"
      , minSize = 2
      }
    , Distro::{
      , name = "alpine-3.14"
      , downloadURL =
          "https://xena.greedo.xeserv.us/pkg/alpine/img/alpine-3.14.2-cloud-init-within.qcow2"
      , sha256Sum =
          "2d2c8da0ac771d8346fc621f26b748e4f892c5f883251415a9670f0d639d4bef"
      , minSize = 2
      }
    , Distro::{
      , name = "alpine-edge"
      , downloadURL =
          "https://xena.greedo.xeserv.us/pkg/alpine/img/alpine-edge-2021-05-18-cloud-init-within.qcow2"
      , sha256Sum =
          "b3bb15311c0bd3beffa1b554f022b75d3b7309b5fdf76fb146fe7c72b83b16d0"
      , minSize = 2
      }
    , Distro::{
      , name = "amazon-linux"
      , downloadURL =
          "https://cdn.amazonlinux.com/os-images/2.0.20210427.0/kvm/amzn2-kvm-2.0.20210427.0-x86_64.xfs.gpt.qcow2"
      , sha256Sum =
          "6ef9daef32cec69b2d0088626ec96410cd24afc504d57278bbf2f2ba2b7e529b"
      , minSize = 25
      }
    , Distro::{
      , name = "android-9"
      , downloadURL = "http://tailnology/vms/android-x86-9.0.qcow2"
      , sha256Sum =
          "4dd5362025ee8c925299b7e88c2a10dba3ea29182f9405386e14601f76364cd3"
      , minSize = 20
      }
    , Distro::{
      , name = "arch"
      , downloadURL =
          "https://mirror.pkgbuild.com/images/v20210815.31636/Arch-Linux-x86_64-cloudimg-20210815.31636.qcow2"
      , sha256Sum =
          "ded266a65c6327ec7fe2d15fc86408c68e3a662b759a4c0a905583e3d3e71816"
      , minSize = 2
      }
    , Distro::{
      , name = "centos-7"
      , downloadURL =
          "https://cloud.centos.org/centos/7/images/CentOS-7-x86_64-GenericCloud-2009.qcow2c"
      , sha256Sum =
          "7ba4513d96591496213a07bbe25e3eb643d61491924c8548a91815b420fd9827"
      , minSize = 8
      }
    , Distro::{
      , name = "centos-8"
      , downloadURL =
          "https://cloud.centos.org/centos/8/x86_64/images/CentOS-8-GenericCloud-8.3.2011-20201204.2.x86_64.qcow2"
      , sha256Sum =
          "7ec97062618dc0a7ebf211864abf63629da1f325578868579ee70c495bed3ba0"
      , minSize = 10
      }
    , Distro::{
      , name = "debian-9"
      , downloadURL =
          "https://cdimage.debian.org/cdimage/openstack/9.13.21-20210511/debian-9.13.21-20210511-openstack-amd64.qcow2"
      , sha256Sum =
          "0667a08e2d947b331aee068db4bbf3a703e03edaf5afa52e23d534adff44b62a"
      , minSize = 2
      }
    , Distro::{
      , name = "debian-10"
      , downloadURL =
          "https://cdimage.debian.org/images/cloud/buster/20210329-591/debian-10-generic-amd64-20210329-591.qcow2"
      , sha256Sum =
          "70c61956095870c4082103d1a7a1cb5925293f8405fc6cb348588ec97e8611b0"
      , minSize = 2
      }
    , Distro::{
      , name = "debian-11"
      , downloadURL =
          "https://cdimage.debian.org/images/cloud/bullseye/daily/20210515-638/debian-11-generic-amd64-daily-20210515-638.qcow2"
      , sha256Sum =
          "0e77c13bd5f15759916d1e60e4925d8a3307bcd80af373fa929cdf419b602694"
      , minSize = 2
      }
    , Distro::{
      , name = "fedora-34"
      , downloadURL =
          "https://download.fedoraproject.org/pub/fedora/linux/releases/34/Cloud/x86_64/images/Fedora-Cloud-Base-34-1.2.x86_64.qcow2"
      , sha256Sum =
          "b9b621b26725ba95442d9a56cbaa054784e0779a9522ec6eafff07c6e6f717ea"
      , minSize = 5
      }
    , Distro::{
      , name = "oracle-linux-7.9"
      , downloadURL =
          "https://yum.oracle.com/templates/OracleLinux/OL7/u9/x86_64/OL7U9_x86_64-olvm-b86.qcow2"
      , sha256Sum =
          "2ef4c10c0f6a0b17844742adc9ede7eb64a2c326e374068b7175f2ecbb1956fb"
      , minSize = 40
      }
    , Distro::{
      , name = "oracle-linux-8.4"
      , downloadURL =
          "https://yum.oracle.com/templates/OracleLinux/OL8/u4/x86_64/OL8U4_x86_64-olvm-b85.qcow2"
      , sha256Sum =
          "b86e1f1ea8fc904ed763a85ba12e9f12f4291c019c8435d0e4e6133392182b0b"
      , minSize = 40
      }
    , Distro::{
      , name = "opensuse-leap-15.1"
      , downloadURL =
          "https://download.opensuse.org/repositories/Cloud:/Images:/Leap_15.1/images/openSUSE-Leap-15.1-OpenStack.x86_64.qcow2"
      , sha256Sum =
          "3203e256dab5981ca3301408574b63bc522a69972fbe9850b65b54ff44a96e0a"
      , minSize = 10
      }
    , Distro::{
      , name = "opensuse-leap-15.2"
      , downloadURL =
          "https://download.opensuse.org/repositories/Cloud:/Images:/Leap_15.2/images/openSUSE-Leap-15.2-OpenStack.x86_64.qcow2"
      , sha256Sum =
          "a293bf8ca21d4a8c2c146f2c42327ad27032afc2e15f61e0f5c05be46613e991"
      , minSize = 10
      }
    , Distro::{
      , name = "opensuse-leap-15.3"
      , downloadURL =
          "http://mirror.its.dal.ca/opensuse/distribution/leap/15.3/appliances/openSUSE-Leap-15.3-JeOS.x86_64-OpenStack-Cloud.qcow2"
      , sha256Sum =
          "aab4dbf6f7e40d6f3605ee1545e2858349bbefb6b09b28a1369900537a1133a1"
      , minSize = 5
      }
    , Distro::{
      , name = "opensuse-tumbleweed"
      , downloadURL =
          "https://download.opensuse.org/tumbleweed/appliances/openSUSE-Tumbleweed-JeOS.x86_64-OpenStack-Cloud.qcow2"
      , sha256Sum =
          "79e610bba3ed116556608f031c06e4b9260e3be2b193ce1727914ba213afac3f"
      , minSize = 5
      }
    , Distro::{
      , name = "ubuntu-14.04"
      , downloadURL =
          "http://cloud-images.ubuntu.com/trusty/20191107/trusty-server-cloudimg-amd64-disk1.img"
      , sha256Sum =
          "3c4ad0defbe729dd3c16d2851d775575d1c5351c85734418d3b89bfdfd28ebd1"
      , minSize = 5
      }
    , Distro::{
      , name = "ubuntu-16.04"
      , downloadURL =
          "https://cloud-images.ubuntu.com/xenial/20210429/xenial-server-cloudimg-amd64-disk1.img"
      , sha256Sum =
          "50a21bc067c05e0c73bf5d8727ab61152340d93073b3dc32eff18b626f7d813b"
      , minSize = 5
      }
    , Distro::{
      , name = "ubuntu-18.04"
      , downloadURL =
          "http://cloud-images.ubuntu.com/bionic/20210720/bionic-server-cloudimg-amd64.img"
      , sha256Sum =
          "c0bd5923b8edef2b32610a70ca99d92faebfbb1d8784c80328a45f7768433256"
      , minSize = 5
      }
    , Distro::{
      , name = "ubuntu-19.10"
      , downloadURL =
          "https://cloud-images.ubuntu.com/releases/eoan/release/ubuntu-19.10-server-cloudimg-amd64.img"
      , sha256Sum =
          "f0b499f0a7c8b5ca90ad12aa0b11a3643d5d272de02fabfd799eecb6227ec456"
      , minSize = 5
      }
    , Distro::{
      , name = "ubuntu-20.04"
      , downloadURL =
          "https://cloud-images.ubuntu.com/focal/20210916/focal-server-cloudimg-amd64.img"
      , sha256Sum =
          "5e10fb23ecb10123496dd21934a0a5b53b03936d0ab59060c7a95bf60b19152d"
      , minSize = 5
      }
    , Distro::{
      , name = "ubuntu-20.10"
      , downloadURL =
          "https://cloud-images.ubuntu.com/groovy/20210604/groovy-server-cloudimg-amd64.img"
      , sha256Sum =
          "2196df5f153faf96443e5502bfdbcaa0baaefbaec614348fec344a241855b0ef"
      , minSize = 5
      }
    , Distro::{
      , name = "ubuntu-21.04"
      , downloadURL =
          "https://cloud-images.ubuntu.com/hirsute/20210909/hirsute-server-cloudimg-amd64.img"
      , sha256Sum =
          "259f5225fa45029210befad9f43f704bac8c27babcea1f97db30c54c14b98cb2"
      , minSize = 5
      }
    , Distro::{
      , name = "windows-server-2012-r2"
      , downloadURL = "http://localhost:42069/windows2012r2.qcow2"
      , sha256Sum = "xxxfakewindowsserver2012r2"
      , minSize = 13
      }
    , Distro::{
      , name = "windows-10-21h1"
      , downloadURL = "http://tailnology/vms/win10-21h1-prepped.qcow2"
      , sha256Sum =
          "6637c1fe58d5ce13429e94d315e5a1c3580d1f0e379378c0bd49c5d7c4531100"
      , minSize = 40
      }
    , Distro::{
      , name = "dragonflybsd-6.0"
      , downloadURL =
          "https://object-storage.public.mtl1.vexxhost.net/swift/v1/1dbafeefbd4f4c80864414a441e72dd2/bsd-cloud-image.org/images/dragonflybsd/6.0.0/dragonflybsd-6.0.0-ufs.qcow2"
      , sha256Sum =
          "839c75b33e3c3a18ed3792dfd7123c3d39e9183fac085cda67cf4f4133c292e3"
      , minSize = 4
      }
    , Distro::{
      , name = "freebsd-12"
      , downloadURL =
          "https://object-storage.public.mtl1.vexxhost.net/swift/v1/1dbafeefbd4f4c80864414a441e72dd2/bsd-cloud-image.org/images/freebsd/12.2/freebsd-12.2.qcow2"
      , sha256Sum =
          "3c7c7fafe5c389b9295dcaab7a71c47cc30ad6e79e3a0c9cb164933ad2fb9814"
      , minSize = 3
      }
    , Distro::{
      , name = "freebsd-13"
      , downloadURL =
          "https://object-storage.public.mtl1.vexxhost.net/swift/v1/1dbafeefbd4f4c80864414a441e72dd2/bsd-cloud-image.org/images/freebsd/13.0/freebsd-13.0-ufs.qcow2"
      , sha256Sum =
          "64d1d9a3aa4b0cf118c7338bf57ec62005a436a23d3f82499a786690275ee5ee"
      , minSize = 3
      }
    , Distro::{
      , name = "netbsd-9.1"
      , downloadURL =
          "https://object-storage.public.mtl1.vexxhost.net/swift/v1/1dbafeefbd4f4c80864414a441e72dd2/bsd-cloud-image.org/images/netbsd/9.2/netbsd-9.2.qcow2"
      , sha256Sum =
          "653552edeba7b70041bdd1877f30c39ad933d80f1de1cc8d8b3f4e3eaf849687"
      , minSize = 3
      }
    , Distro::{
      , name = "openbsd-6.9"
      , downloadURL =
          "https://object-storage.public.mtl1.vexxhost.net/swift/v1/1dbafeefbd4f4c80864414a441e72dd2/bsd-cloud-image.org/images/openbsd/6.9/openbsd-6.9.qcow2"
      , sha256Sum =
          "bb52acb3770ef5caf3e941a38a0e47b785c5f148d4998c1d2945fc45920336e2"
      , minSize = 4
      }
    ]
