mod blaseball;
mod elfs;
mod mlp_fim;
mod pokemon;
mod xc1;
mod xc2;

lazy_static::lazy_static! {
    static ref COMBINED_ADJ: Vec<&'static str> = {
        let mut adjs: Vec<&str> = vec![];
        adjs.extend(xc1::ADJECTIVES.iter());
        adjs.extend(xc2::ADJECTIVES.iter());
        adjs.extend(elfs::ADJECTIVES.iter());
        adjs.extend(blaseball::FIRST_NAMES.iter());
        adjs.sort();
        adjs.dedup();
        adjs
    };

    static ref COMBINED_NOUN: Vec<&'static str> = {
        let mut nouns: Vec<&str> = vec![];
        nouns.extend(xc1::NOUNS.iter());
        nouns.extend(xc2::NOUNS.iter());
        nouns.extend(xc2::COMMON_BLADES.iter());
        nouns.extend(elfs::NOUNS.iter());
        nouns.extend(mlp_fim::PONIES.iter());
        nouns.extend(pokemon::POKEDEX.iter());
        nouns.extend(blaseball::LAST_NAMES.iter());
        nouns.sort();
        nouns.dedup();
        nouns
    };
}

pub fn unique_monster() -> Option<String> {
    let mut generator = names::Generator::new(&COMBINED_ADJ, &COMBINED_NOUN, names::Name::Plain);
    generator.next()
}
