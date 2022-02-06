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
    time::Duration,
};
use tokio::{net::lookup_host, process::Command, task::spawn_blocking};
use uuid::Uuid;
use virt::{
    connect::Connect,
    domain::{Domain, VIR_DOMAIN_INTERFACE_ADDRESSES_SRC_LEASE},
};

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

    let details = NewInstance {
        name: details.name.or(Some(namegen::next())),
        memory_mb: details.memory_mb.or(Some(512)),
        host: details.host.clone(),
        disk_size_gb: details.disk_size_gb.or(Some(distro.min_size)),
        zvol_prefix: details.zvol_prefix.or(Some("rpool/safe/vms".into())),
        distro: distro.name,
        sata: details.sata.or(Some(false)),
        cpus: details.cpus.or(Some(2)),
        user_data: details
            .user_data
            .or(Some(include_str!("../../var/xe-base.yaml").into())),
    };

    debug!("name: {}", details.name.as_ref().unwrap());

    debug!("checking if image exists");
    let output = Command::new("ssh")
        .args([
            &details.host.clone(),
            "stat",
            &format!("$HOME/.cache/within/mkvm/qcow2/{}", distro.sha256sum),
        ])
        .output()
        .await?;
    if output.status != ExitStatus::from_raw(0) {
        debug!("downloading image");
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

            Command::new("ssh")
                .args([
                    &details.host.clone(),
                    "rm",
                    &format!("$HOME/.cache/within/mkvm/qcow2/{}", distro.sha256sum),
                ])
                .status()
                .await?;

            return Err(Error::CantDownloadImage(
                distro.download_url.clone(),
                stderr,
            ));
        }
    }
    debug!("making zvol");
    let output = Command::new("ssh")
        .args([
            &details.host.clone(),
            "sudo",
            "zfs",
            "create",
            "-V",
            &format!("{}G", details.disk_size_gb.unwrap()),
            &format!(
                "{}/{}",
                details.zvol_prefix.as_ref().unwrap(),
                &details.name.as_ref().unwrap()
            ),
        ])
        .output()
        .await?;
    if output.status != ExitStatus::from_raw(0) {
        let stderr = String::from_utf8(output.stderr).unwrap();
        return Err(Error::CantMakeZvol(details.host.clone(), stderr));
    }
    debug!("hydrating zvol");
    let output = Command::new("ssh")
        .args([
            &details.host.clone(),
            "sudo",
            "qemu-img",
            "convert",
            "-O",
            "raw",
            &format!("$HOME/.cache/within/mkvm/qcow2/{}", distro.sha256sum),
            &format!(
                "/dev/zvol/{}/{}",
                details.zvol_prefix.as_ref().unwrap(),
                &details.name.as_ref().unwrap()
            ),
        ])
        .output()
        .await?;
    if output.status != ExitStatus::from_raw(0) {
        let stderr = String::from_utf8(output.stderr).unwrap();
        return Err(Error::CantHydrateZvol(details.host.clone(), stderr));
    }
    debug!("making init snapshot");
    let output = Command::new("ssh")
        .args([
            &details.host.clone(),
            "sudo",
            "zfs",
            "snapshot",
            &format!(
                "{}/{}@init",
                details.zvol_prefix.as_ref().unwrap(),
                &details.name.as_ref().unwrap()
            ),
        ])
        .output()
        .await?;
    if output.status != ExitStatus::from_raw(0) {
        let stderr = String::from_utf8(output.stderr).unwrap();
        return Err(Error::CantMakeInitSnapshot(details.host.clone(), stderr));
    }

    conn.execute(
        "INSERT INTO cloudconfig_seeds(uuid, user_data) VALUES (?1, ?2)",
        params![id.clone(), details.user_data.clone().unwrap()],
    )?;

    let mut buf: Vec<u8> = vec![];

    let mac_addr = random_mac();

    debug!("rendering xml");
    crate::templates::base_xml(
        &mut buf,
        details.name.clone().unwrap(),
        id.to_string(),
        mac_addr.clone(),
        details.zvol_prefix.clone().unwrap(),
        details.sata.unwrap(),
        details.memory_mb.unwrap() * 1024,
        details.cpus.unwrap(),
        format!("http://100.87.242.16:23818/api/cloudinit/{}/", id),
    )?;

    let buf = String::from_utf8(buf).unwrap();
    trace!("libvirt xml:\n{}", buf);

    let addr: Result<Option<String>, Error> = {
        let details = details.clone();
        spawn_blocking(move || {
            debug!("connecting to host");
            let lc = Connect::open(&format!("qemu+ssh://root@{}/system", details.host.clone()))?;

            debug!("defining domain");
            let dom = Domain::define_xml(&lc, &buf)?;

            debug!("starting domain");
            dom.create()?;

            debug!("waiting for DHCP");
            std::thread::sleep(Duration::from_millis(2000));

            let mut addr: Vec<String> = dom
                .interface_addresses(VIR_DOMAIN_INTERFACE_ADDRESSES_SRC_LEASE, 0)?
                .into_iter()
                .filter(|iface| iface.addrs.clone().get(0).is_some())
                .map(|iface| iface.addrs.clone().get(0).unwrap().clone().addr)
                .collect();

            if addr.get(0).is_none() {
                Ok(None)
            } else {
                Ok(Some(addr.swap_remove(0)))
            }
        })
        .await?
    };
    let addr = addr?;

    debug!("ip: {:?}", addr);
    let zvol_name = format!(
        "{}/{}",
        details.zvol_prefix.clone().unwrap(),
        details.name.clone().unwrap()
    );

    let ins = Instance {
        uuid: id,
        name: details.name.unwrap(),
        host: details.host,
        memory: details.memory_mb.unwrap(),
        disk_size: details.disk_size_gb.unwrap(),
        mac_address: mac_addr.clone(),
        zvol_name: zvol_name.clone(),
    };

    {
        let ins = ins.clone();
        conn.execute(
            "INSERT INTO instances(uuid, name, host, mac_address, memory, disk_size, zvol_name) VALUES (?1, ?2, ?3, ?4, ?5, ?6, ?7)",
            params![
                ins.uuid,
                ins.name,
                ins.host,
                mac_addr,
                ins.memory,
                ins.disk_size,
                zvol_name,
            ],
        )?;
    }

    Ok(Json(ins))
}
