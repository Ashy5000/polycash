use std::process::Command;
use std::net::IpAddr;
use actix_web::{web, Responder};
use std::str::FromStr;

pub fn add_peer(ip: web::Path<String>) -> impl Responder {
    // Run checks on the IP address
    // First, check if the IP address is valid
    // This also serves as a sanitizer which is important as we are going to use the IP address in a command
    let addr_result = IpAddr::from_str(&ip.to_string());
    if addr_result.is_err() {
        return format!("Invalid IP address: {}", ip);
    }
    // Then, check if the IP address is reachable
    let ping = Command::new("ping")
        .arg("-c")
        .arg("1")
        .arg(ip.to_string())
        .output()
        .expect("failed to execute process");
    if !ping.status.success() {
        return format!("Peer not reachable: {}", ip);
    }
    // Next, check if the IP address is already in the list of peers
    let peers = std::fs::read_to_string("peers.txt").unwrap();
    for peer in peers.lines() {
        if peer == ip.to_string() {
            return format!("Peer already exists: {}", ip);
        }
    }
    // Add the peer to the list of peers in peers.txt
    std::fs::write("peers.txt", ip.to_string() + "\n").unwrap();
    println!("Peer added: {}", ip);
    format!("Peer added: {}", ip)
}
