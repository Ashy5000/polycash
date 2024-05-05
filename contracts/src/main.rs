mod blockutil;
mod buffer;
mod math;
mod read_contract;
mod syntax_tree;
mod vm;

use crate::{buffer::Buffer, vm::run_vm};

use std::{collections::HashMap, process::ExitCode};

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
    let exit_code = run_vm(tree, &mut buffers, blockutil_interface);
    ExitCode::from(exit_code as u8)
}
