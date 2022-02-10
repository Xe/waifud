use crate::{models::Instance, Error, State};
use axum::extract::{Extension, Path};
use rusqlite::params;
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

#[instrument]
pub async fn vendor_data(Path(_id): Path<Uuid>) -> &'static str {
    include_str!("./vendor-data")
}
