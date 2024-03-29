use bb8::PooledConnection;
use bb8_rusqlite::RusqliteConnectionManager;
use rusqlite::params;
use serde::{Deserialize, Serialize};
use uuid::Uuid;

use crate::Result;

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
    pub join_tailnet: bool,
}

impl Instance {
    pub fn from_name(
        conn: &PooledConnection<'_, RusqliteConnectionManager>,
        name: String,
    ) -> Result<Self> {
        let mut stmt = conn.prepare(
            "SELECT uuid, name, host, mac_address, memory, disk_size, zvol_name, status, distro, join_tailnet FROM instances WHERE name = ?1",
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
                join_tailnet: row.get(9)?,
            })
        })?;

        Ok(instance)
    }

    pub fn from_uuid(
        conn: &PooledConnection<'_, RusqliteConnectionManager>,
        id: Uuid,
    ) -> Result<Self> {
        let mut stmt = conn.prepare(
            "SELECT uuid, name, host, mac_address, memory, disk_size, zvol_name, status, distro, join_tailnet FROM instances WHERE uuid = ?1",
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
                join_tailnet: row.get(9)?,
            })
        })?;

        Ok(instance)
    }
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

impl Distro {
    pub fn from_name(
        conn: &PooledConnection<'_, RusqliteConnectionManager>,
        name: String,
    ) -> Result<Self> {
        Ok(conn.query_row(
            "SELECT
           name
         , download_url
         , sha256sum
         , min_size
         , format
         FROM distros
         WHERE name = ?1",
            params![name],
            |row| {
                Ok(Distro {
                    name: row.get(0)?,
                    download_url: row.get(1)?,
                    sha256sum: row.get(2)?,
                    min_size: row.get(3)?,
                    format: row.get(4)?,
                })
            },
        )?)
    }
}

#[derive(Deserialize, Serialize, Debug, Clone)]
pub struct AuditEvent {
    pub id: i32,
    pub ts: i64,
    pub kind: String,
    pub op: String,
    pub data: Option<serde_json::Value>,
    pub uuid: Option<String>,
    pub name: Option<String>,
}

impl AuditEvent {
    pub fn get_for_instance(
        uuid: Uuid,
        conn: &PooledConnection<'_, RusqliteConnectionManager>,
    ) -> Result<Vec<AuditEvent>> {
        let mut stmt =
            conn.prepare("SELECT id, ts, kind, op, data, uuid, name from audit_logs where uuid=? and kind='instance'")?;

        let vals: Vec<AuditEvent> = stmt
            .query_map(params![uuid.to_string()], |row| {
                Ok(AuditEvent {
                    id: row.get(0)?,
                    ts: row.get(1)?,
                    kind: row.get(2)?,
                    op: row.get(3)?,
                    data: row.get(4)?,
                    uuid: row.get(5)?,
                    name: row.get(6)?,
                })
            })?
            .into_iter()
            .map(Result::unwrap)
            .collect();
        Ok(vals)
    }

    pub fn get_all(
        conn: &PooledConnection<'_, RusqliteConnectionManager>,
    ) -> Result<Vec<AuditEvent>> {
        let mut stmt = conn.prepare("SELECT id, ts, kind, op, data, uuid, name from audit_logs")?;
        let vals: Vec<AuditEvent> = stmt
            .query_map(params![], |row| {
                Ok(AuditEvent {
                    id: row.get(0)?,
                    ts: row.get(1)?,
                    kind: row.get(2)?,
                    op: row.get(3)?,
                    data: row.get(4)?,
                    uuid: row.get(5)?,
                    name: row.get(6)?,
                })
            })?
            .into_iter()
            .map(Result::unwrap)
            .collect();
        Ok(vals)
    }
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Session {
    pub uuid: Uuid,
    pub user: String,
    pub expired: bool,
}

impl Session {
    #[instrument(skip(conn), err)]
    pub fn get(
        conn: &PooledConnection<'_, RusqliteConnectionManager>,
        id: Uuid,
    ) -> Result<Session> {
        Ok(conn.query_row(
            "SELECT (uuid, user, expired) FROM sessions WHERE uuid = ?1",
            params![id],
            |row| {
                Ok(Session {
                    uuid: row.get(0)?,
                    user: row.get(1)?,
                    expired: row.get(2)?,
                })
            },
        )?)
    }
}
