use crate::{
    api::libvirt::Machine,
    libvirt::NewInstance,
    models::{AuditEvent, Distro, Instance},
    Result,
};
use reqwest::header;
use std::time::Duration;
use url::Url;
use uuid::Uuid;

pub struct Client {
    base_url: Url,
    cli: reqwest::Client,
}

impl Client {
    pub fn new(base_url: String) -> Result<Self> {
        let mut headers = header::HeaderMap::new();
        headers.insert(
            header::USER_AGENT,
            header::HeaderValue::from_str(crate::APPLICATION_NAME)?,
        );

        let cli = reqwest::Client::builder()
            .default_headers(headers)
            .connect_timeout(Duration::from_millis(500))
            .build()?;

        Ok(Client {
            base_url: Url::parse(&base_url)?,
            cli,
        })
    }

    pub async fn audit_logs(&self) -> Result<Vec<AuditEvent>> {
        let mut u = self.base_url.clone();
        u.set_path("/api/v1/auditlogs");
        Ok(self
            .cli
            .get(u)
            .send()
            .await?
            .error_for_status()?
            .json()
            .await?)
    }

    pub async fn create_instance(&self, ni: NewInstance) -> Result<Instance> {
        let mut u = self.base_url.clone();
        u.set_path("/api/v1/instances");
        Ok(self
            .cli
            .post(u)
            .json(&ni)
            .send()
            .await?
            .error_for_status()?
            .json()
            .await?)
    }

    pub async fn delete_instance(&self, id: Uuid) -> Result {
        let mut u = self.base_url.clone();
        u.set_path(&format!("/api/v1/instances/{}", id));
        self.cli.delete(u).send().await?.error_for_status()?;
        Ok(())
    }

    pub async fn reinit_instance(&self, id: Uuid) -> Result {
        let mut u = self.base_url.clone();
        u.set_path(&format!("/api/v1/instances/{}/reinit", id));
        self.cli.post(u).send().await?.error_for_status()?;
        Ok(())
    }

    pub async fn list_instances(&self) -> Result<Vec<Instance>> {
        let mut u = self.base_url.clone();
        u.set_path("/api/v1/instances");
        Ok(self
            .cli
            .get(u)
            .send()
            .await?
            .error_for_status()?
            .json()
            .await?)
    }

    pub async fn get_instance(&self, id: Uuid) -> Result<Instance> {
        let mut u = self.base_url.clone();
        u.set_path(&format!("/api/v1/instances/{}", id));
        Ok(self
            .cli
            .get(u)
            .send()
            .await?
            .error_for_status()?
            .json()
            .await?)
    }

    pub async fn get_instance_by_name(&self, name: String) -> Result<Instance> {
        let mut u = self.base_url.clone();
        u.set_path(&format!("/api/v1/instances/name/{}", name));
        Ok(self
            .cli
            .get(u)
            .send()
            .await?
            .error_for_status()?
            .json()
            .await?)
    }

    pub async fn get_instance_machine(&self, id: Uuid) -> Result<Machine> {
        let mut u = self.base_url.clone();
        u.set_path(&format!("/api/v1/instances/{}/machine", id));
        Ok(self
            .cli
            .get(u)
            .send()
            .await?
            .error_for_status()?
            .json()
            .await?)
    }

    pub async fn shutdown_instance(&self, id: Uuid) -> Result<()> {
        let mut u = self.base_url.clone();
        u.set_path(&format!("/api/v1/instances/{}/shutdown", id));
        self.cli.post(u).send().await?.error_for_status()?;
        Ok(())
    }

    pub async fn start_instance(&self, id: Uuid) -> Result<()> {
        let mut u = self.base_url.clone();
        u.set_path(&format!("/api/v1/instances/{}/start", id));
        self.cli.post(u).send().await?.error_for_status()?;
        Ok(())
    }

    pub async fn hard_reboot_instance(&self, id: Uuid) -> Result<()> {
        let mut u = self.base_url.clone();
        u.set_path(&format!("/api/v1/instances/{}/hardreboot", id));
        self.cli.post(u).send().await?.error_for_status()?;
        Ok(())
    }

    pub async fn reboot_instance(&self, id: Uuid) -> Result<()> {
        let mut u = self.base_url.clone();
        u.set_path(&format!("/api/v1/instances/{}/reboot", id));
        self.cli.post(u).send().await?.error_for_status()?;
        Ok(())
    }

    pub async fn create_distro(&self, d: Distro) -> Result<Distro> {
        let mut u = self.base_url.clone();
        u.set_path("/api/v1/distros");
        Ok(self
            .cli
            .post(u)
            .json(&d)
            .send()
            .await?
            .error_for_status()?
            .json()
            .await?)
    }

    pub async fn list_distros(&self) -> Result<Vec<Distro>> {
        let mut u = self.base_url.clone();
        u.set_path("/api/v1/distros");
        Ok(self
            .cli
            .get(u)
            .send()
            .await?
            .error_for_status()?
            .json()
            .await?)
    }

    pub async fn update_distro(&self, d: Distro) -> Result<Distro> {
        let mut u = self.base_url.clone();
        u.set_path(&format!("/api/v1/distros/{}", d.name));
        Ok(self
            .cli
            .post(u)
            .json(&d)
            .send()
            .await?
            .error_for_status()?
            .json()
            .await?)
    }

    pub async fn get_distro(&self, name: String) -> Result<Distro> {
        let mut u = self.base_url.clone();
        u.set_path(&format!("/api/v1/distros/{}", name));
        Ok(self
            .cli
            .get(u)
            .send()
            .await?
            .error_for_status()?
            .json()
            .await?)
    }

    pub async fn delete_distro(&self, name: String) -> Result {
        let mut u = self.base_url.clone();
        u.set_path(&format!("/api/v1/distros/{}", name));
        self.cli.delete(u).send().await?.error_for_status()?;
        Ok(())
    }
}
