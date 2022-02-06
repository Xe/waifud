#[macro_use]
extern crate tracing;

use std::{convert::TryInto, fs, path::PathBuf, time::Duration};
use structopt::StructOpt;
use tabular::{row, Table};
use waifud::{
    client::Client,
    libvirt::NewInstance,
    models::{Distro, Instance},
    Error, Result,
};

#[derive(StructOpt, Debug)]
/// waifuctl lets you manage VM instances on waifud.
struct Opt {
    /// waifud host to connect to
    #[structopt(short, long)]
    host: String,
    #[structopt(subcommand)]
    cmd: Command,
}

#[derive(StructOpt, Debug)]
enum Command {
    /// List all instances
    List,
    Create(CreateOpts),
    /// Delete an instance by name
    Delete {
        /// Instance name
        name: String,
    },
    Distro {
        #[structopt(subcommand)]
        cmd: DistroCmd,
    },
}

/// Create a new instance
#[derive(StructOpt, Debug)]
struct CreateOpts {
    /// Instance name, leave blank to autogenerate
    #[structopt(short, long)]
    name: Option<String>,

    /// Memory in megabytes
    #[structopt(short, long, default_value = "512")]
    memory: i32,

    /// CPU cores
    #[structopt(short, long, default_value = "2")]
    cpus: i32,

    /// Host to put the VM on
    #[structopt(short, long)]
    host: String,

    /// Disk size in GB, leave blank to use distribution default
    #[structopt(short = "s", long = "disk-size")]
    disk_size: Option<i32>,

    /// ZFS dataset to put the VM disk in
    #[structopt(short, long = "zvol", default_value = "rpool/local/vms")]
    zvol_prefix: String,

    /// File containing cloud-init user data
    #[structopt(short, long, default_value = "./var/xe-base.yaml")]
    user_data: PathBuf,

    /// Distribution to use
    #[structopt(short, long)]
    distro: String,
}

impl TryInto<NewInstance> for CreateOpts {
    type Error = anyhow::Error;

    fn try_into(self) -> Result<NewInstance, anyhow::Error> {
        let user_data = fs::read_to_string(self.user_data)?;

        Ok(NewInstance {
            name: self.name,
            memory_mb: Some(self.memory),
            cpus: Some(self.cpus),
            host: self.host,
            disk_size_gb: self.disk_size,
            zvol_prefix: Some(self.zvol_prefix),
            distro: self.distro,
            sata: Some(false),
            user_data: Some(user_data),
        })
    }
}

/// Manage distribution images in waifud
#[derive(StructOpt, Debug)]
enum DistroCmd {
    /// Create a new distro from details you provide
    Create(CreateDistroOpts),
    /// List all distros
    List {
        /// Show more information
        #[structopt(short)]
        verbose: bool,
    },
    /// Delete a distro image
    Delete { name: String },
}

/// Create a new distro snapshot for waifud to use
#[derive(StructOpt, Debug)]
struct CreateDistroOpts {
    /// Distribution name, include the version as a suffix
    #[structopt(short, long)]
    pub name: String,

    /// Download URL for the qcow2 base snapshot
    #[structopt(short, long = "download-url")]
    pub download_url: String,

    /// The sha256 of the qcow2 base snapshot
    #[structopt(short, long = "sha256")]
    pub sha256sum: String,

    /// The minimum size of a VM created from this snapshot (gigabytes)
    #[structopt(short, long)]
    pub min_size: i32,

    /// The format of the disk image
    #[structopt(short, long, default_value = "waifud://qcow2")]
    pub format: String,
}

impl Into<Distro> for CreateDistroOpts {
    fn into(self) -> Distro {
        Distro {
            name: self.name,
            download_url: self.download_url,
            sha256sum: self.sha256sum,
            min_size: self.min_size,
            format: self.format,
        }
    }
}

async fn list_instances(cli: Client) -> Result {
    let instances = cli.list_instances().await?;

    let mut table = Table::new("{:>}  {:<}  {:<}  {:<}  {:<}");
    table.add_row(row!("name", "host", "memory", "ip", "id"));
    for instance in instances {
        let m = cli.get_instance_machine(instance.uuid).await?;

        table.add_row(row!(
            instance.name,
            instance.host,
            instance.memory,
            m.addr.unwrap_or("".into()),
            instance.uuid,
        ));
    }

    println!("{}", table);

    Ok(())
}

async fn create_instance(cli: Client, opts: CreateOpts) -> Result {
    let ni: NewInstance = opts.try_into()?;
    let i = cli.create_instance(ni).await?;

    println!(
        "created instance {} on {}, waiting for IP address",
        i.name, i.host
    );

    loop {
        let m = cli.get_instance_machine(i.uuid).await?;
        if m.addr.is_none() {
            tokio::time::sleep(Duration::from_millis(1000)).await;
            continue;
        }

        println!("IP address: {}", m.addr.unwrap());
        break;
    }

    Ok(())
}

async fn delete_instance(cli: Client, name: String) -> Result {
    let instances = cli.list_instances().await?;
    let instances = instances
        .into_iter()
        .filter(|i| i.name == name)
        .collect::<Vec<Instance>>();
    let i = instances.get(0);

    match i {
        Some(i) => cli.delete_instance(i.uuid).await?,
        None => {
            eprintln!("no instance named {} was found", name);
            return Err(Error::InstanceDoesntExist(name));
        }
    };

    Ok(())
}

async fn create_distro(cli: Client, opts: CreateDistroOpts) -> Result {
    let d: Distro = opts.into();
    let d = cli.create_distro(d).await?;
    println!("created {}", d.name);

    Ok(())
}

async fn list_distros(cli: Client, verbose: bool) -> Result {
    let distros = cli.list_distros().await?;

    if verbose {
        let mut table = Table::new("{:>}  {:<}  {:<}  {:<}");
        table.add_row(row!("name", "min size", "sha256", "url"));
        for distro in distros {
            table.add_row(row!(
                distro.name,
                distro.min_size,
                distro.sha256sum,
                distro.download_url,
            ));
        }

        println!("{}", table);
    } else {
        distros.iter().for_each(|d| println!("{}", d.name));
    }

    Ok(())
}

async fn delete_distro(cli: Client, name: String) -> Result<()> {
    cli.delete_distro(name).await?;
    Ok(())
}

#[tokio::main]
async fn main() -> Result<()> {
    tracing_subscriber::fmt::init();

    let opt = Opt::from_args();
    debug!("{:?}", opt);

    let cli = Client::new(opt.host)?;

    if let Err(why) = match opt.cmd {
        Command::Distro { cmd } => match cmd {
            DistroCmd::Create(opts) => create_distro(cli, opts).await,
            DistroCmd::Delete { name } => delete_distro(cli, name).await,
            DistroCmd::List { verbose } => list_distros(cli, verbose).await,
        },
        Command::List => list_instances(cli).await,
        Command::Create(opts) => create_instance(cli, opts).await,
        Command::Delete { name } => delete_instance(cli, name).await,
    } {
        eprintln!("OOPSIE WOOPSIE!! Uwu We made a fucky wucky!! A wittle fucko boingo! The code monkeys at our headquarters are working VEWY HAWD to fix this!");
        eprintln!("{}", why);
    }

    Ok(())
}
