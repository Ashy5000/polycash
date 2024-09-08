// Copyright 2024, Asher Wrobel
/*
This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.
*/
mod blockutil;
mod buffer;
mod math;
mod read_contract;
mod sanitization;
mod syntax_tree;
mod vm;
mod stack;
mod state;

mod msgpack;

use std::env;
use std::fs::File;
use std::process::ExitCode;
use hex::decode;
use crate::msgpack::decode_pending_state;
use crate::vm::run_vm;

#[global_allocator]
static GLOBAL: tikv_jemallocator::Jemalloc = tikv_jemallocator::Jemalloc;

fn main() -> ExitCode {
    let contract_contents = read_contract::read_contract();
    let args: Vec<std::string::String> = env::args().collect();
    let contract_hash = &args[2];
    let gas_limit: f64 = args[3].parse().unwrap();
    let sender: Vec<u8> = args[4].clone().into();
    let pending_state = decode_pending_state();
    let (exit_code, gas_used) = run_vm(contract_contents, contract_hash.parse().unwrap(), gas_limit, sender, pending_state);
    println!("Gas used: {}", gas_used);
    ExitCode::from(exit_code as u8)
}
