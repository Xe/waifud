use crate::Error;
use axum::extract::Path;
use rusqlite::params;
use uuid::Uuid;

#[instrument(err)]
pub async fn user_data(Path(id): Path<Uuid>) -> Result<String, Error> {
    let conn = crate::establish_connection()?;

    Ok(conn.query_row(
        "SELECT user_data FROM cloudconfig_seeds WHERE uuid = ?1",
        params![id],
        |row| Ok(row.get(0)?),
    )?)
}

#[instrument(err)]
pub async fn meta_data(Path(id): Path<Uuid>) -> Result<String, Error> {
    let conn = crate::establish_connection()?;

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
