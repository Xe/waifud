use crate::{models::Distro, Result, State};
use axum::{
    extract::{Extension, Path},
    Json,
};
use rusqlite::params;
use std::sync::Arc;

#[instrument(err)]
pub async fn create(
    Extension(state): Extension<Arc<State>>,
    Json(distro): Json<Distro>,
) -> Result<Json<Distro>> {
    let conn = state.0.lock().await;

    let mut distro = distro;
    if distro.format == "".to_string() {
        distro.format = "waifud://qcow2".into();
    }

    let d = distro.clone();
    conn.execute(
        "INSERT INTO distros (name, download_url, sha256sum, min_size, format) VALUES (?1, ?2, ?3, ?4, ?5)",
        params![d.name, d.download_url, d.sha256sum, d.min_size, d.format],
    )?;

    Ok(Json(distro))
}

#[instrument(err)]
pub async fn delete(
    Extension(state): Extension<Arc<State>>,
    Path(name): Path<String>,
) -> Result<()> {
    let conn = state.0.lock().await;

    conn.execute("DELETE FROM distros WHERE name = ?1", params![name])?;

    Ok(())
}

#[instrument(err)]
pub async fn list(Extension(state): Extension<Arc<State>>) -> Result<Json<Vec<Distro>>> {
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
