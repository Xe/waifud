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
    pub tailscale: Tailscale,
    #[serde(rename = "pasetoKeypair")]
    pub paseto_keypair: String,
    pub yubikey: Yubikey,
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

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Yubikey {
    #[serde(rename = "clientID")]
    pub client_id: String,
    pub key: String,
}
