// Copyright 2024, Asher Wrobel
/*
This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.
*/
use actix_web::web;
use std::net::IpAddr;
use std::str::FromStr;
use std::process::Command;

pub fn verify_peer(ip: &String) -> bool {
    // First, check if the IP address is valid
    // This also serves as a sanitizer which is important as we are going to use the IP address in a command
    let addr_result = IpAddr::from_str(&ip);
    if addr_result.is_err() {
        return false;
    }
    // Then, check if the IP address is reachable
    let ping = Command::new("ping")
        .arg("-c")
        .arg("1")
        .arg(ip)
        .output()
        .expect("failed to execute process");
    if !ping.status.success() {
        return false;
    }
    // Next, check if the IP address is already in the list of peers
    let peers = std::fs::read_to_string("peers.txt").unwrap();
    for peer in peers.lines() {
        if peer == ip {
            return false;
        }
    }
    return true;
}
