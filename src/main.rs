#[macro_use]
extern crate tracing;

use axum::{routing::get, AddExtensionLayer, Json, Router};
use rusqlite::{params, Connection};
use std::sync::Arc;
use tokio::sync::Mutex;
use tower::limit::ConcurrencyLimitLayer;
use tower_http::trace::TraceLayer;
use waifud::{models::Distro, Error};

#[derive(Clone)]
pub struct State(Arc<Mutex<Connection>>);

#[tokio::main]
async fn main() -> Result<(), Error> {
    tracing_subscriber::fmt::init();

    waifud::migrate::run()?;

    let middleware = tower::ServiceBuilder::new()
        .layer(TraceLayer::new_for_http())
        .layer(ConcurrencyLimitLayer::new(64))
        .layer(AddExtensionLayer::new(State(Arc::new(Mutex::new(
            waifud::establish_connection()?,
        )))));

    // build our application with a route
    let app = Router::new()
        .route("/api/v1/distros", get(get_distros))
        .route(
            "/api/v1/libvirt/machines",
            get(waifud::api::libvirt::get_machines),
        )
        .route(
            "/api/cloudinit/:id/meta-data",
            get(waifud::api::cloudinit::meta_data),
        )
        .route(
            "/api/cloudinit/:id/user-data",
            get(waifud::api::cloudinit::user_data),
        )
        .route(
            "/api/cloudinit/:id/vendor-data",
            get(waifud::api::cloudinit::vendor_data),
        )
        .layer(middleware);

    // run it
    let addr = &"[::]:23818".parse()?;
    info!("listening on {}", addr);
    axum::Server::bind(addr)
        .serve(app.into_make_service())
        .await?;

    Ok(())
}

#[instrument(err)]
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
