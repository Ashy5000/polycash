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

use crate::{buffer::Buffer, vm::run_vm};

use std::{collections::HashMap, env, process::ExitCode};

fn main() -> ExitCode {
    let contract_contents = read_contract::read_contract();
    let mut tree = syntax_tree::build_syntax_tree();
    tree.create(contract_contents);
    let mut buffers: HashMap<String, Buffer> = HashMap::new();
    buffers.insert(
        "00000000".to_owned(),
        Buffer {
            contents: Vec::new(),
        },
    );
    let blockutil_interface = blockutil::BlockUtilInterface::new();
    let args: Vec<String> = env::args().collect();
    let contract_hash = &args[2];
    let (exit_code, gas_used) = run_vm(tree, &mut buffers, blockutil_interface, contract_hash.clone());
    println!("Gas used: {}", gas_used);
    ExitCode::from(exit_code as u8)
}
