use crate::{establish_connection, models::Distro};
use anyhow::Result;
use rusqlite::{params, Connection};
use rusqlite_migration::{Migrations, M};

#[instrument(err)]
pub fn run() -> Result<()> {
    info!("running");
    let mut conn = establish_connection()?;

    let migrations = Migrations::new(vec![M::up(include_str!("./base_schema.sql"))]);
    conn.pragma_update(None, "journal_mode", &"WAL").unwrap();

    migrations.to_latest(&mut conn)?;

    load_distros(conn)?;

    Ok(())
}

fn load_distros(conn: Connection) -> Result<()> {
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
