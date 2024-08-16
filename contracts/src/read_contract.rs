// Copyright 2024, Asher Wrobel
/*
This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.
*/
use std::{env, fs};
use smartstring::alias::String;

pub fn read_contract() -> String {
    let args: Vec<std::string::String> = env::args().collect();
    if args.len() != 4 {
        panic!("Expected exactly 4 arguments, got {}", args.len())
    }
    let contract_path = &args[1];
    fs::read_to_string(contract_path)
        .expect("Should have been able to read smart contract file")
        .parse()
        .unwrap()
}
