use clap::Parser;
use names::{Generator, Name};

#[derive(Parser)]
#[clap(author, version, about, long_about = None)]
struct Cli {
    #[clap(short, long, default_value = "5")]
    count: usize,

    #[clap(short, long)]
    add_numbers: bool,
}

fn main() {
    let cli = Cli::parse();

    let generator = names::Generator::new(
        &rotbart::COMBINED_ADJ,
        &rotbart::COMBINED_NOUN,
        if cli.add_numbers {
            Name::Numbered
        } else {
            Name::Plain
        },
    );
    generator
        .take(cli.count)
        .for_each(|name| println!("{name}"));
}
