use crate::{models::Instance, Error, State};
use axum::extract::{Extension, Path};
use rusqlite::params;
use serde::{Deserialize, Serialize};
use std::sync::Arc;
use uuid::Uuid;

#[instrument(err)]
pub async fn user_data(
    Path(id): Path<Uuid>,
    Extension(state): Extension<Arc<State>>,
) -> Result<String, Error> {
    let conn = state.pool.get().await?;

    conn.execute(
        "UPDATE instances SET status = ?1 WHERE uuid = ?2",
        params!["running", id],
    )?;
    let ins = Instance::from_uuid(&conn, id)?;
    conn.execute(
        "INSERT INTO audit_logs(kind, op, data) VALUES (?1, ?2, ?3)",
        params!["instance", "running", serde_json::to_string(&ins)?],
    )?;

    Ok(conn.query_row(
        "SELECT user_data FROM cloudconfig_seeds WHERE uuid = ?1",
        params![id],
        |row| Ok(row.get(0)?),
    )?)
}

#[instrument(err)]
pub async fn meta_data(
    Path(id): Path<Uuid>,
    Extension(state): Extension<Arc<State>>,
) -> Result<String, Error> {
    let conn = state.pool.get().await?;

    let hostname: String = conn.query_row(
        "SELECT name FROM instances WHERE uuid = ?1",
        params![id],
        |row| row.get(0),
    )?;

    Ok(format!(
        "instance-id: {}
local-hostname: {}",
        id, hostname,
    ))
}

#[instrument(err, skip(ts))]
pub async fn vendor_data(
    Path(id): Path<Uuid>,
    Extension(ts): Extension<Arc<tailscale_client::Client>>,
    Extension(state): Extension<Arc<State>>,
) -> Result<String, Error> {
    let conn = state.pool.get().await?;

    let i: Instance = Instance::from_uuid(&conn, id)?;

    if i.join_tailnet {
        let key_info = ts
            .create_key(tailscale_client::Capabilities {
                reusable: false,
                ephemeral: true,
            })
            .await?;

        conn.execute(
            "INSERT INTO audit_logs(kind, op, data) VALUES (?1, ?2, ?3)",
            params![
                "tailnet authkey",
                "create",
                serde_json::to_string(&key_info)?
            ],
        )?;

        if i.distro == "ubuntu-20.04".to_string() || i.distro == "ubuntu-22.04".to_string() {
            Ok(format!("#cloud-config\n{}", serde_yaml::to_string(&CloudConfig{
                write_files: vec![
                    File{
                        owner: "root:root".to_string(),
                        path: "/etc/update-motd.d/69-waifud".to_string(),
                        permissions: "0755".to_string(),
                        content: "#!/bin/sh\n#\n# This file is written by waifud.\necho \"\"\necho \"Welcome to waifud <3\"\n".to_string(),
                    },
                ],
                runcmd: vec![
                    vec!["sh".into(), "-c".into(), "curl -fsSL https://tailscale.com/install.sh | sh".into()],
                    vec!["systemctl".into(), "enable".into(), "--now".into(), "tailscaled.service".into()],
                    vec!["tailscale".into(), "up".into(), "--authkey".into(), key_info.key.unwrap(), "--ssh".into(), "--advertise-tags=tag:vm".into()],
                    vec!["apt".into(), "install".into(), "-y".into(), "systemd-container".into()]
                ],
            })?))
        } else {
            Ok(format!("#cloud-config\n{}", serde_yaml::to_string(&CloudConfig{
                write_files: vec![
                    File{
                        owner: "root:root".into(),
                        path: "/etc/update-motd.d/69-waifud".into(),
                        permissions: "0755".into(),
                        content: "#!/bin/sh\n#\n# This file is written by waifud.\necho \"\"\necho \"Welcome to waifud <3\"\n".into(),
                    },
                ],
                runcmd: vec![
                    vec!["sh".into(), "-c".into(), "curl -fsSL https://tailscale.com/install.sh | sh".into()],
                    vec!["systemctl".into(), "enable".into(), "--now".into(), "tailscaled.service".into()],
                    vec!["tailscale".into(), "up".into(), "--authkey".into(), key_info.key.unwrap(), "--ssh".into()],
                ],
            })?))
        }
    } else {
        Ok(include_str!("./vendor-data").to_string())
    }
}

#[derive(Default, Debug, Clone, PartialEq, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct CloudConfig {
    #[serde(rename = "write_files")]
    pub write_files: Vec<File>,
    pub runcmd: Vec<Vec<String>>,
}

#[derive(Default, Debug, Clone, PartialEq, Serialize, Deserialize)]
pub struct File {
    pub owner: String,
    pub path: String,
    pub permissions: String,
    pub content: String,
}
