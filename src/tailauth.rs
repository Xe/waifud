use crate::Error;
use async_trait::async_trait;
use axum::extract::{FromRequest, RequestParts};

pub struct Tailauth(pub ts_localapi::User, pub ts_localapi::WhoisPeer);

#[async_trait]
impl<B> FromRequest<B> for Tailauth
where
    B: Send,
{
    type Rejection = Error;

    async fn from_request(req: &mut RequestParts<B>) -> Result<Self, Self::Rejection> {
        let addr: axum_client_ip::ClientIp =
            req.extract().await.map_err(|_| Error::BadMiddlewareStack)?;

        let result = ts_localapi::whois((addr.0, 0).into()).await?;

        info!(
            user = result.user_profile.login_name,
            ip = addr.0.to_string(),
            platform = result
                .node
                .clone()
                .hostinfo
                .os
                .unwrap_or("<unknown>".to_string())
        );

        Ok(Tailauth(result.user_profile, result.node))
    }
}
