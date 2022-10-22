# isekaid Design

- Metadata service
- Runs on every machine
- Looks up instance information by remote IP
  - Listen to the DHCP client table
  - Look up mac address of instance
  - ask waifud for machine info by mac address
  - use that to construct metadata for cloud-init
