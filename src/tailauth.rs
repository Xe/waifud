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

        debug!("remote user: {}", result.user_profile.login_name);

        Ok(Tailauth(result.user_profile, result.node))
    }
}
