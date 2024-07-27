// Copyright 2024, Asher Wrobel
/*
This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.
*/
use regex::Regex;
use std::process::Command;
use std::str::FromStr;
use crate::key::get_public_key;

pub(crate) fn sync() -> f64 {
    let binding = get_public_key().expect("Failed to get public key");
    let public_key = binding.as_str();
    if public_key == "" {
        return -1.0
    }
    let output = Command::new("env")
        .arg("-C")
        .arg("..")
        .arg("./builds/node/node")
        .arg("-command")
        .arg("sync;balance ".to_owned() + public_key)
        .output()
        .expect("Failed to run node executable");
    let output_string =
        String::from_utf8(output.stdout).expect("Failed to convert output to string");
    let re = Regex::new(r": (?<balance>[0-9]\.[0-9]*)").unwrap();
    let Some(caps) = re.captures(&output_string) else {
        println!("no match!");
        return -1.0;
    };
    let balance = f64::from_str(&caps["balance"]).expect("Invalid balance format");
    balance
}
