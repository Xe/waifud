use crate::{models::Distro, Error};
use rusqlite::{params, Connection};

#[instrument(err)]
pub fn run() -> Result<(), crate::Error> {
    info!("running");
    let connection = crate::establish_connection()?;

    connection.execute_batch(include_str!("./schema.sql"))?;

    load_distros(connection)?;

    Ok(())
}

fn load_distros(conn: Connection) -> Result<(), Error> {
    let count = conn.query_row("SELECT COUNT(*) FROM distros", params![], |row| {
        row.get::<_, i64>(0)
    })?;

    if count != 0 {
        return Ok(());
    }

    let distros: Vec<Distro> =
        serde_dhall::from_str(include_str!("../../data/distros.dhall")).parse()?;

    let mut stmt = conn.prepare("INSERT INTO distros(name, download_url, sha256sum, min_size, format) VALUES (?1, ?2, ?3, ?4, ?5)")?;

    for distro in distros {
        stmt.execute(params![
            distro.name,
            distro.download_url,
            distro.sha256sum,
            distro.min_size,
            distro.format,
        ])?;
    }

    Ok(())
}
