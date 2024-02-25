use actix_web::{web, Responder};
use crate::verify_peer::verify_peer;

pub fn add_peer(ip: web::Path<String>) -> impl Responder {
    // Validate the IP address
    if !verify_peer(&ip) {
        return format!("Invalid peer: {}", ip);
    }
    // Add the peer to the list of peers in peers.txt
    std::fs::write("peers.txt", ip.to_string() + "\n").unwrap();
    println!("Peer added: {}", ip);
    format!("Peer added: {}", ip)
}
