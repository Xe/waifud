use async_trait::async_trait;
use axum::{
    extract::{FromRequest, RequestParts},
    http::header,
};
use paseto::{validate_public_token, PasetoPublicKey};
use ring::signature::Ed25519KeyPair;
use serde_json::Value;
use std::sync::Arc;

use crate::{Config, Error};

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
            .and_then(|headers| headers.get(header::AUTHORIZATION))
            .and_then(|value| value.to_str().ok());

        let cfg: &Arc<Config> = req
            .extensions()
            .ok_or(Error::BadMiddlewareStack)?
            .get::<Arc<Config>>()
            .ok_or(Error::BadMiddlewareStack)?;
        let kp = Ed25519KeyPair::from_pkcs8(&hex::decode(&cfg.paseto_keypair)?)?;

        #[cfg(debug_assertions)]
        println!(
            "{}",
            paseto::PasetoBuilder::new()
                .set_ed25519_key(&kp)
                .set_claim("test", Value::String("value".to_string()))
                .build()
                .unwrap()
        );

        match auth_header {
            Some(auth_header) if token_is_valid(kp, auth_header) => Ok(Self),
            _ => Err(Error::Unauthorized),
        }
    }
}

fn token_is_valid(kp: Ed25519KeyPair, token: &str) -> bool {
    validate_public_token(
        token,
        None,
        &PasetoPublicKey::ED25519KeyPair(&kp),
        &paseto::TimeBackend::Chrono,
    )
    .is_ok()
}
