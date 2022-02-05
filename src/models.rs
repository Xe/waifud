use serde::{Deserialize, Serialize};
use uuid::Uuid;

pub struct Instance {
    pub uuid: Uuid,
    pub name: String,
    pub host: String,
    pub memory: i32,
    pub disk_size: i32,
}

pub struct CloudconfigSeed {
    pub uuid: Uuid,
    pub user_data: String,
}

#[derive(Deserialize, Serialize)]
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
