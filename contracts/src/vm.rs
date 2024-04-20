use std::collections::HashMap;

use crate::{syntax_tree::SyntaxTree, buffer::Buffer};

pub fn vm_check_buffer_initialization(buffers: &mut HashMap<String, Buffer>, loc: String) -> bool {
    buffers.contains_key(&(loc.clone()))
}

pub fn vm_throw_global_error(buffers: &mut HashMap<String, Buffer>) {
    if let Some(x) = buffers.get_mut(&("00000000".to_owned())) {
        *x = Buffer{
            contents: vec![1],
        };
    }
}

pub fn vm_throw_local_error(buffers: &mut HashMap<String, Buffer>, loc: String) {
    if !vm_check_buffer_initialization(buffers, loc.clone()) {
        vm_throw_global_error(buffers);
        return
    }
    if let Some(x) = buffers.get_mut(&(loc.clone())) {
        *x = Buffer{
            contents: vec![1],
        };
    }
}

pub fn run_vm(syntax_tree: SyntaxTree, buffers: &mut HashMap<String, Buffer>) {
    let mut line_number = 0;
    while line_number < syntax_tree.lines.len() {
        let line = syntax_tree.lines[line_number].clone();
        let should_increment = true;
        match line.command.as_str() {
            "InitBfr" => {
                if vm_check_buffer_initialization(buffers, line.args[0].clone()) {
                    vm_throw_local_error(buffers, line.args[1].clone())
                }
                buffers.insert(
                    line.args[0].clone(),
                    Buffer{
                        contents: Vec::new(),
                    }
                );
            },
            "Stdout" => {
                if !vm_check_buffer_initialization(buffers, line.args[0].clone()) {
                    vm_throw_local_error(buffers, line.args[1].clone())
                }
                println!("{:?}", buffers.get(&line.args[0]).unwrap().contents)
            },
            "SetCnst" => {
                if !vm_check_buffer_initialization(buffers, line.args[0].clone()) {
                    vm_throw_local_error(buffers, line.args[2].clone())
                }
                if let Some(x) = buffers.get_mut(&(line.args[0].clone())) {
                    x.contents = hex::decode(line.args[1].clone()).expect("Failed to parse raw hex value");
                }
            }
            "_" => vm_throw_global_error(buffers),
            &_ => vm_throw_global_error(buffers),
        }
        if should_increment {
            line_number += 1;
        }
    }
}