#[macro_use]
extern crate tracing;

use axum::{
    routing::{delete, get, post},
    AddExtensionLayer, Router,
};
use std::sync::Arc;
use tower::limit::ConcurrencyLimitLayer;
use tower_http::trace::TraceLayer;
use waifud::{Config, Result, State};

#[tokio::main]
async fn main() -> Result {
    tracing_subscriber::fmt::init();

    waifud::migrate::run()?;

    let cfg: Config = serde_dhall::from_file("./config.dhall").parse()?;

    let middleware = tower::ServiceBuilder::new()
        .layer(TraceLayer::new_for_http())
        .layer(ConcurrencyLimitLayer::new(64))
        .layer(AddExtensionLayer::new(Arc::new(State::new().await?)))
        .layer(AddExtensionLayer::new(Arc::new(cfg)));

    let app = Router::new()
        .route("/api/v1/distros", get(waifud::api::distros::list))
        .route("/api/v1/distros", post(waifud::api::distros::create))
        .route("/api/v1/distros/:name", post(waifud::api::distros::update))
        .route("/api/v1/distros/:name", get(waifud::api::distros::get))
        .route(
            "/api/v1/distros/:name",
            delete(waifud::api::distros::delete),
        )
        .route("/api/v1/instances", post(waifud::api::instances::create))
        .route("/api/v1/instances", get(waifud::api::instances::list))
        .route("/api/v1/instances/:id", get(waifud::api::instances::get))
        .route(
            "/api/v1/instances/:id/reboot",
            post(waifud::api::instances::reboot),
        )
        .route(
            "/api/v1/instances/:id/start",
            post(waifud::api::instances::start),
        )
        .route(
            "/api/v1/instances/:id/shutdown",
            post(waifud::api::instances::shutdown),
        )
        .route(
            "/api/v1/instances/:id/hardreboot",
            post(waifud::api::instances::hard_reboot),
        )
        .route(
            "/api/v1/instances/name/:name",
            get(waifud::api::instances::get_by_name),
        )
        .route(
            "/api/v1/instances/:id",
            delete(waifud::api::instances::delete),
        )
        .route(
            "/api/v1/instances/:id/machine",
            get(waifud::api::instances::get_machine),
        )
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

    let addr = &"[::]:23818".parse()?;
    info!("listening on {}", addr);
    axum::Server::bind(addr)
        .serve(app.into_make_service())
        .await?;

    Ok(())
}
