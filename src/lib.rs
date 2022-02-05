#[macro_use]
extern crate tracing;

use axum::{
    http::StatusCode,
    response::{IntoResponse, Response},
};
use rusqlite::Connection;
use std::{env, net::AddrParseError};

pub const APPLICATION_NAME: &str = concat!(env!("CARGO_PKG_NAME"), "/", env!("CARGO_PKG_VERSION"),);

pub fn establish_connection() -> Result<Connection, rusqlite::Error> {
    let database_url = env::var("DATABASE_URL").unwrap_or("./var/waifud.db".to_string());
    Connection::open(&database_url)
}

pub mod api;
pub mod migrate;
pub mod models;
pub mod namegen;

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
            Error::TokioJoin(why) => (StatusCode::INTERNAL_SERVER_ERROR, format!("{}", why)),
            Error::Hyper(why) => (StatusCode::INTERNAL_SERVER_ERROR, format!("{}", why)),
            Error::AddrParse(why) => (StatusCode::INTERNAL_SERVER_ERROR, format!("{}", why)),
        }
    }
}

include!(concat!(env!("OUT_DIR"), "/templates.rs"));
