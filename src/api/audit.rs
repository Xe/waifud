use crate::{models::AuditEvent, tailauth::Tailauth, Result, State};
use axum::{
    extract::{Extension, Path},
    Json,
};
use std::sync::Arc;

#[instrument(err, skip(state))]
pub async fn list(
    Extension(state): Extension<Arc<State>>,
    _: Tailauth,
) -> Result<Json<Vec<AuditEvent>>> {
    let conn = state.pool.get().await?;

    Ok(Json(AuditEvent::get_all(&conn)?))
}

#[instrument(err, skip(state))]
pub async fn list_for_instance(
    Path(id): Path<uuid::Uuid>,
    Extension(state): Extension<Arc<State>>,
    _: Tailauth,
) -> Result<Json<Vec<AuditEvent>>> {
    let conn = state.pool.get().await?;

    Ok(Json(AuditEvent::get_for_instance(id, &conn)?))
}
