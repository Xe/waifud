use scraper::Html;

use crate::models::Distro;

/// # Scraper for https://geo.mirror.pkgbuild.com/images/
///
/// This scrapes the Arch Linux cloud image site and extracts out the URL of the latest
/// release of Arch Linux and its sha256 sum.

const RELEASE_BASE: &'static str = "https://geo.mirror.pkgbuild.com/images/";

pub async fn scrape() -> crate::Result<crate::models::Distro> {
    let sel = scraper::Selector::parse("a").expect("selector to parse");

    let response_html = reqwest::get("https://geo.mirror.pkgbuild.com/images/")
        .await?
        .error_for_status()?
        .text()
        .await?;

    let doc = Html::parse_document(&response_html);

    let link = doc.select(&sel).last().ok_or(crate::Error::Catchall(
        "can't get last element of Arch image list".to_string(),
    ))?;

    let u = url::Url::parse(RELEASE_BASE)?.join(link.value().attr("href").ok_or(
        crate::Error::Catchall("link has no href, how???".to_string()),
    )?)?;

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
        .filter(|path| path.starts_with("Arch-Linux-x86_64-cloudimg-"))
        .map(ToString::to_string)
        .collect();

    if links.len() != 4 {
        return Err(crate::Error::Catchall(
            "wrong number of things in the list, wanted 4".to_string(),
        ));
    }

    let image_url = links.get(0).unwrap();
    let shasum_url = links.get(1).unwrap();

    let shasum = reqwest::get(u.join(&shasum_url)?)
        .await?
        .error_for_status()?
        .text()
        .await?;

    let shasum = shasum
        .split("  ")
        .take(1)
        .map(ToString::to_string)
        .collect::<Vec<String>>();

    let shasum = shasum.get(0).unwrap();

    let image_url = u.join(&image_url)?.as_str().to_string();

    Ok(Distro {
        name: "arch".to_string(),
        download_url: image_url,
        sha256sum: shasum.to_string(),
        min_size: 2,
        format: "waifud://qcow2".to_string(),
    })
}
