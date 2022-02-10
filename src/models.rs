use serde::{Deserialize, Serialize};
use uuid::Uuid;

#[derive(Debug, Default, Deserialize, Serialize, Clone)]
pub struct Instance {
    pub uuid: Uuid,
    pub name: String,
    pub host: String,
    pub mac_address: String,
    pub memory: i32,
    pub disk_size: i32,
    pub zvol_name: String,
    pub status: String,
    pub distro: String,
}

pub struct CloudconfigSeed {
    pub uuid: Uuid,
    pub user_data: String,
}

#[derive(Deserialize, Serialize, Debug, Clone)]
pub struct Distro {
    pub name: String,
    #[serde(rename = "downloadURL")]
    pub download_url: String,
    #[serde(rename = "sha256Sum")]
    pub sha256sum: String,
    #[serde(rename = "minSize")]
    pub min_size: i32,
    pub format: String,
}
