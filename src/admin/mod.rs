use crate::{
    api::libvirt::Machine,
    models::{Distro, Instance},
    tailauth::Tailauth,
    Result, State,
};
use axum::{extract::Path, Extension};
use maud::{html, Markup, PreEscaped};
use rusqlite::params;
use std::sync::Arc;
use ts_localapi::User;
use uuid::Uuid;
use virt::{connect::Connect, domain::Domain};

const CSS: PreEscaped<&'static str> = PreEscaped(include_str!("./xess.css"));

fn import_js(name: &str) -> PreEscaped<String> {
    PreEscaped(format!(
        "<script src=\"/static/{name}\" type =\"module\"></script>"
    ))
}

pub fn base(title: Option<String>, user_data: User, body: Markup) -> Markup {
    let page_title = title.clone().unwrap_or("waifud".to_string());
    let title = title
        .map(|s| format!("{s} - waifud"))
        .unwrap_or("waifud".to_string());

    html! {
        (maud::DOCTYPE)
        html {
            head {
                meta charset="utf-8";
                title {(title)}
                style {(CSS)}
                meta name="viewport" content="width=device-width, initial-scale=1.0";
                link rel="icon" href="data:image/svg+xml,<svg xmlns=%22http://www.w3.org/2000/svg%22 viewBox=%220 0 100 100%22><text y=%22.9em%22 font-size=%2290%22>🔥</text></svg>";
            }
            body.top {
                main {
                    nav.nav {
                        a href="/admin" {"Home"}
                        " "
                        a href="/admin/distros" {"Distros"}
                        " "
                        a href="/admin/instances" {"Instances"}
                        div.right {
                            {(user_data.display_name)}
                            " "
                            img style="width:32px;height:32px" src=(user_data.profile_pic_url);
                        }
                    }
                    br;
                    h1 {(page_title)}
                    br;
                    (body);
                    hr;
                    footer {
                        p {
                            "Powered with dokis by "
                            a href="https://github.com/Xe/waifud" {"waifud"}
                            ". ❤️"
                        }
                    }
                }
            }
        }
    }
}

pub async fn instance_create(Tailauth(user, _): Tailauth) -> Markup {
    base(
        Some("Create instance".to_string()),
        user,
        html! {
            (import_js("instance_create.js"))
            div #root {
                "Loading..."
            }
        },
    )
}

pub async fn instance(
    Extension(state): Extension<Arc<State>>,
    Tailauth(user, _): Tailauth,
    Path(id): Path<Uuid>,
) -> Result<Markup> {
    let conn = state.pool.get().await?;

    let instance = Instance::from_uuid(&conn, id)?;

    let conn = Connect::open(&format!("qemu+ssh://root@{}/system", instance.host))?;
    let machine: Option<Machine> = Domain::lookup_by_uuid_string(&conn, &id.to_string())
        .ok()
        .and_then(|dom| Machine::try_from(dom).ok());

    Ok(base(
        Some(instance.name.clone()),
        user,
        html! {
            (import_js("instance_detail.js"))
            table {
                tr {
                    th {"Status"}
                    td {(instance.status)}
                }
                tr {
                    th {"IP Address"}
                    td {
                        @if let Some(m) = machine {
                            (m.addr.unwrap_or("".to_string()))
                        }
                    }
                }
                tr {
                    th {"Host"}
                    td {(instance.host)}
                }
                tr {
                    th {"Memory"}
                    td {(instance.memory) " MB"}
                }
                tr {
                    th {"Disk size"}
                    td {(instance.disk_size) " GB"}
                }
                tr {
                    th {"ZVol name"}
                    td {(instance.zvol_name)}
                }
                tr {
                    th {"Distro"}
                    td {(instance.distro)}
                }
                tr {
                    th {"UUID"}
                    td #instance_id {(instance.uuid.to_string())}
                }
            }

            h2 {"Quick Actions"}
            div #actions {"Loading..."}
        },
    ))
}

pub async fn instances(
    Extension(state): Extension<Arc<State>>,
    Tailauth(user, _): Tailauth,
) -> Result<Markup> {
    let conn = state.pool.get().await?;

    let mut result: Vec<Instance> = Vec::new();

    let mut stmt = conn.prepare(
        "SELECT uuid, name, host, mac_address, memory, disk_size, zvol_name, status, distro, join_tailnet FROM instances",
    )?;
    let instances = stmt.query_map(params![], |row| {
        Ok(Instance {
            uuid: row.get(0)?,
            name: row.get(1)?,
            host: row.get(2)?,
            mac_address: row.get(3)?,
            memory: row.get(4)?,
            disk_size: row.get(5)?,
            zvol_name: row.get(6)?,
            status: row.get(7)?,
            distro: row.get(8)?,
            join_tailnet: row.get(9)?,
        })
    })?;

    for instance in instances {
        result.push(instance?);
    }

    Ok(base(
        Some("Instances".to_string()),
        user,
        html! {
            p{ a href="/admin/instances/create" {"Create a new instance"} }
            table {
                tr {
                    th {"Name"}
                    th {"Host"}
                    th {"Memory"}
                    th {"Disk"}
                    th {"Distro"}
                    th {"Status"}
                }
                @for i in result {
                    tr {
                        td {a href={"/admin/instances/" (i.uuid.to_string())} {(i.name)}}
                        td {(i.host)}
                        td {(i.memory) " MB"}
                        td {(i.disk_size) " GB"}
                        td {(i.distro)}
                        td {(i.status)}
                    }
                }
            }
        },
    ))
}

pub async fn home(
    Extension(state): Extension<Arc<State>>,
    Tailauth(user, _): Tailauth,
) -> Result<Markup> {
    let conn = state.pool.get().await?;

    let mut stmt = conn.prepare(
        "
WITH distro_count    ( val ) AS ( SELECT COUNT(*) FROM distros )
   , instance_count  ( val ) AS ( SELECT COUNT(*) FROM instances )
   , instance_memory ( amt ) AS ( SELECT SUM(memory) FROM instances )
SELECT dc.val AS distros
     , ic.val AS instances
     , im.amt AS ram_use
FROM distro_count    dc
   , instance_count  ic
   , instance_memory im
",
    )?;

    let (distro_count, instance_count, total_memory): (i32, i32, Option<i32>) =
        stmt.query_row(params![], |row| Ok((row.get(0)?, row.get(1)?, row.get(2)?)))?;

    Ok(base(
        Some("Home".to_string()),
        user.clone(),
        html! {
            p {
                "Hello "
                (user.login_name)
                "! I am tracking "
                (distro_count)
                " distribution image"
                @if distro_count != 1 {
                    "s"
                }
                ", "
                (instance_count)
                " VM instance"
                @if instance_count != 1 {
                    "s"
                }
                 " that use a total of "
                (total_memory.unwrap_or(0))
                " megabytes of RAM."
            }
            p{ a href="/admin/instances/create" {"Create a new instance"} }
        },
    ))
}

pub async fn distro_list(
    Extension(state): Extension<Arc<State>>,
    Tailauth(user, _): Tailauth,
) -> Result<Markup> {
    let conn = state.pool.get().await?;

    let mut stmt = conn.prepare(
        "SELECT name, download_url, sha256sum, min_size, format FROM distros ORDER BY name ASC",
    )?;
    let iter = stmt.query_map(params![], |row| {
        Ok(Distro {
            name: row.get(0)?,
            download_url: row.get(1)?,
            sha256sum: row.get(2)?,
            min_size: row.get(3)?,
            format: row.get(4)?,
        })
    })?;
    let mut result: Vec<Distro> = vec![];

    for distro in iter {
        result.push(distro.unwrap());
    }

    Ok(base(
        Some("Distros".to_string()),
        user,
        html! {
            table {
                tr {
                    th {"Name"}
                    th {"Min. Size (gb)"}
                }
                @for d in result {
                    tr {
                        td {(d.name)}
                        td {(d.min_size)}
                    }
                }
            }
        },
    ))
}

pub async fn test_handler(Tailauth(user, _): Tailauth) -> Result<Markup> {
    Ok(base(
        Some("Test Page lol".to_string()),
        user,
        html! {
            p {"I'm baby tonx narwhal ennui crucifix taiyaki yr farm-to-table lomo locavore chillwave next level. Af palo santo bicycle rights try-hard gentrify jianbing viral heirloom actually sartorial fashion axe pickled artisan selvage cred. Celiac hammock sriracha yes plz, fit migas semiotics bruh shabby chic gluten-free chambray portland pug. Vice activated charcoal cornhole messenger bag enamel pin, put a bird on it blog ascot kale chips green juice sartorial twee retro. Try-hard hashtag umami leggings tote bag chillwave."}
            h2 {"Lumbersexual polaroid"}
            p {
                "Migas trust fund sriracha pop-up occupy. Chicharrones meggings bruh green juice squid. Brunch ennui umami fit gastropub 8-bit dreamcatcher. Bespoke portland pork belly vegan direct trade shoreditch austin franzen same +1 hoodie sustainable pickled celiac succulents. Lo-fi squid pok pok, chillwave master cleanse DIY tbh enamel pin gastropub iPhone yes plz lyft actually lumbersexual."
            }
            p {
                "Next level gastropub intelligentsia flannel tote bag, pug tilde lumbersexual poke mustache occupy. Seitan viral poutine messenger bag, echo park wayfarers af bruh poke distillery jianbing. Chillwave activated charcoal +1, disrupt shoreditch swag humblebrag lyft bushwick readymade same taxidermy kickstarter cold-pressed unicorn. Organic cloud bread polaroid tacos listicle man braid poutine chia skateboard fixie."
            }
        },
    ))
}
