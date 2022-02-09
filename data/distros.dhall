let Distro =
      { Type =
          { name : Text
          , downloadURL : Text
          , sha256Sum : Text
          , minSize : Natural
          , format : Text
          }
      , default =
        { name = ""
        , downloadURL = ""
        , sha256Sum = ""
        , minSize = 5
        , format = "waifud://qcow2"
        }
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
          "4cd35cbede40b165d9d491efee12352abb57f47f007909fe048cdf9cac638b89"
      , minSize = 2
      }
    , Distro::{
      , name = "alpine-3.15"
      , downloadURL =
          "https://xena.greedo.xeserv.us/pkg/alpine/img/alpine-3.15.0-cloud-init-within.qcow2"
      , sha256Sum =
          "13033e43676fc9776fcf228eb6465b8cb3b01ebe552a3c5de87483e00ec7295a"
      , minSize = 2
      }
    , Distro::{
      , name = "amazon-linux"
      , downloadURL =
          "https://cdn.amazonlinux.com/os-images/2.0.20211223.0/kvm/amzn2-kvm-2.0.20211223.0-x86_64.xfs.gpt.qcow2"
      , sha256Sum =
          "093ee88e855e5f13490cbda4ba72c8fef0bec17ecbbede12dc549fe2fbac511c"
      , minSize = 25
      }
    , Distro::{
      , name = "arch"
      , downloadURL =
          "https://mirror.pkgbuild.com/images/v20220204.46656/Arch-Linux-x86_64-cloudimg-20220204.46656.qcow2"
      , sha256Sum =
          "cf936eaededf6aff595dc72c91a9fac9d0577c6269a677a95d4daf80da612c8c"
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
          "https://cloud.centos.org/centos/8/x86_64/images/CentOS-8-GenericCloud-8.4.2105-20210603.0.x86_64.qcow2"
      , sha256Sum =
          "3510fc7deb3e1939dbf3fe6f65a02ab1efcc763480bc352e4c06eca2e4f7c2a2"
      , minSize = 10
      }
    , Distro::{
      , name = "centos-stream-9"
      , downloadURL =
          "https://cloud.centos.org/centos/9-stream/x86_64/images/CentOS-Stream-GenericCloud-9-20220204.0.x86_64.qcow2"
      , sha256Sum =
          "6d561125c73c72f73f057c6a9de45af242b678336f3ff4fafdf0223af96a5f47"
      , minSize = 10
      }
    , Distro::{
      , name = "fedora-35"
      , downloadURL =
          "https://mirror.dst.ca/fedora/releases/35/Cloud/x86_64/images/Fedora-Cloud-Base-35-1.2.x86_64.qcow2"
      , sha256Sum =
          "fe84502779b3477284a8d4c86731f642ca10dd3984d2b5eccdf82630a9ca2de6"
      , minSize = 5
      }
    , Distro::{
      , name = "nixos-21.11"
      , downloadURL =
          "https://xena.greedo.xeserv.us/pkg/nixos/nixos-21.11-within-2022-02-08-08-37.qcow2"
      , sha256Sum =
          "1d7fb8a2dad803e52935e191421f502bbf1b611d22f011d62d02254fb36b627d"
      , minSize = 8
      }
    , Distro::{
      , name = "nixos-unstable"
      , downloadURL =
          "https://xena.greedo.xeserv.us/pkg/nixos/nixos-unstable-within-2022-02-08-08-37.qcow2"
      , sha256Sum =
          "2ba68d1f299f8152993adac809d52ba244d5442840f4ecbdb12f8a69fdcfc672"
      , minSize = 8
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
      , name = "opensuse-leap-15.4"
      , downloadURL =
          "http://mirror.its.dal.ca/opensuse/distribution/leap/15.4/appliances/openSUSE-Leap-15.4-JeOS.x86_64-15.4-OpenStack-Cloud-Build5.46.qcow2"
      , sha256Sum =
          "dd36668e2fc206d104ea43894f0a778f2312adee368dc2f6ef7738a7c8c4a686"
      , minSize = 5
      }
    , Distro::{
      , name = "opensuse-tumbleweed"
      , downloadURL =
          "https://download.opensuse.org/tumbleweed/appliances/openSUSE-Tumbleweed-JeOS.x86_64-15.1.0-OpenStack-Cloud-Snapshot20220204.qcow2"
      , sha256Sum =
          "d3ef72c231a02504732e60a3b3e016d743ecea4964c5b50ec2704e1238118e0a"
      , minSize = 5
      }
    , Distro::{
      , name = "rocky-linux-8"
      , downloadURL =
          "https://download.rockylinux.org/pub/rocky/8.5/images/Rocky-8-GenericCloud-8.5-20211114.2.x86_64.qcow2"
      , sha256Sum =
          "c23f58f26f73fb9ae92bfb4cf881993c23fdce1bbcfd2881a5831f90373ce0c8"
      , minSize = 10
      }
    , Distro::{
      , name = "ubuntu-18.04"
      , downloadURL =
          "http://cloud-images.ubuntu.com/bionic/20220131/bionic-server-cloudimg-amd64.img"
      , sha256Sum =
          "66f5c336b54b668fe64cb42b3b2f794327d64a376c414944234be51d61edeeec"
      , minSize = 5
      }
    , Distro::{
      , name = "ubuntu-20.04"
      , downloadURL =
          "http://cloud-images.ubuntu.com/focal/20220204/focal-server-cloudimg-amd64.img"
      , sha256Sum =
          "1d91ddc804dc201ba9a9f49def83fb4f40a76c3666e7daa50b38446a21b9543f"
      , minSize = 5
      }
    ]
