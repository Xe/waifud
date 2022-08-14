use hyper::{body::Buf, Body, Client, Request, StatusCode};
use hyperlocal::{UnixClientExt, Uri};
use serde::{Deserialize, Serialize};
use std::net::{IpAddr, SocketAddr};

#[derive(Debug, thiserror::Error)]
pub enum Error {
    #[error("hyper error: {0}")]
    Hyper(#[from] hyper::Error),

    #[error("http error: {0}")]
    HTTP(#[from] hyper::http::Error),

    #[error("json error: {0}")]
    JSON(#[from] serde_json::Error),

    #[error("wanted status code {0}, but tailscaled returned status code {1}")]
    WrongStatusCode(StatusCode, StatusCode),
}

#[derive(Clone, PartialEq, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct WhoisResponse {
    #[serde(rename = "Node")]
    pub node: WhoisPeer,
    #[serde(rename = "UserProfile")]
    pub user_profile: User,
}

pub async fn whois(ip_port: SocketAddr) -> Result<WhoisResponse, Error> {
    let ip_port = if let SocketAddr::V6(ip_port) = ip_port {
        let ip = ip_port
            .ip()
            .to_ipv4()
            .map(|ip| IpAddr::V4(ip))
            .unwrap_or(IpAddr::V6(ip_port.ip().clone()));
        (ip, ip_port.port()).into()
    } else {
        ip_port
    };
    let url: hyper::Uri = Uri::new(
        "/var/run/tailscale/tailscaled.sock",
        &format!("/localapi/v0/whois?addr={ip_port}"),
    )
    .into();
    let client = Client::unix();

    let req = Request::builder()
        .uri(url)
        .header("Host", "local-tailscaled.sock")
        .body(Body::empty())
        .unwrap();

    let resp = client.request(req).await?;
    if !resp.status().is_success() {
        return Err(Error::WrongStatusCode(StatusCode::OK, resp.status()));
    }

    let body = hyper::body::aggregate(resp).await?;

    Ok(serde_json::from_reader(body.reader())?)
}

#[derive(Serialize, Deserialize, Clone, PartialEq, Eq)]
pub struct User {
    #[serde(rename = "ID")]
    pub id: u64,
    #[serde(rename = "LoginName")]
    pub login_name: String,
    #[serde(rename = "DisplayName")]
    pub display_name: String,
    #[serde(rename = "ProfilePicURL")]
    pub profile_pic_url: String,
    #[serde(rename = "Roles")]
    pub roles: Vec<Option<serde_json::Value>>,
}

#[derive(Default, Debug, Clone, PartialEq, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct WhoisPeer {
    #[serde(rename = "ID")]
    pub id: i64,
    #[serde(rename = "StableID")]
    pub stable_id: String,
    #[serde(rename = "Name")]
    pub name: String,
    #[serde(rename = "User")]
    pub user: i64,
    #[serde(rename = "Key")]
    pub key: String,
    #[serde(rename = "KeyExpiry")]
    pub key_expiry: String,
    #[serde(rename = "Machine")]
    pub machine: String,
    #[serde(rename = "DiscoKey")]
    pub disco_key: String,
    #[serde(rename = "Addresses")]
    pub addresses: Vec<String>,
    #[serde(rename = "AllowedIPs")]
    pub allowed_ips: Vec<String>,
    #[serde(rename = "Endpoints")]
    pub endpoints: Vec<String>,
    #[serde(rename = "Hostinfo")]
    pub hostinfo: Hostinfo,
    #[serde(rename = "Created")]
    pub created: String,
    #[serde(rename = "PrimaryRoutes")]
    pub primary_routes: Option<Vec<String>>,
    #[serde(rename = "MachineAuthorized")]
    pub machine_authorized: Option<bool>,
    #[serde(rename = "Capabilities")]
    pub capabilities: Option<Vec<String>>,
    #[serde(rename = "ComputedName")]
    pub computed_name: String,
    #[serde(rename = "ComputedNameWithHost")]
    pub computed_name_with_host: String,
}

#[derive(Default, Debug, Clone, PartialEq, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct Hostinfo {
    #[serde(rename = "OS")]
    pub os: Option<String>,
    #[serde(rename = "Hostname")]
    pub hostname: String,
    #[serde(rename = "RoutableIPs")]
    pub routable_ips: Option<Vec<String>>,
    #[serde(rename = "Services")]
    pub services: Vec<Service>,
}

#[derive(Default, Debug, Clone, PartialEq, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct Service {
    #[serde(rename = "Proto")]
    pub proto: String,
    #[serde(rename = "Port")]
    pub port: i64,
    #[serde(rename = "Description")]
    pub description: Option<String>,
}
