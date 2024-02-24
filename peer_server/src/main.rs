use actix_web::{get, web, App, HttpServer, Responder};

mod add_peer;

#[get("/add_peer/{ip}")]
async fn collect_peer_list(ip: web::Path<String>) -> impl Responder {
    add_peer::add_peer(ip)
}

#[get("/get_peers")]
async fn serve_peer_list() -> impl Responder {
    // Read the list of peers from peers.txt
    let peers = std::fs::read_to_string("peers.txt").unwrap();
    peers
}

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    HttpServer::new(|| {
        App::new()
            .service(collect_peer_list)
            .service(serve_peer_list)
    })
    .bind(("127.0.0.1", 8080))?
    .run()
    .await
}
