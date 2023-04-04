# waifud

![enbyware](https://pride-badges.pony.workers.dev/static/v1?label=enbyware&labelColor=%23555&stripeWidth=8&stripeColors=FCF434%2CFFFFFF%2C9C59D1%2C2C2C2C)
![made with Nix](https://img.shields.io/badge/made%20with-Nix-blue?logo=nixos)
![built with Garnix](https://img.shields.io/static/v1?label=Built%20with&message=Garnix&color=blue&style=flat&logo=nixos&link=https://garnix.io&labelColor=111212)
![license](https://img.shields.io/github/license/Xe/waifud)
![language count](https://img.shields.io/github/languages/count/Xe/waifud)
![repo size](https://img.shields.io/github/repo-size/Xe/waifud)

A few tools to help me manage and run virtual machines across a homelab cluster.

waifud was made for my own personal use and I do not expect it to be very useful
outside that context. If you do want to run this on your
infrastructure anyways, please [contact me](https://xeiaso.net/contact).

<big>THIS IS EXPERIMENTAL! USE IT AT YOUR OWN PERIL!</big>

TODO(Xe): Link to blogpost on the design/implementation once it is a thing.

Blogposts about waifud:
 - [waifud Plans](https://xeiaso.net/blog/waifud-plans-2021-06-19)
 - [waifud Progress Report #1](https://xeiaso.net/blog/waifud-progress-2022-02-06)
 - [waifud Progress Report #2](https://xeiaso.net/blog/waifud-progress-report-2)

Overall architecture diagram (with incomplete components marked with a
clock):

```mermaid
flowchart TD
    subgraph control plane
    WD[fa:fa-rust waifud]
    WC[fa:fa-rust waifuctl]
    ID[fa:fa-golang fa:fa-clock isekaid]
    MD[fa:fa-golang fa:fa-clock megamid]
    PD[fa:fa-golang fa:fa-clock portald]
    end
    subgraph VM plane
    LV[fa:fa-c libvirt]
    WH[fa:fa-linux runner\nnodes]
    VM[fa:fa-linux virtual\nmachines]
    end
    subgraph external
    TS[fa:fa-golang Tailscale]
    end

    PD --> |tailnet ingress for| WD
    WC --> |operator tool for| WD
    WC --> |usually connects via|PD
    ID --> |fetches node metadata\nand secrets for| WD
    VM --> |cloud-init\nmetadata| ID
    WD --> |manages libvirt on| WH
    LV --> |actually runs VMs| VM
    VM --> |network storage| MD
    WD --> |sets limits for\nrequests metrics from| MD  
    WH --> |runs| LV
    WH <--> |subnet router\ninterconnect| TS
    TS --> |network layer for| PD
    VM --> |usually a part of| TS
```
