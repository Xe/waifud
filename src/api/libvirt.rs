use crate::{Error, Result};
use axum::Json;
use serde::{Deserialize, Serialize};
use std::convert::TryFrom;
use virt::{connect::Connect, domain::Domain};

#[derive(Default, Debug, Clone, PartialEq, Serialize, Deserialize)]
pub struct Machine {
    pub name: String,
    pub host: String,
    pub active: bool,
    pub uuid: String,
    pub addr: Option<String>,
    pub memory_megs: u64,
    pub cpus: u32,
}

impl TryFrom<Domain> for Machine {
    type Error = Error;

    fn try_from(dom: Domain) -> Result<Self, Self::Error> {
        let addr: Option<String> = if dom.is_active()? {
            let mut addr: Vec<String> = dom
                .interface_addresses(virt::domain::VIR_DOMAIN_INTERFACE_ADDRESSES_SRC_LEASE, 0)?
                .into_iter()
                .map(|iface| iface.addrs.clone())
                .filter(|addrs| addrs.get(0).is_some())
                .map(|addrs| addrs.get(0).unwrap().clone().addr)
                .collect();

            if addr.get(0).is_none() {
                None
            } else {
                Some(addr.swap_remove(0))
            }
        } else {
            None
        };

        let info = dom.get_info()?;
        let max_memory = dom.get_max_memory()? / 1024;
        let conn = dom.get_connect()?;
        let host = conn.get_hostname()?;

        Ok(Machine {
            name: dom.get_name()?,
            host,
            active: dom.is_active()?,
            uuid: dom.get_uuid_string()?,
            addr,
            memory_megs: max_memory,
            cpus: info.nr_virt_cpu,
        })
    }
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
        result.push(Machine::try_from(dom)?);
    }
    conn.close()?;

    Ok(result)
}
