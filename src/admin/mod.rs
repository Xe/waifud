use maud::{html, Markup, PreEscaped};
use ts_localapi::User;

use crate::{tailauth::Tailauth, Result};

const CSS: PreEscaped<&'static str> = PreEscaped(include_str!("./xess.css"));

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
            }
            body.top {
                main {
                    nav.nav {
                        a.left href="/admin/" {"Home"}
                        div.right {
                            {(user_data.display_name)}
                            " "
                            img style="width:32px;height:32px" src=(user_data.profile_pic_url);
                        }
                    }
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
