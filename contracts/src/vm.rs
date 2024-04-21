use std::collections::HashMap;

use crate::{buffer::Buffer, syntax_tree::SyntaxTree};

pub fn vm_access_buffer(
    buffers: &mut HashMap<String, Buffer>,
    loc: String,
    err_loc: String,
) -> Vec<u8> {
    if !vm_check_buffer_initialization(buffers, loc.to_owned()) {
        vm_throw_local_error(buffers, err_loc);
        return Vec::new();
    }
    if let Some(x) = buffers.get(&(loc.to_owned())) {
        return x.contents.clone();
    }
    Vec::new()
}

pub fn vm_check_buffer_initialization(buffers: &mut HashMap<String, Buffer>, loc: String) -> bool {
    buffers.contains_key(&(loc.clone()))
}

pub fn vm_throw_global_error(buffers: &mut HashMap<String, Buffer>) {
    if let Some(x) = buffers.get_mut(&("00000000".to_owned())) {
        *x = Buffer { contents: vec![1] };
    }
}

pub fn vm_throw_local_error(buffers: &mut HashMap<String, Buffer>, loc: String) {
    if !vm_check_buffer_initialization(buffers, loc.clone()) {
        vm_throw_global_error(buffers);
        return;
    }
    if let Some(x) = buffers.get_mut(&(loc.clone())) {
        *x = Buffer { contents: vec![1] };
    }
}

pub fn run_vm(syntax_tree: SyntaxTree, buffers: &mut HashMap<String, Buffer>) -> i64 {
    let mut line_number = 0;
    while line_number < syntax_tree.lines.len() {
        let line = syntax_tree.lines[line_number].clone();
        let should_increment = true;
        match line.command.as_str() {
            "Exit" => return line.args[0].parse::<i64>().unwrap(),
            "InitBfr" => {
                if vm_check_buffer_initialization(buffers, line.args[0].clone()) {
                    vm_throw_local_error(buffers, line.args[1].clone())
                }
                buffers.insert(
                    line.args[0].clone(),
                    Buffer {
                        contents: Vec::new(),
                    },
                );
            }
            "CpyBfr" => {
                if !vm_check_buffer_initialization(buffers, line.args[0].clone())
                    || !vm_check_buffer_initialization(buffers, line.args[1].clone())
                {
                    vm_throw_local_error(buffers, line.args[2].clone());
                }
                let src_contents: Vec<u8> =
                    vm_access_buffer(buffers, line.args[0].clone(), line.args[2].clone());
                if let Some(dst) = buffers.get_mut(&(line.args[1].clone())) {
                    dst.contents = src_contents;
                }
            }
            "FreeBfr" => {
                if !vm_check_buffer_initialization(buffers, line.args[0].clone()) {
                    vm_throw_local_error(buffers, line.args[1].clone())
                }
                buffers
                    .remove(&line.args[0].clone())
                    .expect("Failed to free memory");
            }
            "BfrStat" => {
                let status = vm_check_buffer_initialization(buffers, line.args[0].clone());
                if let Some(x) = buffers.get_mut(&(line.args[1].clone())) {
                    if status {
                        x.contents = vec![1];
                    } else {
                        x.contents = vec![0];
                    }
                }
            }
            "Add" => {
                let status_0 = vm_check_buffer_initialization(buffers, line.args[0].clone());
                let status_1 = vm_check_buffer_initialization(buffers, line.args[1].clone());
                let status_2 = vm_check_buffer_initialization(buffers, line.args[2].clone());
                if !status_0 || !status_1 || !status_2 {
                    vm_throw_local_error(buffers, line.args[3].clone());
                }
                let buffer_0 = buffers.get(&(line.args[0].clone())).unwrap();
                let buffer_1 = buffers.get(&(line.args[1].clone())).unwrap();
                let buffer_0_u64 = buffer_0.as_u64().unwrap();
                let buffer_1_u64 = buffer_1.as_u64().unwrap();
                let result_u64 = buffer_0_u64 + buffer_1_u64;
                let buffer_result = buffers.get_mut(&(line.args[2].clone())).unwrap();
                buffer_result.load_u64(result_u64);
            }
            "Stdout" => {
                if !vm_check_buffer_initialization(buffers, line.args[0].clone()) {
                    vm_throw_local_error(buffers, line.args[1].clone())
                }
                println!("{:?}", buffers.get(&line.args[0]).unwrap().contents)
            }
            "SetCnst" => {
                if !vm_check_buffer_initialization(buffers, line.args[0].clone()) {
                    vm_throw_local_error(buffers, line.args[2].clone())
                }
                if let Some(x) = buffers.get_mut(&(line.args[0].clone())) {
                    x.contents =
                        hex::decode(line.args[1].clone()).expect("Failed to parse raw hex value");
                }
            }
            &_ => vm_throw_global_error(buffers),
        }
        if should_increment {
            line_number += 1;
        }
    }
    return 0;
}
