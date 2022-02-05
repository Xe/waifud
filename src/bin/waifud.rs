#[macro_use]
extern crate tracing;

use axum::{
    extract::Path,
    http::StatusCode,
    response::{IntoResponse, Response},
    routing::get,
    AddExtensionLayer, Json, Router,
};
use rusqlite::{params, Connection};
use serde::{Deserialize, Serialize};
use std::{net::AddrParseError, sync::Arc};
use tokio::{sync::Mutex, task::JoinError};
use tower::limit::ConcurrencyLimitLayer;
use tower_http::trace::TraceLayer;
use uuid::Uuid;
use virt::connect::Connect;
use waifud::models::Distro;

#[derive(Clone)]
pub struct State(Arc<Mutex<Connection>>);

#[tokio::main]
async fn main() -> Result<(), Error> {
    tracing_subscriber::fmt::init();
    let middleware = tower::ServiceBuilder::new()
        .layer(TraceLayer::new_for_http())
        .layer(ConcurrencyLimitLayer::new(64))
        .layer(AddExtensionLayer::new(State(Arc::new(Mutex::new(
            waifud::establish_connection()?,
        )))));

    // build our application with a route
    let app = Router::new()
        .route("/api/v1/distros", get(get_distros))
        .route("/api/v1/libvirt/machines", get(get_machines))
        .route("/api/cloudinit/:id/user-data", get(user_data))
        .route("/api/cloudinit/:id/vendor-data", get(vendor_data))
        .route("/api/cloudinit/:id/meta-data", get(meta_data))
        .layer(middleware);

    // run it
    let addr = &"[::]:23818".parse()?;
    info!("listening on {}", addr);
    axum::Server::bind(addr)
        .serve(app.into_make_service())
        .await?;

    Ok(())
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
    TokioJoin(#[from] JoinError),

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
            Error::SQLite(why) => (StatusCode::INTERNAL_SERVER_ERROR, format!("{}", why)),
            Error::TokioJoin(why) => (StatusCode::INTERNAL_SERVER_ERROR, format!("{}", why)),
            Error::Hyper(why) => (StatusCode::INTERNAL_SERVER_ERROR, format!("{}", why)),
            Error::AddrParse(why) => (StatusCode::INTERNAL_SERVER_ERROR, format!("{}", why)),
        }
    }
}

#[instrument(skip(uri))]
fn list_all_vms(uri: &str, host: String) -> Result<Vec<Machine>, Error> {
    debug!("connecting to {}: {}", host, uri);
    let mut conn = Connect::open(uri)?;
    let mut result = vec![];

    for dom in conn.list_all_domains(0)? {
        let addr: Option<String> = if dom.is_active()? {
            let mut addr: Vec<String> = dom
                .interface_addresses(virt::domain::VIR_DOMAIN_INTERFACE_ADDRESSES_SRC_LEASE, 0)?
                .into_iter()
                .filter(|iface| iface.addrs.clone().get(0).is_some())
                .map(|iface| iface.addrs.clone().get(0).unwrap().clone().addr)
                .collect();

            if addr.get(0).is_none() {
                None
            } else {
                Some(addr.swap_remove(0))
            }
        } else {
            None
        };

        let max_memory = dom.get_max_memory()? / 1024;

        result.push(Machine {
            name: dom.get_name()?,
            host: host.clone(),
            active: dom.is_active()?,
            uuid: dom.get_uuid_string()?,
            addr,
            memory_megs: max_memory,
        });
    }
    conn.close()?;

    Ok(result)
}

#[instrument]
fn list_distros() -> Result<Vec<Distro>, Error> {
    let conn = waifud::establish_connection()?;

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

async fn get_distros() -> Result<Json<Vec<Distro>>, Error> {
    let result = list_distros()?;
    Ok(Json(result))
}

#[derive(Default, Debug, Clone, PartialEq, Serialize, Deserialize)]
pub struct Machine {
    pub name: String,
    pub host: String,
    pub active: bool,
    pub uuid: String,
    pub addr: Option<String>,
    pub memory_megs: u64,
}

async fn get_machines() -> Result<Json<Vec<Machine>>, Error> {
    let mut result = Vec::new();
    for host in &["kos-mos", "logos", "ontos", "pneuma"] {
        result.extend_from_slice(&list_all_vms(
            &format!("qemu+ssh://root@{}/system", host),
            host.to_string(),
        )?);
    }
    Ok(Json(result))
}

async fn user_data(Path(id): Path<Uuid>) -> Result<String, Error> {
    let conn = waifud::establish_connection()?;

    Ok(conn.query_row(
        "SELECT user_data FROM cloudconfig_seeds WHERE uuid = ?1",
        params![id],
        |row| Ok(row.get(0)?),
    )?)
}

async fn meta_data(Path(id): Path<Uuid>) -> String {
    format!(
        "instance-id: hunter2
local-hostname: {}",
        id,
    )
}

async fn vendor_data(Path(_id): Path<Uuid>) -> &'static str {
    include_str!("./vendor-data")
}
