#[macro_use]
extern crate tracing;

use serde::{Deserialize, Serialize};
use serde_dhall::StaticType;
use std::{
    convert::TryInto,
    fs,
    io::{self, Write},
    path::PathBuf,
    process::exit,
    time::Duration,
};
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
    pub host: Option<String>,
    #[structopt(subcommand)]
    cmd: Command,
}

#[derive(StructOpt, Deserialize, Serialize, Debug, StaticType)]
struct Config {
    /// waifud host to connect to
    #[structopt(short, long)]
    pub host: String,
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
    /// Turn an instance on
    Start {
        /// Instance name
        name: String,
    },
    /// Turn an instance off
    Shutdown {
        /// Instance name
        name: String,
    },
    /// Manually trigger instance reboot
    Reboot {
        /// Instance name
        name: String,

        /// Unsafely force reboot
        #[structopt(short, long)]
        hard: bool,
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
    /// Create a new base distro snapshot
    Create(CreateDistroOpts),
    /// Delete a distro image
    Delete { name: String },
    /// List all distros
    List {
        /// Show more information
        #[structopt(short)]
        verbose: bool,
    },
    /// Updates a base distro snapshot
    Update(CreateDistroOpts),
}

/// Defines a base distro snapshot for waifud to use
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

    let mut table = Table::new("{:>}  {:<}  {:<}  {:<}  {:<}  {:<}  {:<}");
    table.add_row(row!(
        "name", "host", "distro", "memory", "ip", "status", "id"
    ));
    for instance in instances {
        let m = cli.get_instance_machine(instance.uuid).await?;

        table.add_row(row!(
            instance.name,
            instance.host,
            instance.distro,
            instance.memory,
            m.addr.unwrap_or("".into()),
            instance.status,
            instance.uuid,
        ));
    }

    println!("{}", table);

    Ok(())
}

async fn wait_until_status<T>(cli: &Client, i: Instance, want: T) -> Result
where
    T: Into<String>,
{
    let want = want.into();
    let mut i = i.clone();

    loop {
        i = cli.get_instance(i.uuid).await?;
        io::stdout().flush()?;
        print!("{}: {}   \r", i.name, i.status);
        if i.status == want {
            break;
        }

        tokio::time::sleep(Duration::from_millis(1000)).await;
    }

    io::stdout().flush()?;
    print!("\n");
    Ok(())
}

async fn start_instance(cli: Client, name: String) -> Result {
    let i = cli.get_instance_by_name(name).await?;

    cli.start_instance(i.uuid).await?;

    wait_until_status(&cli, i.clone(), "running").await?;
    println!("{} is running", i.name);

    Ok(())
}

async fn shutdown_instance(cli: Client, name: String) -> Result {
    let i = cli.get_instance_by_name(name).await?;

    cli.shutdown_instance(i.uuid).await?;

    println!("shut down {}", i.name);

    Ok(())
}

async fn reboot_instance(cli: Client, name: String, hard: bool) -> Result {
    let i = cli.get_instance_by_name(name).await?;

    if hard {
        cli.hard_reboot_instance(i.uuid).await
    } else {
        cli.reboot_instance(i.uuid).await
    }?;

    wait_until_status(&cli, i, "running").await?;

    Ok(())
}

async fn create_instance(cli: Client, opts: CreateOpts) -> Result {
    let ni: NewInstance = opts.try_into()?;
    let i = cli.create_instance(ni).await?;

    println!("created instance {} on {}", i.name, i.host);

    wait_until_status(&cli, i.clone(), "running").await?;

    let m = cli.get_instance_machine(i.uuid).await?;

    println!(
        "\r{}: {}: IP address: {}",
        i.name,
        i.status,
        m.addr.unwrap()
    );

    Ok(())
}

async fn delete_instance(cli: Client, name: String) -> Result {
    let i = cli.get_instance_by_name(name.clone()).await;

    match i {
        Ok(i) => cli.delete_instance(i.uuid).await?,
        Err(why) => {
            eprintln!("no instance named {} was found: {}", name, why);
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

async fn update_distro(cli: Client, opts: CreateDistroOpts) -> Result {
    if let Err(why) = cli.get_distro(opts.name.clone()).await {
        println!("can't get distro {}: {}", opts.name, why);
        exit(1);
    }

    let d: Distro = opts.into();
    let d = cli.update_distro(d).await?;
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
        let mut table = Table::new("{:<}  {:<}");
        table.add_row(row!("name", "disk GB"));
        distros.into_iter().for_each(|d| {
            table.add_row(row!(d.name, d.min_size.to_string()));
        });
        println!("{}", table);
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

    let mut opt = Opt::from_args();
    debug!("{:?}", opt);

    if opt.host.is_none() {
        let mut fname = dirs::config_dir().unwrap();
        fname.push("xeserv");
        let _ = fs::create_dir_all(&fname);
        fname.push("waifuctl");
        fname.set_extension("dhall");

        if let Err(_) = fs::metadata(&fname) {
            let mut fout = fs::File::create(&fname).unwrap();
            let cfg = serde_dhall::serialize(&Config {
                host: "http://[::]:23818".into(),
            })
            .static_type_annotation()
            .to_string()?;
            fout.write_all(cfg.as_bytes())?;
        }

        opt.host = Some(serde_dhall::from_file(&fname).parse::<Config>()?.host);
    }

    let cli = Client::new(opt.host.unwrap())?;

    if let Err(why) = match opt.cmd {
        Command::Distro { cmd } => match cmd {
            DistroCmd::Create(opts) => create_distro(cli, opts).await,
            DistroCmd::Delete { name } => delete_distro(cli, name).await,
            DistroCmd::List { verbose } => list_distros(cli, verbose).await,
            DistroCmd::Update(opts) => update_distro(cli, opts).await,
        },
        Command::List => list_instances(cli).await,
        Command::Create(opts) => create_instance(cli, opts).await,
        Command::Delete { name } => delete_instance(cli, name).await,
        Command::Reboot { name, hard } => reboot_instance(cli, name, hard).await,
        Command::Start { name } => start_instance(cli, name).await,
        Command::Shutdown { name } => shutdown_instance(cli, name).await,
    } {
        eprintln!("OOPSIE WOOPSIE!! Uwu We made a fucky wucky!! A wittle fucko boingo! The code monkeys at our headquarters are working VEWY HAWD to fix this!");
        eprintln!("{}", why);
    }

    Ok(())
}
