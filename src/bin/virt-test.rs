use anyhow::Result;
use virt::connect::Connect;

fn list_all_vms(uri: &str) -> Result<()> {
    let mut conn = Connect::open(uri)?;

    for dom in conn.list_all_domains(0)? {
        let addr: Option<String> = if dom.is_active()? {
            let mut addr: Vec<String> = dom
                .interface_addresses(virt::domain::VIR_DOMAIN_INTERFACE_ADDRESSES_SRC_LEASE, 0)?
                .into_iter()
                .filter(|iface| iface.addrs.clone().get(0).is_some())
                .map(|iface| iface.addrs.clone().get(0).unwrap().clone().addr)
                .collect();

            if addr.get(0).is_none() {
                None
            } else {
                Some(addr.swap_remove(0))
            }
        } else {
            None
        };
        println!(
            "{:<30} {:<15}{:<6}{:<34} {:<6}",
            uri,
            dom.get_name()?,
            if dom.is_active()? { "on" } else { "off" },
            dom.get_uuid_string()?,
            addr.unwrap_or("".into())
        );
    }
    conn.close()?;

    Ok(())
}

fn main() {
    for host in &["kos-mos", "logos", "ontos", "pneuma"] {
        list_all_vms(&format!("qemu+ssh://root@{}/system", host)).unwrap()
    }
}
