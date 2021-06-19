use anyhow::Result;
use virt::connect::Connect;

fn list_all_vms(uri: &str) -> Result<()> {
    let mut conn = Connect::open(uri)?;
    for dom in conn.list_all_domains(0)? {
        // if dom.is_active().unwrap() {
        println!(
            "{:<30} {:<15}{:<6}{}",
            uri,
            dom.get_name()?,
            if dom.is_active()? { "on" } else { "off" },
            dom.get_uuid_string()?
        );
        // }
    }
    conn.close()?;

    Ok(())
}

fn main() {
    for host in &["kos-mos", "logos", "ontos", "pneuma"] {
        list_all_vms(&format!("qemu+ssh://root@{}/system", host)).unwrap()
    }
}
