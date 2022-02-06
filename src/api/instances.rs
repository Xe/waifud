use crate::{
    establish_connection,
    libvirt::{random_mac, NewInstance},
    models::{Distro, Instance},
    namegen, Error,
};
use axum::Json;
use rusqlite::params;
use std::{
    io::{self, Write},
    net::SocketAddr,
    os::unix::prelude::ExitStatusExt,
    process::ExitStatus,
};
use tokio::{net::lookup_host, process::Command};
use uuid::Uuid;

#[instrument(err)]
#[axum_macros::debug_handler]
pub async fn make(Json(details): Json<NewInstance>) -> Result<Json<Instance>, Error> {
    let id = Uuid::new_v4();

    let addrs: Vec<SocketAddr> = lookup_host(details.host.clone() + ":22".into())
        .await?
        .collect();
    if addrs.len() == 0 {
        return Err(Error::HostDoesntExist(details.host));
    }

    let conn = establish_connection()?;

    let distro = conn.query_row(
        "SELECT name, download_url, sha256sum, min_size, format FROM distros WHERE name = ?1",
        params![details.distro.clone()],
        |row| {
            Ok(Distro {
                name: row.get(0)?,
                download_url: row.get(1)?,
                sha256sum: row.get(2)?,
                min_size: row.get(3)?,
                format: row.get(4)?,
            })
        },
    )?;

    let output = Command::new("ssh")
        .args([
            &details.host.clone(),
            "stat",
            &format!("$HOME/.cache/within/mkvm/qcow2/{}", distro.sha256sum),
        ])
        .output()
        .await?;
    if output.status != ExitStatus::from_raw(0) {
        let output = Command::new("ssh")
            .args([
                &details.host.clone(),
                "wget",
                "-O",
                &format!("$HOME/.cache/within/mkvm/qcow2/{}", distro.sha256sum),
                &distro.download_url,
            ])
            .output()
            .await?;
        if output.status != ExitStatus::from_raw(0) {
            let stderr = String::from_utf8(output.stderr).unwrap();
            return Err(Error::CantDownloadImage(
                distro.download_url.clone(),
                stderr,
            ));
        }
    }

    let mut buf: Vec<u8> = vec![];

    let details = NewInstance {
        name: details.name.or(Some(namegen::next())),
        memory_mb: details.memory_mb.or(Some(512)),
        host: details.host,
        disk_size_gb: details.disk_size_gb.or(Some(distro.min_size)),
        zvol_prefix: details.zvol_prefix.or(Some("rpool/safe/vms".into())),
        distro: distro.name,
        sata: details.sata.or(Some(false)),
        cpus: details.cpus.or(Some(2)),
        user_data: details
            .user_data
            .or(Some(include_str!("../../var/xe-base.yaml").into())),
    };

    crate::templates::base_xml(
        &mut buf,
        details.name.unwrap(),
        id.to_string(),
        random_mac(),
        details.zvol_prefix.unwrap(),
        details.sata.unwrap(),
        details.memory_mb.unwrap() * 1024,
        details.cpus.unwrap(),
        format!("http://100.87.242.16:23818/api/cloudinit/{}/", id),
    )?;

    Ok(Json(Instance::default()))
}
