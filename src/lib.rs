#[macro_use]
extern crate diesel;
#[macro_use]
extern crate tracing;

use serde::{Deserialize, Serialize};

pub const APPLICATION_NAME: &str = concat!(env!("CARGO_PKG_NAME"), "/", env!("CARGO_PKG_VERSION"),);

pub mod schema;

#[derive(Debug, Serialize, Deserialize)]
pub struct NetworkConfig {
    pub name: String,
    pub uuid: uuid::Uuid,
    pub iface: String,
    pub mac_address: String,
    pub domain: String,
    pub ipv4_addr_base: String,
    pub ipv4_mask: String,
    pub ipv4_dstart: String,
    pub ipv4_dend: String,
    pub ipv6: Option<(String, u16)>,
}

include!(concat!(env!("OUT_DIR"), "/templates.rs"));
