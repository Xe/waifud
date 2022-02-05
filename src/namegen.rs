use rand::prelude::*;

pub fn next() -> String {
    let data: &[u8] = include_bytes!("../data/names.json");
    let names: Vec<String> = serde_json::from_slice(data).unwrap();

    let i = thread_rng().gen::<usize>() % names.len();
    names[i].clone()
}
