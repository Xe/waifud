use crate::{models::Distro, Result, State};
use axum::{extract::Extension, Json};
use rusqlite::params;
use std::sync::Arc;

#[instrument(err)]
pub async fn get_distros(Extension(state): Extension<Arc<State>>) -> Result<Json<Vec<Distro>>> {
    let conn = state.0.lock().await;

    let mut stmt =
        conn.prepare("SELECT name, download_url, sha256sum, min_size, format FROM distros")?;
    let iter = stmt.query_map(params![], |row| {
        Ok(Distro {
            name: row.get(0)?,
            download_url: row.get(1)?,
            sha256sum: row.get(2)?,
            min_size: row.get(3)?,
            format: row.get(4)?,
        })
    })?;
    let mut result: Vec<Distro> = vec![];

    for distro in iter {
        result.push(distro.unwrap());
    }

    Ok(Json(result))
}
