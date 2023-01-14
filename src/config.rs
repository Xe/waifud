use serde::{Deserialize, Serialize};
use std::{fmt, net::IpAddr};

#[derive(Clone, Serialize, Deserialize)]
pub struct Config {
    #[serde(rename = "baseURL")]
    pub base_url: String,
    pub hosts: Vec<String>,
    #[serde(rename = "bindHost")]
    pub bind_host: IpAddr,
    pub port: u16,
    #[serde(rename = "rpoolBase")]
    pub rpool_base: String,
    #[serde(rename = "qemuPath")]
    pub qemu_path: String,
    #[serde(skip_serializing)]
    pub tailscale: Tailscale,
}

impl fmt::Debug for Config {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        write!(f, "Config()")
    }
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Tailscale {
    #[serde(rename = "apiKey")]
    pub api_key: String,
    pub tailnet: String,
}
