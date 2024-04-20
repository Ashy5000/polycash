mod read_contract;
mod syntax_tree;
mod buffer;
mod vm;

use crate::{buffer::Buffer, vm::run_vm};

use std::collections::HashMap;

fn main() {
    let contract_contents = read_contract::read_contract();
    let mut tree = syntax_tree::build_syntax_tree();
    tree.create(contract_contents);
    let mut buffers: HashMap<String, Buffer> = HashMap::new();
    buffers.insert(
        "00000000".to_owned(),
        Buffer {
            contents: Vec::new(),
        }
    );
    run_vm(tree, &mut buffers);
}
