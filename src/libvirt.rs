use mac_address::MacAddress;
use rand::Rng;
use serde::{Deserialize, Serialize};

#[derive(Deserialize, Serialize, Clone, Debug)]
pub struct NewInstance {
    pub name: Option<String>,
    pub memory_mb: Option<i32>,
    pub cpus: Option<i32>,
    pub host: String,
    pub disk_size_gb: Option<i32>,
    pub zvol_prefix: Option<String>,
    pub distro: String,
    pub sata: Option<bool>,
    pub user_data: Option<String>,
}

pub fn random_mac() -> String {
    let mut addr = rand::thread_rng().gen::<[u8; 6]>();
    addr[0] = (addr[0] | 2) & 0xfe;
    MacAddress::new(addr).to_string()
}
