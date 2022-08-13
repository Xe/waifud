use crate::{models::Distro, Result};
use futures::future::join_all;

pub mod arch;
pub mod ubuntu;

pub async fn get_all() -> Result<Vec<Distro>> {
    let distro_scrapes = join_all(vec![
        ubuntu::scrape(("22.04", "jammy")),
        ubuntu::scrape(("20.04", "focal")),
        ubuntu::scrape(("18.04", "bionic")),
    ])
    .await;
    let mut result: Vec<Distro> = vec![arch::scrape().await?];

    for distro in distro_scrapes {
        result.push(distro?);
    }

    Ok(result)
}
