use crate::{
    api::libvirt::Machine,
    libvirt::NewInstance,
    models::{Distro, Instance},
    Result,
};
use std::time::Duration;
use url::Url;
use uuid::Uuid;

pub struct Client {
    base_url: Url,
    cli: reqwest::Client,
}

impl Client {
    pub fn new(base_url: String) -> Result<Self> {
        let cli = reqwest::Client::builder()
            .connect_timeout(Duration::from_millis(500))
            .build()?;

        Ok(Client {
            base_url: Url::parse(&base_url)?,
            cli,
        })
    }

    pub async fn create_instance(&self, ni: NewInstance) -> Result<Instance> {
        let mut u = self.base_url.clone();

        u.set_path("/api/v1/instances");
        let i: Instance = self
            .cli
            .post(u)
            .json(&ni)
            .send()
            .await?
            .error_for_status()?
            .json()
            .await?;

        Ok(i)
    }

    pub async fn delete_instance(&self, id: Uuid) -> Result {
        let mut u = self.base_url.clone();
        u.set_path(&format!("/api/v1/instances/{}", id));
        self.cli.delete(u).send().await?.error_for_status()?;
        Ok(())
    }

    pub async fn list_instances(&self) -> Result<Vec<Instance>> {
        let mut u = self.base_url.clone();
        u.set_path("/api/v1/instances");
        let instances: Vec<Instance> = self
            .cli
            .get(u)
            .send()
            .await?
            .error_for_status()?
            .json()
            .await?;
        Ok(instances)
    }

    pub async fn get_instance(&self, id: Uuid) -> Result<Instance> {
        let mut u = self.base_url.clone();
        u.set_path(&format!("/api/v1/instances/{}", id));
        let i: Instance = self
            .cli
            .get(u)
            .send()
            .await?
            .error_for_status()?
            .json()
            .await?;
        Ok(i)
    }

    pub async fn get_instance_machine(&self, id: Uuid) -> Result<Machine> {
        let mut u = self.base_url.clone();
        u.set_path(&format!("/api/v1/instances/{}/machine", id));
        let m: Machine = self
            .cli
            .get(u)
            .send()
            .await?
            .error_for_status()?
            .json()
            .await?;
        Ok(m)
    }

    pub async fn create_distro(&self, d: Distro) -> Result<Distro> {
        let mut u = self.base_url.clone();
        u.set_path("/api/v1/distros");
        let d: Distro = self
            .cli
            .post(u)
            .json(&d)
            .send()
            .await?
            .error_for_status()?
            .json()
            .await?;
        Ok(d)
    }

    pub async fn list_distros(&self) -> Result<Vec<Distro>> {
        let mut u = self.base_url.clone();
        u.set_path("/api/v1/distros");
        let distros: Vec<Distro> = self
            .cli
            .get(u)
            .send()
            .await?
            .error_for_status()?
            .json()
            .await?;
        Ok(distros)
    }

    pub async fn delete_distro(&self, name: String) -> Result {
        let mut u = self.base_url.clone();
        u.set_path(&format!("/api/v1/distros/{}", name));
        self.cli.delete(u).send().await?.error_for_status()?;
        Ok(())
    }
}
