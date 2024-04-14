// Copyright 2024, Asher Wrobel
/*
This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.
*/
use std::process::Command;

#[derive(Debug)]
pub(crate) struct PublicKeyError {}

pub(crate) fn get_public_key() -> Result<String, PublicKeyError> {
    let output = Command::new("env")
        .arg("-C")
        .arg("..")
        .arg("./builds/node/node_linux-amd64")
        .arg("-command")
        .arg("getpubkey")
        .output()
        .expect("Failed to run node executable");
    let s = match String::from_utf8(output.stdout) {
        Ok(v) => v,
        Err(_) => return Err(PublicKeyError{}),
};
    return Ok(s)
}

pub(crate) fn generate_key() {
    let status = Command::new("env")
        .arg("-C")
        .arg("..")
        .arg("./builds/node/node_linux-amd64")
        .arg("-command")
        .arg("keygen")
        .status()
        .expect("Failed to run node executable");
    if status.success() {
        println!("Key generated successfully!");
    } else {
        println!("Failed to generate key");
    }
}
