use crate::{Error, State};
use axum::extract::{Extension, Path};
use rusqlite::params;
use std::sync::Arc;
use uuid::Uuid;

#[instrument(err)]
pub async fn user_data(
    Path(id): Path<Uuid>,
    Extension(state): Extension<Arc<State>>,
) -> Result<String, Error> {
    let conn = state.0.lock().await;

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
    let conn = state.0.lock().await;

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
