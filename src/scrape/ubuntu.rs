use scraper::{ElementRef, Html};
use std::collections::HashMap;

use crate::{models::Distro, Error};

/// # Scraper for Ubuntu cloud images
///
/// This scrapes the Ubuntu cloud image site and extracts out the URL of the latest
/// release of Ubuntu and its sha256 sum.

const RELEASE_BASE: &'static str = "http://cloud-images.ubuntu.com/daily/server/";

pub async fn scrape((version, name): (&str, &str)) -> crate::Result<crate::models::Distro> {
    let sel = scraper::Selector::parse("a").expect("selector to parse");

    let mut base: String = RELEASE_BASE.to_string();
    base.push_str(&name);
    base.push_str("/");

    let u = url::Url::parse(&base)?;
    debug!("url: {u}");

    let response_html = reqwest::get(u.as_str())
        .await?
        .error_for_status()?
        .text()
        .await?;

    let doc = Html::parse_document(&response_html);

    let elems = doc.select(&sel).collect::<Vec<ElementRef>>();
    let elems = elems.into_iter().rev().collect::<Vec<ElementRef>>();
    let link = elems.get(2).ok_or(crate::Error::Catchall(
        "can't get second to last element of image list".to_string(),
    ))?;

    let u = u
        .join(name)?
        .join(link.value().attr("href").ok_or(crate::Error::Catchall(
            "link has no href, how???".to_string(),
        ))?)?;
    debug!("url: {u}");

    let response_html = reqwest::get(u.as_str())
        .await?
        .error_for_status()?
        .text()
        .await?;

    let doc = Html::parse_document(&response_html);

    let links: Vec<String> = doc
        .select(&sel)
        .filter(|elem| elem.value().attr("href").is_some())
        .map(|elem| elem.value().attr("href"))
        .map(Option::unwrap)
        .filter(|path| path.ends_with("-server-cloudimg-amd64.img"))
        .map(ToString::to_string)
        .collect();

    let image_url = links.get(0).unwrap();
    let shasum_url = u.join("SHA256SUMS")?;

    let shasums = reqwest::get(shasum_url)
        .await?
        .error_for_status()?
        .text()
        .await?;

    let mut sha_map = HashMap::<String, String>::new();

    for line in shasums.split("\n") {
        let sides: Vec<&str> = line.split(" *").collect();
        if sides.len() != 2 {
            error!("Somehow this doesn't have two spaces in it {line:?}");
            continue;
        }

        sha_map.insert(
            sides.get(1).unwrap().to_string(),
            sides.get(0).unwrap().to_string(),
        );
    }

    // println!("{}", serde_dhall::serialize(&sha_map).to_string()?);

    let mut key: String = name.clone().to_string();
    key.push_str("-server-cloudimg-amd64.img");
    let shasum = sha_map
        .get(&key)
        .ok_or(Error::Catchall(format!("can't find shasum for {name}")))?;

    let image_url = u.join(&image_url)?.as_str().to_string();

    Ok(Distro {
        name: format!("ubuntu-{version}"),
        download_url: image_url,
        sha256sum: shasum.to_string(),
        min_size: 5,
        format: "waifud://qcow2".to_string(),
    })
}
