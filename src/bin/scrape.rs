#[tokio::main]
async fn main() -> waifud::Result<()> {
    tracing_subscriber::fmt::init();

    let arch_info = waifud::scrape::get_all().await?;

    println!("{}", serde_dhall::serialize(&arch_info).to_string()?);

    Ok(())
}
