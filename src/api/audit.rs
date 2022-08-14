use crate::{models::AuditEvent, tailauth::Tailauth, Result, State};
use axum::{extract::Extension, Json};
use std::sync::Arc;

#[instrument(err)]
pub async fn list(
    Extension(state): Extension<Arc<State>>,
    _: Tailauth,
) -> Result<Json<Vec<AuditEvent>>> {
    let conn = state.pool.get().await?;

    Ok(Json(AuditEvent::get_all(&conn)?))
}
