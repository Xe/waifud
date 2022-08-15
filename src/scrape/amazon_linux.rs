use crate::{models::Distro, Error};
use scraper::Html;
use url::Url;

/// # Scraper for Amazon Linux
///
/// This scrapes the Amazon Linux cloud image site and extracts out the URL of the latest
/// release of Amazon Linux and its sha256 sum.

pub async fn scrape() -> crate::Result<crate::models::Distro> {
    let sel = scraper::Selector::parse("a").expect("selector to parse");

    let res = {
        let https = hyper_tls::HttpsConnector::new();
        let cli = hyper::Client::builder().build::<_, hyper::Body>(https);
        cli.get(hyper::Uri::from_static(
            "https://cdn.amazonlinux.com/os-images/latest/kvm/",
        ))
        .await
    }?;

    let release_base = res
        .headers()
        .get(axum::http::header::LOCATION)
        .ok_or(Error::Catchall(
            "why did the redirect not work?".to_string(),
        ))?
        .to_str()?;

    let response_html = reqwest::get(release_base)
        .await?
        .error_for_status()?
        .text()
        .await?;

    let doc = Html::parse_document(&response_html);

    let link = doc
        .select(&sel)
        .filter(|elem| elem.value().attr("href").is_some())
        .map(|elem| elem.value().attr("href").unwrap())
        .filter(|link| link.starts_with("amzn2-kvm"))
        .take(1)
        .map(ToString::to_string)
        .collect::<Vec<String>>();
    let link = link.get(0).ok_or(crate::Error::Catchall(
        "can't get last element of Amazon Linux image list".to_string(),
    ))?;

    let image_url = Url::parse(release_base)?.join(link)?;
    let shasum_url = Url::parse(release_base)?.join("SHA256SUMS")?;

    let shasum = reqwest::get(shasum_url)
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

    Ok(Distro {
        name: "amazon-linux-2".to_string(),
        download_url: image_url.to_string(),
        sha256sum: shasum.to_string(),
        min_size: 25,
        format: "waifud://qcow2".to_string(),
    })
}
