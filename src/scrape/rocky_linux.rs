use scraper::{ElementRef, Html};

use crate::models::Distro;

/// # Scraper for Rocky Linux cloud images
///
/// This scrapes the Rocky Linux cloud image site and extracts out the URL of the latest
/// release of Rocky Linux and its sha256 sum.

const RELEASE_BASE: &'static str = "http://download.rockylinux.org/pub/rocky/";

#[instrument]
pub async fn scrape(version: i32) -> crate::Result<crate::models::Distro> {
    let sel = scraper::Selector::parse("a").expect("selector to parse");

    let mut base: String = RELEASE_BASE.to_string();
    base.push_str(&format!("{}", version));
    base.push_str("/images/");

    if version == 9 {
        base.push_str("x86_64/");
    }

    let u = url::Url::parse(&base)?;
    debug!("url: {u}");

    let response_html = reqwest::get(u.as_str())
        .await?
        .error_for_status()?
        .text()
        .await?;

    let doc = Html::parse_document(&response_html);

    let elems = doc.select(&sel).collect::<Vec<ElementRef>>();
    let elems = elems
        .into_iter()
        .rev()
        .filter(|elem| elem.value().attr("href").is_some())
        .map(|elem| elem.value().attr("href").unwrap())
        .filter(|link| link.contains("GenericCloud"))
        .filter(|link| link.contains("x86_64"))
        .filter(|link| !link.contains("latest"))
        .filter(|link| link.ends_with(".qcow2"))
        .collect::<Vec<&str>>();
    let link = elems.get(0).ok_or(crate::Error::Catchall(
        "can't get second to last element of image list".to_string(),
    ))?;

    let u = u.join(link)?;
    debug!(url = u.to_string(), link = link);

    let image_url = u.to_string();
    let shasum_url = u.join("./CHECKSUM")?;

    debug!("shasum url: {shasum_url}");
    let mut shasums = reqwest::get(shasum_url)
        .await?
        .error_for_status()?
        .text()
        .await?;
    if version != 8 {
        shasums = shasums
            .split("\n")
            .filter(|line| line.contains("SHA256"))
            .collect::<Vec<&str>>()
            .join("\n");
    }

    let mut shasum = String::new();

    for line in shasums.split("\n").filter(|line| line.contains(link)) {
        if line == "" {
            break;
        }
        if version != 8 {
            let sides: Vec<&str> = line.split(" ").collect();
            if sides.len() != 4 {
                error!("Somehow this doesn't have 3 spaces in it {line:?}");
                continue;
            }
            shasum = sides.get(3).unwrap().to_string()
        } else {
            let sides: Vec<&str> = line.split("  ").collect();
            if sides.len() != 2 {
                error!("Somehow this doesn't have two spaces in it {line:?}");
                continue;
            }
            shasum = sides.get(0).unwrap().to_string();
        }
    }

    let image_url = u.join(&image_url)?.as_str().to_string();

    Ok(Distro {
        name: format!("rocky-linux-{version}"),
        download_url: image_url,
        sha256sum: shasum.to_string(),
        min_size: 10,
        format: "waifud://qcow2".to_string(),
    })
}
