use std::collections::HashMap;

use crate::{
    buffer::Buffer,
    math::{execute_math_operation, Add, And, Divide, Eq, Multiply, Not, Or, Subtract},
    syntax_tree::SyntaxTree,
};

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
    let mut should_increment;
    while line_number < syntax_tree.lines.len() {
        let line = syntax_tree.lines[line_number].clone();
        should_increment = true;
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
            "Add" => execute_math_operation(
                Add {},
                buffers,
                line.args[0].clone(),
                line.args[1].clone(),
                line.args[2].clone(),
                line.args[3].clone(),
            ),
            "Sub" => execute_math_operation(
                Subtract {},
                buffers,
                line.args[0].clone(),
                line.args[1].clone(),
                line.args[2].clone(),
                line.args[3].clone(),
            ),
            "Mul" => execute_math_operation(
                Multiply {},
                buffers,
                line.args[0].clone(),
                line.args[1].clone(),
                line.args[2].clone(),
                line.args[3].clone(),
            ),
            "Div" => execute_math_operation(
                Divide {},
                buffers,
                line.args[0].clone(),
                line.args[1].clone(),
                line.args[2].clone(),
                line.args[3].clone(),
            ),
            "Eq" => execute_math_operation(
                Eq {},
                buffers,
                line.args[0].clone(),
                line.args[1].clone(),
                line.args[2].clone(),
                line.args[3].clone(),
            ),
            "And" => execute_math_operation(
                And {},
                buffers,
                line.args[0].clone(),
                line.args[1].clone(),
                line.args[2].clone(),
                line.args[3].clone(),
            ),
            "Or" => execute_math_operation(
                Or {},
                buffers,
                line.args[0].clone(),
                line.args[1].clone(),
                line.args[2].clone(),
                line.args[3].clone(),
            ),
            "Not" => execute_math_operation(
                Not {},
                buffers,
                line.args[0].clone(),
                "".to_owned(),
                line.args[1].clone(),
                line.args[2].clone(),
            ),
            "Jmp" => {
                line_number = line.args[0].parse::<usize>().unwrap() - 1;
                should_increment = false;
            }
            "JmpCond" => {
                if !vm_check_buffer_initialization(buffers, line.args[0].clone()) {
                    vm_throw_local_error(buffers, line.args[2].clone())
                }
                if buffers.get(&line.args[0]).unwrap().as_u64() != Ok(0) {
                    line_number = line.args[1].parse::<usize>().unwrap() - 1;
                    should_increment = false;
                }
            }
            "Stdout" => {
                if !vm_check_buffer_initialization(buffers, line.args[0].clone()) {
                    vm_throw_local_error(buffers, line.args[1].clone())
                }
                println!("{:?}", buffers.get(&line.args[0]).unwrap().contents)
            }
            "Stderr" => {
                if !vm_check_buffer_initialization(buffers, line.args[0].clone()) {
                    vm_throw_local_error(buffers, line.args[1].clone())
                }
                eprintln!("{:?}", buffers.get(&line.args[0]).unwrap().contents)
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
            "Tx" => {
                let sender_bytes =
                    vm_access_buffer(buffers, line.args[0].clone(), line.args[3].clone());
                let sender = match String::from_utf8(sender_bytes) {
                    Ok(v) => v,
                    Err(_) => {
                        vm_throw_local_error(buffers, line.args[3].clone());
                        "".to_owned()
                    }
                };
                let receiver_bytes =
                    vm_access_buffer(buffers, line.args[1].clone(), line.args[3].clone());
                let receiver = match String::from_utf8(receiver_bytes) {
                    Ok(v) => v,
                    Err(_) => {
                        vm_throw_local_error(buffers, line.args[3].clone());
                        "".to_owned()
                    }
                };
                let amount_bytes =
                    vm_access_buffer(buffers, line.args[2].clone(), line.args[3].clone());
                let amount = match String::from_utf8(amount_bytes) {
                    Ok(v) => v,
                    Err(_) => {
                        vm_throw_local_error(buffers, line.args[3].clone());
                        "".to_owned()
                    }
                };
                println!("TX {} {} {}", sender, receiver, amount);
            }
            &_ => vm_throw_global_error(buffers),
        }
        if should_increment {
            line_number += 1;
        }
    }
    return 0;
}
