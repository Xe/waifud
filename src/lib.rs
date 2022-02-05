use rusqlite::{params, Connection};
use std::env;

pub const APPLICATION_NAME: &str = concat!(env!("CARGO_PKG_NAME"), "/", env!("CARGO_PKG_VERSION"),);

pub fn establish_connection() -> Result<Connection, rusqlite::Error> {
    let database_url = env::var("DATABASE_URL").unwrap_or("./var/waifud.db".to_string());
    Connection::open(&database_url)
}

pub mod models;

include!(concat!(env!("OUT_DIR"), "/templates.rs"));
