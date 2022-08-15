use crate::establish_connection;
use anyhow::Result;
use rusqlite_migration::{Migrations, M};

#[instrument(err)]
pub fn run() -> Result<()> {
    info!("running");
    let mut conn = establish_connection()?;

    let migrations = Migrations::new(vec![
        M::up(include_str!("./base_schema.sql")),
        M::up(include_str!("./20220225-session.sql")),
        M::up(include_str!("./20220814-no-session.sql")),
    ]);
    conn.pragma_update(None, "journal_mode", &"WAL").unwrap();

    migrations.to_latest(&mut conn)?;

    Ok(())
}
