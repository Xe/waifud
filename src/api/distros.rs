use crate::{establish_connection, models::Distro, Result};
use axum::Json;
use rusqlite::params;

#[instrument(err)]
fn list_distros() -> Result<Vec<Distro>> {
    let conn = establish_connection()?;

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

    Ok(result)
}

pub async fn get_distros() -> Result<Json<Vec<Distro>>> {
    let result = list_distros()?;
    Ok(Json(result))
}
