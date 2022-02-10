use crate::{
    api::libvirt::Machine,
    libvirt::{random_mac, NewInstance},
    models::{Distro, Instance},
    namegen, Config, Error, State,
};
use axum::{
    extract::{Extension, Path},
    Json,
};
use rusqlite::params;
use std::{
    convert::TryFrom, net::SocketAddr, os::unix::prelude::ExitStatusExt, process::ExitStatus,
    sync::Arc, time::Duration,
};
use tokio::{net::lookup_host, process::Command, task::spawn_blocking, time::sleep};
use uuid::Uuid;
use virt::{connect::Connect, domain::Domain};

#[instrument(err)]
#[axum_macros::debug_handler]
pub async fn delete(
    Path(id): Path<Uuid>,
    Extension(state): Extension<Arc<State>>,
) -> Result<(), Error> {
    let conn = state.pool.get().await?;

    let info: Result<(String, String), Error> = {
        let id = id.clone();
        let mut stmt = conn.prepare("SELECT host, zvol_name FROM instances WHERE uuid = ?1")?;
        let info: (String, String) =
            stmt.query_row(params![id], |row| Ok((row.get(0)?, row.get(1)?)))?;
        Ok(info)
    };
    let (host, zvol_name): (String, String) = info?;

    let nuke: Result<(), Error> = {
        let host = host.clone();
        let id = id.clone();

        spawn_blocking(move || {
            let conn = Connect::open(&format!("qemu+ssh://root@{}/system", host))?;

            let dom = Domain::lookup_by_uuid_string(&conn, &id.to_string())?;

            dom.destroy()?;
            dom.undefine()?;
            Ok(())
        })
        .await?
    };
    nuke?;

    sleep(Duration::from_millis(500)).await;

    debug!("destroying zvol");
    let output = Command::new("ssh")
        .args([&host, "sudo", "zfs", "destroy", "-rf", &zvol_name])
        .output()
        .await?;
    if output.status != ExitStatus::from_raw(0) {
        let stderr = String::from_utf8(output.stderr).unwrap();
        return Err(Error::CantDeleteZvol(host.clone(), stderr));
    }

    conn.execute("DELETE FROM instances WHERE uuid = ?1", params![id])?;
    conn.execute("DELETE FROM cloudconfig_seeds WHERE uuid = ?1", params![id])?;

    Ok(())
}

#[instrument(err)]
#[axum_macros::debug_handler]
pub async fn get_machine(
    Path(id): Path<Uuid>,
    Extension(state): Extension<Arc<State>>,
) -> Result<Json<Machine>, Error> {
    let conn = state.pool.get().await?;

    let mut stmt = conn.prepare("SELECT host FROM instances WHERE uuid = ?1")?;
    let host: String = stmt.query_row(params![id], |row| row.get(0))?;

    let conn = Connect::open(&format!("qemu+ssh://root@{}/system", host))?;

    let dom = Domain::lookup_by_uuid_string(&conn, &id.to_string())?;
    Ok(Json(Machine::try_from(dom)?))
}

#[instrument(err)]
#[axum_macros::debug_handler]
pub async fn hard_reboot(
    Path(id): Path<Uuid>,
    Extension(state): Extension<Arc<State>>,
) -> Result<(), Error> {
    let conn = state.pool.get().await?;

    let mut stmt = conn.prepare("SELECT host FROM instances WHERE uuid = ?1")?;
    let host: String = stmt.query_row(params![id], |row| row.get(0))?;

    let vc = Connect::open(&format!("qemu+ssh://root@{}/system", host))?;

    let dom = Domain::lookup_by_uuid_string(&vc, &id.to_string())?;

    dom.destroy()?;
    dom.create()?;

    conn.execute(
        "UPDATE instances SET status = ?1 WHERE uuid = ?2",
        params!["rebooting", id],
    )?;

    Ok(())
}

#[instrument(err)]
#[axum_macros::debug_handler]
pub async fn shutdown(
    Path(id): Path<Uuid>,
    Extension(state): Extension<Arc<State>>,
) -> Result<(), Error> {
    let conn = state.pool.get().await?;

    let mut stmt = conn.prepare("SELECT host FROM instances WHERE uuid = ?1")?;
    let host: String = stmt.query_row(params![id], |row| row.get(0))?;

    let vc = Connect::open(&format!("qemu+ssh://root@{}/system", host))?;

    let dom = Domain::lookup_by_uuid_string(&vc, &id.to_string())?;
    dom.shutdown()?;

    conn.execute(
        "UPDATE instances SET status = ?1 WHERE uuid = ?2",
        params!["off", id],
    )?;

    Ok(())
}

#[instrument(err)]
#[axum_macros::debug_handler]
pub async fn start(
    Path(id): Path<Uuid>,
    Extension(state): Extension<Arc<State>>,
) -> Result<(), Error> {
    let conn = state.pool.get().await?;

    let mut stmt = conn.prepare("SELECT host FROM instances WHERE uuid = ?1")?;
    let host: String = stmt.query_row(params![id], |row| row.get(0))?;

    let vc = Connect::open(&format!("qemu+ssh://root@{}/system", host))?;

    let dom = Domain::lookup_by_uuid_string(&vc, &id.to_string())?;
    dom.create()?;

    conn.execute(
        "UPDATE instances SET status = ?1 WHERE uuid = ?2",
        params!["starting", id],
    )?;

    Ok(())
}

#[instrument(err)]
#[axum_macros::debug_handler]
pub async fn reboot(
    Path(id): Path<Uuid>,
    Extension(state): Extension<Arc<State>>,
) -> Result<(), Error> {
    let conn = state.pool.get().await?;

    let mut stmt = conn.prepare("SELECT host FROM instances WHERE uuid = ?1")?;
    let host: String = stmt.query_row(params![id], |row| row.get(0))?;

    let vc = Connect::open(&format!("qemu+ssh://root@{}/system", host))?;

    let dom = Domain::lookup_by_uuid_string(&vc, &id.to_string())?;
    dom.reboot(0)?;

    conn.execute(
        "UPDATE instances SET status = ?1 WHERE uuid = ?2",
        params!["rebooting", id],
    )?;

    Ok(())
}

#[instrument(err)]
#[axum_macros::debug_handler]
pub async fn get_by_name(
    Path(name): Path<String>,
    Extension(state): Extension<Arc<State>>,
) -> Result<Json<Instance>, Error> {
    let conn = state.pool.get().await?;

    let mut stmt = conn.prepare(
        "SELECT uuid, name, host, mac_address, memory, disk_size, zvol_name, status, distro FROM instances WHERE name = ?1",
    )?;
    let instance = stmt.query_row(params![name], |row| {
        Ok(Instance {
            uuid: row.get(0)?,
            name: row.get(1)?,
            host: row.get(2)?,
            mac_address: row.get(3)?,
            memory: row.get(4)?,
            disk_size: row.get(5)?,
            zvol_name: row.get(6)?,
            status: row.get(7)?,
            distro: row.get(8)?,
        })
    })?;

    Ok(Json(instance))
}

#[instrument(err)]
#[axum_macros::debug_handler]
pub async fn get(
    Path(id): Path<Uuid>,
    Extension(state): Extension<Arc<State>>,
) -> Result<Json<Instance>, Error> {
    let conn = state.pool.get().await?;

    let mut stmt = conn.prepare(
        "SELECT uuid, name, host, mac_address, memory, disk_size, zvol_name, status, distro FROM instances WHERE uuid = ?1",
    )?;
    let instance = stmt.query_row(params![id], |row| {
        Ok(Instance {
            uuid: row.get(0)?,
            name: row.get(1)?,
            host: row.get(2)?,
            mac_address: row.get(3)?,
            memory: row.get(4)?,
            disk_size: row.get(5)?,
            zvol_name: row.get(6)?,
            status: row.get(7)?,
            distro: row.get(8)?,
        })
    })?;

    Ok(Json(instance))
}

#[instrument(err)]
#[axum_macros::debug_handler]
pub async fn list(Extension(state): Extension<Arc<State>>) -> Result<Json<Vec<Instance>>, Error> {
    let conn = state.pool.get().await?;

    let mut result: Vec<Instance> = Vec::new();

    let mut stmt = conn.prepare(
        "SELECT uuid, name, host, mac_address, memory, disk_size, zvol_name, status, distro FROM instances",
    )?;
    let instances = stmt.query_map(params![], |row| {
        Ok(Instance {
            uuid: row.get(0)?,
            name: row.get(1)?,
            host: row.get(2)?,
            mac_address: row.get(3)?,
            memory: row.get(4)?,
            disk_size: row.get(5)?,
            zvol_name: row.get(6)?,
            status: row.get(7)?,
            distro: row.get(8)?,
        })
    })?;

    for instance in instances {
        result.push(instance?);
    }

    Ok(Json(result))
}

#[instrument(err)]
#[axum_macros::debug_handler]
pub async fn create(
    Json(details): Json<NewInstance>,
    Extension(state): Extension<Arc<State>>,
    Extension(config): Extension<Arc<Config>>,
) -> Result<Json<Instance>, Error> {
    let id = Uuid::new_v4();

    let addrs: Vec<SocketAddr> = lookup_host(details.host.clone() + ":22".into())
        .await?
        .collect();
    if addrs.len() == 0 {
        return Err(Error::HostDoesntExist(details.host));
    }

    let conn = state.pool.get().await?;

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
        distro: distro.name.clone(),
        sata: details.sata.or(Some(false)),
        cpus: details.cpus.or(Some(2)),
        user_data: details
            .user_data
            .or(Some(include_str!("../../var/xe-base.yaml").into())),
    };

    let mac_addr = random_mac();
    let zvol_name = format!(
        "{}/{}",
        details.zvol_prefix.clone().unwrap(),
        details.name.clone().unwrap()
    );

    let ins = Instance {
        uuid: id,
        name: details.name.clone().unwrap(),
        host: details.host.clone(),
        memory: details.memory_mb.unwrap(),
        disk_size: details.disk_size_gb.unwrap(),
        mac_address: mac_addr.clone(),
        zvol_name: zvol_name.clone(),
        status: "init".into(),
        distro: details.distro.clone(),
    };

    {
        let ins = ins.clone();
        conn.execute(
            "INSERT INTO instances(uuid, name, host, mac_address, memory, disk_size, zvol_name, status, distro) VALUES (?1, ?2, ?3, ?4, ?5, ?6, ?7, ?8, ?9)",
            params![
                ins.uuid,
                ins.name,
                ins.host,
                mac_addr,
                ins.memory,
                ins.disk_size,
                zvol_name,
                ins.status,
                ins.distro,
            ],
        )?;

        conn.execute(
            "INSERT INTO cloudconfig_seeds(uuid, user_data) VALUES (?1, ?2)",
            params![id.clone(), details.user_data.clone().unwrap()],
        )?;
    }

    drop(conn);

    tokio::spawn(async move {
        if let Err(why) = make_instance(config, state, details, distro, mac_addr, id).await {
            error!("can't make instance: {}", why);
        }
    });

    Ok(Json(ins))
}

async fn make_instance(
    config: Arc<Config>,
    state: Arc<State>,
    details: NewInstance,
    distro: Distro,
    mac_addr: String,
    id: Uuid,
) -> Result<(), Error> {
    let conn = state.pool.get().await?;

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
        conn.execute(
            "UPDATE instances SET status = ?1 WHERE uuid = ?2",
            params!["downloading image", id],
        )?;
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
    conn.execute(
        "UPDATE instances SET status = ?1 WHERE uuid = ?2",
        params!["hydrating zvol", id],
    )?;
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

    let mut buf: Vec<u8> = vec![];

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
        format!("{}/api/cloudinit/{}/", config.clone().base_url, id),
    )?;

    let buf = String::from_utf8(buf).unwrap();
    trace!("libvirt xml:\n{}", buf);

    let addr: Result<(), Error> = {
        let details = details.clone();
        spawn_blocking(move || {
            debug!("connecting to host");
            let lc = Connect::open(&format!("qemu+ssh://root@{}/system", details.host.clone()))?;

            debug!("defining domain");
            let dom = Domain::define_xml(&lc, &buf)?;

            debug!("starting domain");
            dom.create()?;
            Ok(())
        })
        .await?
    };
    let addr = addr?;

    debug!("ip: {:?}", addr);

    conn.execute(
        "UPDATE instances SET status = ?1 WHERE uuid = ?2",
        params!["waiting for cloud-init", id],
    )?;

    Ok(())
}
