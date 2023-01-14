use crate::{models::Distro, tailauth::Tailauth, Result, State};
use axum::{
    extract::{Extension, Path},
    Json,
};
use rusqlite::params;
use std::sync::Arc;

#[instrument(err)]
pub async fn create(
    Extension(state): Extension<Arc<State>>,
    _: Tailauth,
    Json(distro): Json<Distro>,
) -> Result<Json<Distro>> {
    let conn = state.pool.get().await?;

    let mut distro = distro;
    if distro.format == "".to_string() {
        distro.format = "waifud://qcow2".into();
    }

    {
        let d = distro.clone();
        conn.execute(
            "INSERT INTO distros
                   ( name
                   , download_url
                   , sha256sum
                   , min_size
                   , format
                   )
             VALUES
                 ( ?1
                 , ?2
                 , ?3
                 , ?4
                 , ?5
                 )",
            params![d.name, d.download_url, d.sha256sum, d.min_size, d.format],
        )?;
    }

    conn.execute(
        "INSERT INTO audit_logs(kind, op, data) VALUES (?1, ?2, ?3)",
        params!["distro", "create", serde_json::to_string(&distro)?],
    )?;

    Ok(Json(distro))
}

#[instrument(err)]
pub async fn update(
    Path(name): Path<String>,
    Extension(state): Extension<Arc<State>>,
    _: Tailauth,
    Json(distro): Json<Distro>,
) -> Result<Json<Distro>> {
    let conn = state.pool.get().await?;

    let mut distro = distro;
    if distro.format == "".to_string() {
        distro.format = "waifud://qcow2".into();
    }

    let d = distro.clone();
    conn.execute(
        "
INSERT INTO
  distros( name
         , download_url
         , sha256sum
         , min_size
         , format
         )
VALUES ( ?5
       , ?1
       , ?2
       , ?3
       , ?4
       )
ON CONFLICT DO
  UPDATE SET download_url=?1
           , sha256sum=?2
           , min_size=?3
           , format=?4
",
        params![d.download_url, d.sha256sum, d.min_size, d.format, d.name],
    )?;

    conn.execute(
        "INSERT INTO audit_logs(kind, op, data) VALUES (?1, ?2, ?3)",
        params!["distro", "update", serde_json::to_string(&d)?],
    )?;

    Ok(Json(distro))
}

#[instrument(err)]
pub async fn delete(
    Extension(state): Extension<Arc<State>>,
    Path(name): Path<String>,
    _: Tailauth,
) -> Result<()> {
    let conn = state.pool.get().await?;

    let d = Distro::from_name(&conn, name)?;
    conn.execute("DELETE FROM distros WHERE name = ?1", params![d.name])?;

    conn.execute(
        "INSERT INTO audit_logs(kind, op, data) VALUES (?1, ?2, ?3)",
        params!["distro", "update", serde_json::to_string(&d)?],
    )?;

    Ok(())
}

#[instrument(err)]
pub async fn get(
    Extension(state): Extension<Arc<State>>,
    Path(name): Path<String>,
    _: Tailauth,
) -> Result<Json<Distro>> {
    let conn = state.pool.get().await?;

    Ok(Json(conn.query_row(
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
    )?))
}

#[instrument(err)]
pub async fn list(
    Extension(state): Extension<Arc<State>>,
    _: Tailauth,
) -> Result<Json<Vec<Distro>>> {
    let conn = state.pool.get().await?;

    let mut stmt = conn.prepare(
        "SELECT name, download_url, sha256sum, min_size, format FROM distros ORDER BY name ASC",
    )?;
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
