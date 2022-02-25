#[macro_use]
extern crate tracing;

use axum::{
    extract::extractor_middleware,
    routing::{delete, get, post},
    AddExtensionLayer, Router,
};
use std::sync::Arc;
use tower::limit::ConcurrencyLimitLayer;
use tower_http::trace::TraceLayer;
use waifud::{
    api::{self, audit, cloudinit, distros, instances},
    Config, Result, State,
};

#[tokio::main]
async fn main() -> Result {
    tracing_subscriber::fmt::init();

    waifud::migrate::run()?;

    let cfg: Config = serde_dhall::from_file("./config.dhall").parse()?;

    let middleware = tower::ServiceBuilder::new()
        .layer(TraceLayer::new_for_http())
        .layer(ConcurrencyLimitLayer::new(64))
        .layer(AddExtensionLayer::new(Arc::new(
            tailscale_client::Client::new(
                waifud::APPLICATION_NAME.to_string(),
                cfg.tailscale.api_key.clone(),
                cfg.tailscale.tailnet.clone(),
            )?,
        )))
        .layer(AddExtensionLayer::new(Arc::new(State::new().await?)))
        .layer(AddExtensionLayer::new(Arc::new(cfg)));

    let auth_middleware = middleware
        .clone()
        .layer(extractor_middleware::<waifud::paseto::RequireAuth>());

    let api = Router::new()
        .route("/auditlogs", get(audit::list))
        .route("/distros", get(distros::list))
        .route("/distros", post(distros::create))
        .route("/distros/:name", post(distros::update))
        .route("/distros/:name", get(distros::get))
        .route("/distros/:name", delete(distros::delete))
        .route("/instances", post(instances::create))
        .route("/instances", get(instances::list))
        .route("/instances/:id", get(instances::get))
        .route("/instances/:id/reinit", post(instances::reinit))
        .route("/instances/:id/hardreboot", post(instances::hard_reboot))
        .route("/instances/:id/reboot", post(instances::reboot))
        .route("/instances/:id/start", post(instances::start))
        .route("/instances/:id/shutdown", post(instances::shutdown))
        .route("/instances/name/:name", get(instances::get_by_name))
        .route("/instances/:id", delete(instances::delete))
        .route("/instances/:id/machine", get(instances::get_machine))
        .route("/libvirt/machines", get(api::libvirt::get_machines))
        .layer(auth_middleware);

    let app = Router::new()
        .nest("/api/v1", api)
        .route("/api/cloudinit/:id/meta-data", get(cloudinit::meta_data))
        .route("/api/cloudinit/:id/user-data", get(cloudinit::user_data))
        .route(
            "/api/cloudinit/:id/vendor-data",
            get(cloudinit::vendor_data),
        )
        .layer(middleware);

    let addr = &"[::]:23818".parse()?;
    info!("listening on {}", addr);
    axum::Server::bind(addr)
        .serve(app.into_make_service())
        .await?;

    Ok(())
}
