#[macro_use]
extern crate tracing;

use axum::{
    routing::{delete, get, get_service, post},
    Extension, Router,
};
use hyper::StatusCode;
use std::{io, net::SocketAddr, sync::Arc};
use tower::limit::ConcurrencyLimitLayer;
use tower_http::{services::ServeDir, trace::TraceLayer};
use waifud::{
    admin,
    api::{self, audit, cloudinit, distros, instances},
    Config, Result, State,
};
use yubico::config::Config as YKConfig;

#[tokio::main]
async fn main() -> Result {
    tracing_subscriber::fmt::init();

    waifud::migrate::run()?;

    let cfg: Config = serde_dhall::from_file("./config.dhall").parse()?;

    let yk: YKConfig = YKConfig::default()
        .set_client_id(cfg.yubikey.client_id.clone())
        .set_key(cfg.yubikey.key.clone());

    let middleware = tower::ServiceBuilder::new()
        .layer(TraceLayer::new_for_http())
        .layer(ConcurrencyLimitLayer::new(64))
        .layer(Extension(Arc::new(tailscale_client::Client::new(
            waifud::APPLICATION_NAME.to_string(),
            cfg.tailscale.api_key.clone(),
            cfg.tailscale.tailnet.clone(),
        )?)))
        .layer(Extension(Arc::new(State::new().await?)))
        .layer(Extension(Arc::new(cfg)))
        .layer(Extension(Arc::new(yk)));

    let admin_panel = Router::new()
        .route("/", get(admin::home))
        .route("/api/config", get(admin::config))
        .route("/test", get(admin::test_handler))
        .route("/instances", get(admin::instances))
        .route("/instances/create", get(admin::instance_create))
        .route("/instances/:id", get(admin::instance))
        .route("/distros", get(admin::distro_list));

    let api = Router::new()
        .route("/auditlogs", get(audit::list))
        .route("/auditlogs/instance/:id", get(audit::list_for_instance))
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
        .layer(middleware.clone());

    let app = Router::new()
        .nest("/api/v1", api)
        .nest("/admin", admin_panel)
        .route("/api/cloudinit/:id/meta-data", get(cloudinit::meta_data))
        .route("/api/cloudinit/:id/user-data", get(cloudinit::user_data))
        .route(
            "/api/cloudinit/:id/vendor-data",
            get(cloudinit::vendor_data),
        )
        .nest(
            "/static",
            get_service(ServeDir::new("./static")).handle_error(|err: io::Error| async move {
                (
                    StatusCode::INTERNAL_SERVER_ERROR,
                    format!("unhandled internal server error: {}", err),
                )
            }),
        )
        .layer(middleware);

    // tokio::spawn(waifud::scrape::cron());

    let addr = &"[::]:23818".parse()?;
    info!("listening on {}", addr);
    axum::Server::bind(addr)
        .serve(app.into_make_service_with_connect_info::<SocketAddr>())
        .await?;

    Ok(())
}
