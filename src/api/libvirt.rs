use crate::Result;
use axum::Json;
use serde::{Deserialize, Serialize};
use virt::connect::Connect;

#[derive(Default, Debug, Clone, PartialEq, Serialize, Deserialize)]
pub struct Machine {
    pub name: String,
    pub host: String,
    pub active: bool,
    pub uuid: String,
    pub addr: Option<String>,
    pub memory_megs: u64,
}

#[instrument(err)]
pub async fn get_machines() -> Result<Json<Vec<Machine>>> {
    let mut result = Vec::new();
    for host in &["kos-mos", "logos", "ontos", "pneuma"] {
        result.extend_from_slice(&list_all_vms(
            &format!("qemu+ssh://root@{}/system", host),
            host.to_string(),
        )?);
    }
    Ok(Json(result))
}

#[instrument(skip(uri), err)]
fn list_all_vms(uri: &str, host: String) -> Result<Vec<Machine>> {
    debug!("connecting to {}: {}", host, uri);
    let mut conn = Connect::open(uri)?;
    let mut result = vec![];

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

        let max_memory = dom.get_max_memory()? / 1024;

        result.push(Machine {
            name: dom.get_name()?,
            host: host.clone(),
            active: dom.is_active()?,
            uuid: dom.get_uuid_string()?,
            addr,
            memory_megs: max_memory,
        });
    }
    conn.close()?;

    Ok(result)
}
