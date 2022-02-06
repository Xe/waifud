#[macro_use]
extern crate tracing;

use axum::{
    http::StatusCode,
    response::{IntoResponse, Response},
};
use rusqlite::Connection;
use std::{env, fmt, net::AddrParseError};
use tokio::sync::Mutex;

pub const APPLICATION_NAME: &str = concat!(env!("CARGO_PKG_NAME"), "/", env!("CARGO_PKG_VERSION"),);

pub fn establish_connection() -> Result<Connection> {
    let database_url = env::var("DATABASE_URL").unwrap_or("./var/waifud.db".to_string());
    Ok(Connection::open(&database_url)?)
}

pub type Result<T = (), E = Error> = std::result::Result<T, E>;

pub mod api;
pub mod libvirt;
pub mod migrate;
pub mod models;
pub mod namegen;

pub struct State(Mutex<Connection>);

impl fmt::Debug for State {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        write!(f, "State()")
    }
}

impl State {
    pub fn new() -> Result<Self> {
        Ok(State(Mutex::new(establish_connection()?)))
    }
}

#[derive(thiserror::Error, Debug)]
pub enum Error {
    #[error("libvirt error: {0}")]
    Libvirt(#[from] virt::error::Error),

    #[error("dhall parsing error: {0}")]
    Dhall(#[from] serde_dhall::Error),

    #[error("database error: {0}")]
    SQLite(#[from] rusqlite::Error),

    #[error("internal tokio error: {0}")]
    TokioJoin(#[from] tokio::task::JoinError),

    #[error("hyper error: {0}")]
    Hyper(#[from] hyper::Error),

    #[error("address parse error: {0}")]
    AddrParse(#[from] AddrParseError),

    #[error("io error: {0}")]
    IO(#[from] std::io::Error),

    // Application errors
    #[error("host {0} doesn't exist")]
    HostDoesntExist(String),

    #[error("can't download {0}:\n\n{1}")]
    CantDownloadImage(String, String),

    #[error("can't create zfs zvol on {0}:\n\n{1}")]
    CantMakeZvol(String, String),

    #[error("can't hydrate zfs zvol on {0}:\n\n{1}")]
    CantHydrateZvol(String, String),

    #[error("can't create zfs init snapshot on {0}:\n\n{1}")]
    CantMakeInitSnapshot(String, String),
}

impl IntoResponse for Error {
    fn into_response(self) -> Response {
        let interm: (StatusCode, String) = self.into();
        interm.into_response()
    }
}

impl Into<(StatusCode, String)> for Error {
    fn into(self) -> (StatusCode, String) {
        match self {
            Error::Libvirt(why) => (StatusCode::INTERNAL_SERVER_ERROR, why.message),
            Error::Dhall(why) => (StatusCode::BAD_REQUEST, format!("{}", why)),
            Error::SQLite(err) => match err {
                rusqlite::Error::QueryReturnedNoRows => {
                    (StatusCode::NOT_FOUND, "404 not found".into())
                }
                _ => (StatusCode::INTERNAL_SERVER_ERROR, format!("{}", err)),
            },
            _ => (StatusCode::INTERNAL_SERVER_ERROR, format!("{}", self)),
        }
    }
}

include!(concat!(env!("OUT_DIR"), "/templates.rs"));
