use rocket::{fs::FileServer, get, launch, routes};
use waifud::MainDatabase;

#[get("/")]
fn hello() -> &'static str {
    "Hello, world!"
}

#[launch]
fn rocket() -> _ {
    rocket::build()
        .attach(MainDatabase::fairing())
        .mount("/", routes![hello])
        .mount("/static", FileServer::from("public"))
}
