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
      , name = "arch"
      , downloadURL =
          "https://mirror.pkgbuild.com/images/v20210515.22945/Arch-Linux-x86_64-cloudimg-20210515.22945.qcow2"
      , sha256Sum =
          "e4077f5ba3c5d545478f64834bc4852f9f7a2e05950fce8ecd0df84193162a27"
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
          "https://object-storage.public.mtl1.vexxhost.net/swift/v1/1dbafeefbd4f4c80864414a441e72dd2/bsd-cloud-image.org/images/freebsd/13.0/freebsd-13.0.qcow2"
      , sha256Sum =
          "f7d16fb927f836f94cda37955314506e8507476d0b2d985acf572f1c7ce90e6a"
      , minSize = 3
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
          "22e0392e4d0becb523d1bc5f709366140b7ee20d6faf26de3d0f9046d1ee15d5"
      , minSize = 5
      }
    , Distro::{
      , name = "opensuse-tumbleweed"
      , downloadURL =
          "https://download.opensuse.org/tumbleweed/appliances/openSUSE-Tumbleweed-JeOS.x86_64-OpenStack-Cloud.qcow2"
      , sha256Sum =
          "8bc3aafa7bfb6b7ab0f69a40dea05827b13faca4169b755af44edfeed3742372"
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
          "https://cloud-images.ubuntu.com/bionic/20210604/bionic-server-cloudimg-amd64.img"
      , sha256Sum =
          "50c38d3f7307fe770c15a69b316d0001ac28e484239218d23e1ca8c8e7ec9a10"
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
          "https://cloud-images.ubuntu.com/focal/20210603/focal-server-cloudimg-amd64.img"
      , sha256Sum =
          "1c0969323b058ba8b91fec245527069c2f0502fc119b9138b213b6bfebd965cb"
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
          "https://cloud-images.ubuntu.com/hirsute/20210603/hirsute-server-cloudimg-amd64.img"
      , sha256Sum =
          "bf07f36fc99ff521d3426e7d257e28f0c81feebc9780b0c4f4e25ae594ff4d3b"
      , minSize = 5
      }
    , Distro::{
      , name = "windows-server-2012-r2"
      , downloadURL = "http://localhost:42069/windows2012r2.qcow2"
      , sha256Sum = "xxxfakewindowsserver2012r2"
      , minSize = 13
      }
    ]
