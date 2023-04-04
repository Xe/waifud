use crate::models::Distro;
use serde::{Deserialize, Serialize};
use std::collections::HashMap;

#[derive(Deserialize, Serialize, Debug)]
pub struct Metadata {
    pub fname: String,
    pub sha256: String,
}

pub async fn scrape() -> crate::Result<Vec<Distro>> {
    let md: HashMap<String, Metadata> =
        reqwest::get("https://xena.greedo.xeserv.us/pkg/nixos/metadata.json")
            .await?
            .error_for_status()?
            .json()
            .await?;

    let mut result: Vec<Distro> = vec![];

    for (key, val) in md.iter() {
        result.push(Distro {
            name: format!("nixos-{key}"),
            download_url: format!("https://xena.greedo.xeserv.us/pkg/nixos/{}", val.fname),
            sha256sum: val.sha256.clone(),
            min_size: 8,
            format: "waifud://qcow2".to_string(),
        })
    }

    Ok(result)
}
