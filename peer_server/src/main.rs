// Copyright 2024, Asher Wrobel
/*
This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.
*/
use actix_web::{get, web, App, HttpServer, Responder};

mod add_peer;
mod verify_peer;

#[get("/add_peer/{ip}")]
async fn collect_peer_list(ip: web::Path<String>) -> impl Responder {
    add_peer::add_peer("http://".to_string() + &ip)
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
    .bind(("0.0.0.0", 6060))? // Note the default peer server runs on port 8080 due to an outdated version of the Github repository. This is the only difference in the source code.
    .run()
    .await
}
