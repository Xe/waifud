use chrono::prelude::*;
use serde::{Deserialize, Serialize};

#[derive(thiserror::Error, Debug)]
pub enum Error {
    #[error("error making request: {0}")]
    Reqwest(#[from] reqwest::Error),

    #[error("error parsing json: {0}")]
    JSON(#[from] serde_json::Error),
}

pub type Result<T = (), E = Error> = std::result::Result<T, E>;

/// The Tailscale API Client. Each call will be its own method.
pub struct Client {
    cli: reqwest::Client,
    api_key: String,
    tailnet: String,
    base_url: String,
}

/// The minimal form of a Tailscale API key, only shows the ID.
#[derive(Debug, Clone, Deserialize, Serialize)]
pub struct Key {
    pub id: String,
}

/// Full information about a Tailscale API key.
#[derive(Debug, Clone, Deserialize, Serialize)]
pub struct KeyInfo {
    pub id: String,
    pub key: Option<String>,
    pub created: DateTime<Utc>,
    pub expires: DateTime<Utc>,
    /// TODO(Xe): make this into a better value
    pub capabilities: serde_json::Value,
}

#[test]
fn test_keyinfo_full() {
    let inp = r#"{
	"id":           "k123456CNTRL",
	"key":          "tskey-k123456CNTRL-abcdefghijklmnopqrstuvwxyz",
	"created":      "2021-12-09T23:22:39Z",
	"expires":      "2022-03-09T23:22:39Z",
	"capabilities": {"devices": {"create": {"reusable": false, "ephemeral": false}}}
}"#;
    let _: KeyInfo = serde_json::from_str(inp).unwrap();
}

#[test]
fn test_keyinfo_partial() {
    let inp = r#"{
	"id":           "k123456CNTRL",
	"created":      "2021-12-09T23:22:39Z",
	"expires":      "2022-03-09T23:22:39Z",
	"capabilities": {"devices": {"create": {"reusable": false, "ephemeral": false}}}
}"#;
    let _: KeyInfo = serde_json::from_str(inp).unwrap();
}

#[derive(Debug, Clone, Deserialize, Serialize)]
pub struct Capabilities {
    pub reusable: bool,
    pub ephemeral: bool,
}

impl Client {
    /// Maybe construct a new client with a given user agent string.
    pub fn new(user_agent: String, api_key: String, tailnet: String) -> Result<Self> {
        let cli = reqwest::Client::builder().user_agent(user_agent).build()?;

        Ok(Client {
            cli,
            api_key,
            tailnet,
            base_url: "https://api.tailscale.com".to_string(),
        })
    }

    /// List all active node authentication keys in the tailnet.
    pub async fn list_keys(&self) -> Result<Vec<Key>> {
        #[derive(Debug, Clone, Deserialize, Serialize)]
        struct Wrapper {
            keys: Vec<Key>,
        }

        let w: Wrapper = self
            .cli
            .get(&format!(
                "{}/api/v2/tailnet/{}/keys",
                self.base_url, self.tailnet
            ))
            .basic_auth(&self.api_key, None::<&String>)
            .send()
            .await?
            .error_for_status()?
            .json()
            .await?;

        Ok(w.keys)
    }

    /// Creates a new machine authkey for the tailnet. This cannot be used
    /// to create API keys.
    pub async fn create_key(&self, caps: Capabilities) -> Result<KeyInfo> {
        #[derive(Debug, Clone, Deserialize, Serialize)]
        struct Outer {
            capabilities: Inner,
        }

        #[derive(Debug, Clone, Deserialize, Serialize)]
        struct Inner {
            devices: Inner2,
        }

        #[derive(Debug, Clone, Deserialize, Serialize)]
        struct Inner2 {
            create: Capabilities,
        }

        let msg = Outer {
            capabilities: Inner {
                devices: Inner2 { create: caps },
            },
        };

        Ok(self
            .cli
            .post(&format!(
                "{}/api/v2/tailnet/{}/keys",
                self.base_url, self.tailnet
            ))
            .basic_auth(&self.api_key, None::<&String>)
            .json(&msg)
            .send()
            .await?
            .error_for_status()?
            .json()
            .await?)
    }

    /// Get detailed information for a given key by ID.
    pub async fn key_info(&self, key_id: String) -> Result<KeyInfo> {
        Ok(self
            .cli
            .get(&format!(
                "{}/api/v2/tailnet/{}/keys/{}",
                self.base_url, self.tailnet, key_id,
            ))
            .basic_auth(&self.api_key, None::<&String>)
            .send()
            .await?
            .error_for_status()?
            .json()
            .await?)
    }

    /// Delete a key by ID.
    pub async fn delete_key(&self, key_id: String) -> Result {
        self.cli
            .delete(&format!(
                "{}/api/v2/tailnet/{}/keys/{}",
                self.base_url, self.tailnet, key_id,
            ))
            .basic_auth(&self.api_key, None::<&String>)
            .send()
            .await?
            .error_for_status()?
            .json()
            .await?;

        Ok(())
    }
}
