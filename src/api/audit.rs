use crate::{models::AuditEvent, Result, State};
use axum::{extract::Extension, Json};
use std::sync::Arc;

#[instrument(err)]
pub async fn list(Extension(state): Extension<Arc<State>>) -> Result<Json<Vec<AuditEvent>>> {
    let conn = state.pool.get().await?;

    Ok(Json(AuditEvent::get_all(&conn)?))
}
