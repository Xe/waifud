#[macro_use]
extern crate diesel_migrations;
#[macro_use]
extern crate tracing;

use anyhow::{anyhow, Result};
use diesel::prelude::*;
use std::env;

diesel_migrations::embed_migrations!("./migrations");

pub fn establish_connection() -> Result<SqliteConnection> {
    let database_url = env::var("DATABASE_URL").unwrap_or("./var/waifud.db".to_string());
    SqliteConnection::establish(&database_url)
        .map_err(|why| anyhow!("can't connect to {}: {}", database_url, why))
}

fn main() -> Result<()> {
    tracing_subscriber::fmt::init();

    info!("{} migrator starting up", waifud::APPLICATION_NAME);

    info!("running migrations");
    let connection = establish_connection()?;
    embedded_migrations::run_with_output(&connection, &mut std::io::stdout()).map_err(|why| {
        error!("migration error: {}", why);
        why
    })?;
    info!("migrations succeeded");

    Ok(())
}
