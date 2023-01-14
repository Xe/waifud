use std::time::Duration;

use crate::{models::Distro, Result};
use futures::future::join_all;
use rusqlite::params;
use tokio::time::Instant;

pub mod amazon_linux;
pub mod arch;
pub mod rocky_linux;
pub mod ubuntu;

pub async fn get_all() -> Result<Vec<Distro>> {
    let ubuntus = join_all(vec![
        ubuntu::scrape(("22.04", "jammy")),
        ubuntu::scrape(("20.04", "focal")),
        ubuntu::scrape(("18.04", "bionic")),
    ])
    .await;
    let mut result: Vec<Distro> = vec![arch::scrape().await?, amazon_linux::scrape().await?];

    for distro in ubuntus {
        result.push(distro?);
    }

    Ok(result)
}

pub async fn cron() {
    let mut conn = crate::establish_connection().unwrap();
    'outer: loop {
        tokio::time::sleep_until(
            Instant::now()
                .checked_add(Duration::from_secs(24 * 60 * 60))
                .unwrap(),
        )
        .await;

        debug!("scraping distro images from upstream");

        let distros = get_all().await.unwrap();

        let tx = conn.transaction().unwrap();
        for d in distros {
            if let Err(why) = tx.execute(
                "UPDATE distros
                 SET download_url = ?1
                   , sha256sum    = ?2
                   , min_size     = ?3
                   , format       = ?4
                 WHERE name = ?5",
                params![d.download_url, d.sha256sum, d.min_size, d.format, d.name],
            ) {
                error!("can't update distros: {why}");
                continue 'outer;
            }

            if let Err(why) = tx.execute(
                "INSERT INTO audit_logs(kind, op, data) VALUES (?1, ?2, ?3)",
                params!["distro", "update", serde_json::to_string(&d).unwrap()],
            ) {
                error!("can't update audit logs: {why}");
                continue 'outer;
            }
        }

        if let Err(why) = tx.commit() {
            error!("can't commit transaction: {why}");
        }
    }
}
