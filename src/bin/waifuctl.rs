#![deny(missing_docs)]

//! waifuctl lets you manage VM instances on waifud.

#[macro_use]
extern crate tracing;

use chrono::prelude::*;
use clap::{Args, Parser, Subcommand};
use clap_complete::{generate, Shell};
use serde::{Deserialize, Serialize};
use serde_dhall::StaticType;
use std::{
    convert::TryInto,
    fs,
    io::{self, stdout, Write},
    path::PathBuf,
    process::exit,
    time::Duration,
};
use tabular::{row, Table};
use waifud::{
    client::Client,
    libvirt::NewInstance,
    models::{Distro, Instance},
    Error, Result,
};

#[derive(Debug, Parser)]
#[clap(author, version, about, long_about = None)]
#[clap(propagate_version = true)]
/// waifuctl lets you manage VM instances on waifud.
struct Opt {
    /// waifud host to connect to, formatted as a http/https URL
    #[clap(short = 'H', long)]
    pub host: Option<String>,

    #[clap(subcommand)]
    cmd: Command,
}

#[derive(Deserialize, Serialize, Debug, StaticType, Clone)]
struct Config {
    /// waifud host to connect to, formatted as a http/https URL
    pub host: String,

    /// Default cloudconfig to preload into every VM
    pub userdata: String,
}

#[derive(Subcommand, Debug)]
enum ConfigCmd {
    /// Shows current config
    Show,

    /// Set the waifud host to an arbitrary URL
    SetHost {
        /// The waifud host
        url: String,
    },

    /// Set the default cloudconfig added to instances
    SetUserdata,
}

#[derive(Subcommand, Debug)]
enum Command {
    /// Manage audit logs
    Audit {
        /// Format all audit logs in JSON
        #[clap(long)]
        json: bool,
    },
    /// Manage waifuctl configuration
    Config {
        #[clap(subcommand)]
        cmd: ConfigCmd,
    },
    /// List all instances
    List,
    Create(CreateOpts),
    /// Delete an instance by name
    Delete {
        /// Instance name
        name: String,
    },
    Distro {
        #[clap(subcommand)]
        cmd: DistroCmd,
    },
    /// Reset a VM back to factory settings
    Reinit {
        /// Instance name
        name: String,
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
        #[clap(short, long)]
        hard: bool,
    },
    /// Utilities to help with managing the waifud project
    Utils {
        #[clap(subcommand)]
        cmd: UtilsCmd,
    },
}

/// Create a new instance
#[derive(Args, Debug)]
struct CreateOpts {
    /// Instance name, leave blank to autogenerate
    #[clap(short, long)]
    name: Option<String>,

    /// Memory in megabytes
    #[clap(short, long, default_value = "512")]
    memory: i32,

    /// CPU cores
    #[clap(short, long, default_value = "2")]
    cpus: i32,

    /// Host to put the VM on
    #[clap(short, long)]
    host: String,

    /// Disk size in GB, leave blank to use distribution default
    #[clap(short = 's', long = "disk-size")]
    disk_size: Option<i32>,

    /// ZFS dataset to put the VM disk in
    #[clap(short, long = "zvol", default_value = "rpool/local/vms")]
    zvol_prefix: String,

    /// File containing cloud-init user data, if not set will default to configured value
    #[clap(short, long)]
    user_data: Option<PathBuf>,

    /// Distribution to use
    #[clap(short, long)]
    distro: String,

    /// Automagically join the tailnet
    #[clap(short, long)]
    join_tailnet: bool,
}

impl TryInto<NewInstance> for CreateOpts {
    type Error = anyhow::Error;

    fn try_into(self) -> Result<NewInstance, anyhow::Error> {
        let user_data = match self.user_data {
            Some(user_data) => Some(fs::read_to_string(user_data)?),
            None => None,
        };

        Ok(NewInstance {
            name: self.name,
            memory_mb: Some(self.memory),
            cpus: Some(self.cpus),
            host: self.host,
            disk_size_gb: self.disk_size,
            zvol_prefix: Some(self.zvol_prefix),
            distro: self.distro,
            sata: Some(false),
            user_data,
            join_tailnet: self.join_tailnet,
        })
    }
}

/// Manage distribution images in waifud
#[derive(Subcommand, Debug)]
enum DistroCmd {
    /// Create a new base distro snapshot
    Create(CreateDistroOpts),
    /// Delete a distro image
    Delete { name: String },
    /// List all distros
    List {
        /// Show more information
        #[clap(short)]
        verbose: bool,
    },
    /// Scrapes current versions for distributions
    Scrape,
    /// Updates a base distro snapshot
    Update(CreateDistroOpts),
}

/// Defines a base distro snapshot for waifud to use
#[derive(Args, Debug)]
struct CreateDistroOpts {
    /// Distribution name, include the version as a suffix
    #[clap(short, long)]
    pub name: String,

    /// Download URL for the qcow2 base snapshot
    #[clap(short, long = "download-url")]
    pub download_url: String,

    /// The sha256 of the qcow2 base snapshot
    #[clap(short, long = "sha256")]
    pub sha256sum: String,

    /// The minimum size of a VM created from this snapshot (gigabytes)
    #[clap(short, long)]
    pub min_size: i32,

    /// The format of the disk image
    #[clap(short, long, default_value = "waifud://qcow2")]
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

#[derive(Subcommand, Debug)]
enum UtilsCmd {
    /// Generate shell completions
    Completions {
        #[clap(value_parser)]
        shell: Shell,
    },
    /// Generate manpages to a given folder
    Manpage { path: PathBuf },
}

async fn list_instances(cli: Client) -> Result {
    let instances = cli.list_instances().await?;

    let mut table = Table::new("{:>}  {:<}  {:<}  {:<}  {:<}  {:<}  {:<}");
    table.add_row(row!(
        "name", "host", "distro", "memory", "ip", "status", "id"
    ));
    for instance in instances {
        let m = cli.get_instance_machine(instance.uuid).await;

        table.add_row(row!(
            instance.name,
            instance.host,
            instance.distro,
            instance.memory,
            match m {
                Ok(m) => m.addr.unwrap_or("".into()),
                Err(_) => "".to_string(),
            },
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
        print!(
            "{}: {}                                        \r",
            i.name, i.status
        );
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

#[instrument(ret, level = "debug", err, skip(cli))]
async fn create_instance(cli: Client, cfg: Config, opts: CreateOpts) -> Result {
    let mut ni: NewInstance = opts.try_into()?;

    if ni.user_data.is_none() {
        ni.user_data = Some(cfg.userdata);
    }

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

async fn reinit_instance(cli: Client, name: String) -> Result<()> {
    let i = cli.get_instance_by_name(name.clone()).await?;
    cli.reinit_instance(i.uuid).await?;

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

async fn scrape_distros(cli: Client) -> Result {
    let distros = waifud::scrape::get_all().await?;
    for distro in distros {
        cli.update_distro(distro.clone()).await?;
        println!("updated {}", distro.name);
    }

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

async fn audit_list(cli: Client, json: bool) -> Result<()> {
    let logs = cli.audit_logs().await?;

    if json {
        serde_json::to_writer(stdout(), &logs)?;
        return Ok(());
    }

    let mut table = Table::new("{:>}  {:<}  {:<}  {:<}");
    table.add_row(row!("timestamp", "kind", "name", "op"));

    for log in logs {
        let ts = NaiveDateTime::from_timestamp(log.ts, 0);
        table.add_row(row!(
            ts.to_string(),
            log.kind,
            log.name.unwrap_or("".into()),
            log.op
        ));
    }

    println!("{}", table);

    Ok(())
}

fn config_show(cfg: Config) -> Result {
    println!("waifud host: {}", cfg.host);
    println!("default cloudconfig:\n\n{}", cfg.userdata);

    Ok(())
}

fn config_set_host(cfg: Config, url: String) -> Result {
    let mut cfg = cfg.clone();

    cfg.host = url.clone();

    let mut fname = dirs::config_dir().unwrap();
    fname.push("xeserv");
    let _ = fs::create_dir_all(&fname);
    fname.push("waifuctl");
    fname.set_extension("dhall");

    let mut fout = fs::File::create(&fname).unwrap();
    let cfg = serde_dhall::serialize(&cfg)
        .static_type_annotation()
        .to_string()?;
    fout.write_all(cfg.as_bytes())?;

    println!("set host to {} in {}", url, fname.to_str().unwrap());
    Ok(())
}

fn config_set_userdata(cfg: Config) -> Result {
    let userdata = edit::edit(&cfg.userdata)?;
    let mut cfg = cfg.clone();
    cfg.userdata = userdata;

    let mut fname = dirs::config_dir().unwrap();
    fname.push("xeserv");
    let _ = fs::create_dir_all(&fname);
    fname.push("waifuctl");
    fname.set_extension("dhall");

    let mut fout = fs::File::create(&fname).unwrap();
    let cfg = serde_dhall::serialize(&cfg)
        .static_type_annotation()
        .to_string()?;
    fout.write_all(cfg.as_bytes())?;

    println!("wrote default cloudconfig to {}", fname.to_str().unwrap());

    Ok(())
}

fn utils_completions(shell: Shell) -> Result {
    let cmd = clap::Command::new("waifuctl");
    let mut cmd = Opt::augment_args(cmd);

    generate(
        shell,
        &mut cmd,
        "waifuctl".to_string(),
        &mut std::io::stdout(),
    );

    Ok(())
}

fn utils_gen_manpage(path: PathBuf) -> Result {
    let cmd = clap::Command::new("waifuctl");
    let cmd = Opt::augment_args(cmd);

    let man = clap_mangen::Man::new(cmd.clone());
    let mut buffer: Vec<u8> = Default::default();
    man.render(&mut buffer)?;
    std::fs::write(path.join("waifuctl.1"), buffer)?;

    for scmd in cmd.get_subcommands() {
        let man = clap_mangen::Man::new(scmd.clone());
        let mut buffer: Vec<u8> = Default::default();
        man.render(&mut buffer)?;

        std::fs::write(
            path.join(&format!("waifuctl-{}.1", scmd.get_name())),
            buffer,
        )?;

        if scmd.has_subcommands() {
            for sscmd in scmd.get_subcommands() {
                let man = clap_mangen::Man::new(sscmd.clone());
                let mut buffer: Vec<u8> = Default::default();
                man.render(&mut buffer)?;

                std::fs::write(
                    path.join(&format!(
                        "waifuctl-{}-{}.1",
                        scmd.get_name(),
                        sscmd.get_name()
                    )),
                    buffer,
                )?;
            }
        }
    }

    Ok(())
}

#[tokio::main]
async fn main() -> Result<()> {
    tracing_subscriber::fmt::init();

    let mut opt = Opt::parse();

    let cfg = {
        let mut fname = dirs::config_dir().unwrap();
        fname.push("xeserv");
        let _ = fs::create_dir_all(&fname);
        fname.push("waifuctl");
        fname.set_extension("dhall");

        if opt.host.is_none() {
            if let Err(_) = fs::metadata(&fname) {
                let mut fout = fs::File::create(&fname).unwrap();
                let cfg = serde_dhall::serialize(&Config {
                    host: "http://[::]:23818".into(),
                    userdata: include_str!("../../var/base.yaml").to_string(),
                })
                .static_type_annotation()
                .to_string()?;
                fout.write_all(cfg.as_bytes())?;
            }
        }

        let cfg = serde_dhall::from_file(&fname).parse::<Config>()?;
        debug!("config: {:?}", cfg);

        if cfg.host.len() == 0 {
            println!("welcome to waifud, you may want to run `waifuctl config set-host` to point waifuctl to your waifud server");
        }

        cfg
    };
    if let None = opt.host {
        opt.host = Some(cfg.host.clone());
    }

    debug!("{:?}", opt);

    let cli = Client::new(opt.host.unwrap())?;

    if let Err(why) = match opt.cmd {
        Command::Audit { json } => audit_list(cli, json).await,
        Command::Distro { cmd } => match cmd {
            DistroCmd::Create(opts) => create_distro(cli, opts).await,
            DistroCmd::Delete { name } => delete_distro(cli, name).await,
            DistroCmd::List { verbose } => list_distros(cli, verbose).await,
            DistroCmd::Scrape => scrape_distros(cli).await,
            DistroCmd::Update(opts) => update_distro(cli, opts).await,
        },
        Command::List => list_instances(cli).await,
        Command::Create(opts) => create_instance(cli, cfg, opts).await,
        Command::Delete { name } => delete_instance(cli, name).await,
        Command::Reboot { name, hard } => reboot_instance(cli, name, hard).await,
        Command::Reinit { name } => reinit_instance(cli, name).await,
        Command::Start { name } => start_instance(cli, name).await,
        Command::Shutdown { name } => shutdown_instance(cli, name).await,
        Command::Config { cmd } => match cmd {
            ConfigCmd::Show => config_show(cfg),
            ConfigCmd::SetHost { url } => config_set_host(cfg, url),
            ConfigCmd::SetUserdata => config_set_userdata(cfg),
        },
        Command::Utils { cmd } => match cmd {
            UtilsCmd::Completions { shell } => utils_completions(shell),
            UtilsCmd::Manpage { path } => utils_gen_manpage(path),
        },
    } {
        eprintln!("OOPSIE WOOPSIE!! Uwu We made a fucky wucky!! A wittle fucko boingo! The code monkeys at our headquarters are working VEWY HAWD to fix this!");
        eprintln!("{}", why);
    }

    Ok(())
}
