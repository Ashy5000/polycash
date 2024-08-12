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

use std::env;
use std::process::ExitCode;
use crate::{buffer::Buffer, vm::run_vm};
use smartstring::alias::String;

use rustc_hash::FxHashMap;

#[global_allocator]
static GLOBAL: tikv_jemallocator::Jemalloc = tikv_jemallocator::Jemalloc;

fn main() -> ExitCode {
    let contract_contents = read_contract::read_contract();
    let mut tree = syntax_tree::build_syntax_tree();
    tree.create(contract_contents);
    let mut buffers: FxHashMap<String, Buffer> = FxHashMap::default();
    buffers.insert(
        "00000000".into(),
        Buffer {
            contents: Vec::new(),
        },
    );
    let blockutil_interface = blockutil::BlockUtilInterface::default();
    let args: Vec<std::string::String> = env::args().collect();
    let contract_hash = &args[2];
    let gas_limit: f64 = args[3].parse().unwrap();
    let (exit_code, gas_used) = run_vm(tree, &mut buffers, blockutil_interface, contract_hash.parse().unwrap(), gas_limit);
    println!("Gas used: {}", gas_used);
    ExitCode::from(exit_code as u8)
}
