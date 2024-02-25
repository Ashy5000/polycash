use actix_web::web;
use std::net::IpAddr;
use std::str::FromStr;
use std::process::Command;

pub fn verify_peer(ip: &web::Path<String>) -> bool {
    // First, check if the IP address is valid
    // This also serves as a sanitizer which is important as we are going to use the IP address in a command
    let addr_result = IpAddr::from_str(&ip.to_string());
    if addr_result.is_err() {
        return false;
    }
    // Then, check if the IP address is reachable
    let ping = Command::new("ping")
        .arg("-c")
        .arg("1")
        .arg(ip.to_string())
        .output()
        .expect("failed to execute process");
    if !ping.status.success() {
        return false;
    }
    // Next, check if the IP address is already in the list of peers
    let peers = std::fs::read_to_string("peers.txt").unwrap();
    for peer in peers.lines() {
        if peer == ip.to_string() {
            return false;
        }
    }
    return true;
}
