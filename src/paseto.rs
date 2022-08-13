use crate::{models::Session, Config, Error, State};
use async_trait::async_trait;
use axum::{
    extract::{Extension, FromRequest, RequestParts},
    http::header,
    Json,
};
use paseto::{validate_public_token, PasetoBuilder, PasetoPublicKey};
use ring::signature::Ed25519KeyPair;
use rusqlite::params;
use serde::{Deserialize, Serialize};
use serde_json::Value;
use std::sync::Arc;
use uuid::Uuid;
use yubico::config::Config as YKConfig;

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct LoginRequest {
    pub otp: String,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct TokenSubset {
    pub jti: Uuid,
    pub sub: String,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct LoginResponse {
    pub token: String,
    pub session: Session,
}

#[axum_macros::debug_handler]
pub async fn login(
    Extension(cfg): Extension<Arc<Config>>,
    Extension(state): Extension<Arc<State>>,
    Json(inp): Json<LoginRequest>,
) -> Result<Json<LoginResponse>, Error> {
    let conn = state.pool.get().await?;

    let uid = &inp.otp.clone()[0..12];
    let yk: YKConfig = YKConfig::default()
        .set_client_id(cfg.yubikey.client_id.clone())
        .set_key(cfg.yubikey.key.clone());
    yubico::verify_async(inp.otp, yk).await?;

    let kp = Ed25519KeyPair::from_pkcs8(&hex::decode(&cfg.paseto_keypair)?)?;
    let id = uuid::Uuid::new_v4();

    let token = PasetoBuilder::new()
        .set_ed25519_key(&kp)
        .set_jti(&id.to_string())
        .set_subject(uid.clone())
        .build()
        .map_err(|why| Error::CantMakeToken(format!("{}", why)))?;

    let session = Session {
        uuid: id,
        user: uid.to_string(),
        expired: false,
    };

    conn.execute(
        "INSERT INTO sessions(uuid, user, expired) VALUES (?1, ?2, ?3)",
        params![session.uuid, session.user, session.expired],
    )?;
    conn.execute(
        "INSERT INTO audit_logs(kind, op, data) VALUES (?1, ?2, ?3)",
        params!["session", "create", serde_json::to_value(&session)?],
    )?;

    Ok(Json(LoginResponse { token, session }))
}

/// An extractor that performs authorization.
pub struct RequireAuth;

#[async_trait]
impl<B> FromRequest<B> for RequireAuth
where
    B: Send,
{
    type Rejection = Error;

    async fn from_request(req: &mut RequestParts<B>) -> Result<Self, Self::Rejection> {
        let auth_header = req
            .headers()
            .get(header::AUTHORIZATION)
            .and_then(|value| value.to_str().ok());

        let cfg: &Arc<Config> = req
            .extensions()
            .get::<Arc<Config>>()
            .ok_or(Error::BadMiddlewareStack)?;

        let state: &Arc<State> = req
            .extensions()
            .get::<Arc<State>>()
            .ok_or(Error::BadMiddlewareStack)?;
        let conn = state.pool.get().await?;

        let kp = Ed25519KeyPair::from_pkcs8(&hex::decode(&cfg.paseto_keypair)?)?;

        let tok_data = token_is_valid(kp, auth_header.ok_or(Error::Unauthorized)?)?;

        let ts: TokenSubset = serde_json::from_value(tok_data)?;

        let expired = conn.query_row(
            "SELECT expired FROM sessions WHERE uuid = ?1 LIMIT 1",
            params![ts.jti],
            |row| row.get(0),
        )?;

        if expired {
            Err(Error::Unauthorized)
        } else {
            Ok(RequireAuth)
        }
    }
}

fn token_is_valid(kp: Ed25519KeyPair, token: &str) -> Result<Value, Error> {
    validate_public_token(
        token,
        None,
        &PasetoPublicKey::ED25519KeyPair(&kp),
        &paseto::TimeBackend::Chrono,
    )
    .map_err(|_| Error::Unauthorized)
}
