// Copyright 2024, Asher Wrobel
/*
This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.
*/
use rustc_hash::FxHashMap;
use std::ffi::c_void;
use crate::{
    blockutil::BlockUtilInterface,
    buffer::Buffer,
    math::{execute_math_operation, Add, And, Divide, Eq, Less, Multiply, Not, Or, Subtract, Modulo, Exp},
    syntax_tree::SyntaxTree,
};
use crate::stack::Stack;
use smartstring::alias::String;

const VM_NIL: *const c_void = 0x0000 as *const c_void;

#[inline(always)]
pub fn vm_access_buffer_contents(
    buffers: &mut FxHashMap<String, Buffer>,
    loc: String,
    err_loc: String,
) -> Vec<u8> {
    if !vm_check_buffer_initialization(buffers, loc.to_owned()) {
        vm_throw_local_error(buffers, err_loc);
        return vec![];
    }
    if let Some(x) = buffers.get(&(loc.to_owned())) {
        return x.contents.clone();
    }
    vec![]
}

#[inline(always)]
pub fn vm_access_buffer(
    buffers: &mut FxHashMap<String, Buffer>,
    loc: String,
    err_loc: String,
) -> *mut Buffer {
    if !vm_check_buffer_initialization(buffers, loc.to_owned()) {
        vm_throw_local_error(buffers, err_loc);
        return VM_NIL as *mut Buffer;
    }
    if let Some(x) = buffers.get_mut(&(loc.to_owned())) {
        return x;
    }
    VM_NIL as *mut Buffer
}

#[inline(always)]
pub fn vm_check_buffer_initialization(buffers: &mut FxHashMap<String, Buffer>, loc: String) -> bool {
    buffers.contains_key(&(loc.clone()))
}

#[inline(always)]
pub fn vm_throw_global_error(buffers: &mut FxHashMap<String, Buffer>) {
    if let Some(x) = buffers.get_mut::<String>(&("00000000".into())) {
        *x = Buffer { contents: vec![1] };
    }
}

#[inline(always)]
pub fn vm_throw_local_error(buffers: &mut FxHashMap<String, Buffer>, loc: String) {
    if !vm_check_buffer_initialization(buffers, loc.clone()) {
        vm_throw_global_error(buffers);
        return;
    }
    if let Some(x) = buffers.get_mut(&(loc.clone())) {
        *x = Buffer { contents: vec![1] };
    }
}

pub fn run_vm(
    syntax_tree: SyntaxTree,
    buffers: &mut FxHashMap<String, Buffer>,
    blockutil_interface: BlockUtilInterface,
    contract_hash: String
) -> (i64, f64) {
    let mut line_number = 0;
    let mut should_increment;
    let mut gas_used = 0.0;
    let mut origins = vec![];
    let mut stack = Stack{frames: vec![]};
    while line_number < syntax_tree.lines.len() {
        let line = &syntax_tree.lines[line_number];
        should_increment = true;
        match line.command.as_str() {
            "Exit" => {
                gas_used += 1.0;
                return (line.args[0].parse::<i64>().unwrap(), gas_used);
            }
            "ExitBfr" => {
                gas_used += 1.0;
                if !vm_check_buffer_initialization(buffers, line.args[0].clone()) {
                    vm_throw_local_error(buffers, line.args[1].clone())
                }
                let exit_code = buffers
                    .get(&line.args[0].clone())
                    .expect("Unable to read exit code buffer")
                    .as_u64()
                    .expect("Could not convert exit code to u64");
                return (exit_code as i64, gas_used);
            }
            "InitBfr" => {
                buffers.insert(
                    line.args[0].clone(),
                    Buffer {
                        contents: Vec::new(),
                    },
                );
                gas_used += 2.0;
                gas_used += line.args[0].len() as f64 / 10.0;
            }
            "CpyBfr" => {
                if !vm_check_buffer_initialization(buffers, line.args[0].clone())
                {
                    vm_throw_local_error(buffers, line.args[2].clone());
                }
                if !vm_check_buffer_initialization(buffers, line.args[1].clone()) {
                    buffers.insert(
                        line.args[1].clone(),
                        Buffer {
                            contents: Vec::new(),
                        },
                    );
                    gas_used += 2.0;
                }
                let src_contents: Vec<u8> =
                    vm_access_buffer_contents(buffers, line.args[0].clone(), line.args[2].clone());
                if let Some(dst) = buffers.get_mut(&(line.args[1].clone())) {
                    dst.contents = src_contents.clone();
                    gas_used += src_contents.len() as f64 / 10.0;
                }
                gas_used += 2.0;
            }
            "FreeBfr" => {
                if !vm_check_buffer_initialization(buffers, line.args[0].clone()) {
                    vm_throw_local_error(buffers, line.args[1].clone())
                }
                buffers
                    .remove(&line.args[0].clone())
                    .expect("Failed to free memory");
                gas_used += 2.0;
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
                gas_used += 1.0;
            }
            "BfrLen" => {
                if !vm_check_buffer_initialization(buffers, line.args[0].clone())
                    || !vm_check_buffer_initialization(buffers, line.args[1].clone())
                {
                    vm_throw_local_error(buffers, line.args[2].clone())
                }
                let x = vm_access_buffer_contents(buffers, line.args[0].clone(), line.args[2].clone());
                if let Some(y) = buffers.get_mut(&(line.args[1].clone())) {
                    y.load_u64(x.len().try_into().unwrap());
                }
            }
            "Add" => {
                execute_math_operation(
                    Add {},
                    buffers,
                    line.args[0].clone(),
                    line.args[1].clone(),
                    line.args[2].clone(),
                    line.args[3].clone(),
                );
                gas_used += 0.5;
            }
            "Sub" => {
                execute_math_operation(
                    Subtract {},
                    buffers,
                    line.args[0].clone(),
                    line.args[1].clone(),
                    line.args[2].clone(),
                    line.args[3].clone(),
                );
                gas_used += 0.5;
            }
            "Mul" => {
                execute_math_operation(
                    Multiply {},
                    buffers,
                    line.args[0].clone(),
                    line.args[1].clone(),
                    line.args[2].clone(),
                    line.args[3].clone(),
                );
                gas_used += 1.0;
            }
            "Div" => {
                execute_math_operation(
                    Divide {},
                    buffers,
                    line.args[0].clone(),
                    line.args[1].clone(),
                    line.args[2].clone(),
                    line.args[3].clone(),
                );
                gas_used += 1.0;
            }
            "Exp" => {
                execute_math_operation(
                    Exp {},
                    buffers,
                    line.args[0].clone(),
                    line.args[1].clone(),
                    line.args[2].clone(),
                    line.args[3].clone(),
                );
                gas_used += 1.5;
            }
            "Mod" => execute_math_operation(
                Modulo {},
                buffers,
                line.args[0].clone(),
                line.args[1].clone(),
                line.args[2].clone(),
                line.args[3].clone(),
            ),
            "Eq" => {
                execute_math_operation(
                    Eq {},
                    buffers,
                    line.args[0].clone(),
                    line.args[1].clone(),
                    line.args[2].clone(),
                    line.args[3].clone(),
                );
                gas_used += 0.5;
            }
            "Less" => {
                execute_math_operation(
                    Less {},
                    buffers,
                    line.args[0].clone(),
                    line.args[1].clone(),
                    line.args[2].clone(),
                    line.args[3].clone(),
                );
                gas_used += 0.5;
            }
            "And" => {
                execute_math_operation(
                    And {},
                    buffers,
                    line.args[0].clone(),
                    line.args[1].clone(),
                    line.args[2].clone(),
                    line.args[3].clone(),
                );
                gas_used += 1.0;
            }
            "Or" => {
                execute_math_operation(
                    Or {},
                    buffers,
                    line.args[0].clone(),
                    line.args[1].clone(),
                    line.args[2].clone(),
                    line.args[3].clone(),
                );
                gas_used += 1.0;
            }
            "Not" => {
                execute_math_operation(
                    Not {},
                    buffers,
                    line.args[0].clone(),
                    "".into(),
                    line.args[1].clone(),
                    line.args[2].clone(),
                );
                gas_used += 0.5
            }
            "App" => {
                if !vm_check_buffer_initialization(buffers, line.args[0].clone())
                    || !vm_check_buffer_initialization(buffers, line.args[1].clone())
                {
                    vm_throw_local_error(buffers, line.args[1].clone())
                }
                let y = vm_access_buffer_contents(buffers, line.args[1].clone(), line.args[2].clone());
                if let Some(x) = buffers.get_mut(&(line.args[0].clone())) {
                    x.contents.extend(y);
                    gas_used += x.contents.len() as f64 / 10.0;
                }
                gas_used += 2.0;
            }
            "Slice" => unsafe {
                if !vm_check_buffer_initialization(buffers, line.args[0].clone())
                    || !vm_check_buffer_initialization(buffers, line.args[1].clone())
                    || !vm_check_buffer_initialization(buffers, line.args[2].clone())
                {
                    vm_throw_local_error(buffers, line.args[3].clone())
                }
                let start_buf = (*vm_access_buffer(buffers, line.args[1].clone(), line.args[3].clone())).as_u64().unwrap() as usize;
                let end_buf = (*vm_access_buffer(buffers, line.args[2].clone(), line.args[3].clone())).as_u64().unwrap() as usize;
                let buf_to_slice =
                    vm_access_buffer_contents(buffers, line.args[0].clone(), line.args[3].clone());
                let sliced_buf = buf_to_slice[start_buf..end_buf].to_vec();
                if let Some(x) = buffers.get_mut(&(line.args[0].clone())) {
                    x.contents = sliced_buf;
                    gas_used += x.contents.len() as f64 / 10.0;
                }
                gas_used += 2.0;
            }
            "Shiftl" => unsafe {
                if !vm_check_buffer_initialization(buffers, line.args[0].clone())
                    || !vm_check_buffer_initialization(buffers, line.args[1].clone())
                {
                    vm_throw_local_error(buffers, line.args[2].clone())
                }
                let mut buf_to_shift =
                    vm_access_buffer_contents(buffers, line.args[0].clone(), line.args[2].clone());
                let shift_amount = (*vm_access_buffer(buffers, line.args[1].clone(), line.args[2].clone())).as_u64().unwrap() as usize;
                buf_to_shift.drain(buf_to_shift.len() - shift_amount..);
                let zeroes = vec![0; shift_amount];
                buf_to_shift.extend(zeroes);
                if let Some(x) = buffers.get_mut(&(line.args[0].clone())) {
                    x.contents = buf_to_shift;
                }
                gas_used += 2.0;
            }
            "Shiftr" => unsafe {
                if !vm_check_buffer_initialization(buffers, line.args[0].clone())
                    || !vm_check_buffer_initialization(buffers, line.args[1].clone())
                {
                    vm_throw_local_error(buffers, line.args[1].clone())
                }
                let mut buf_to_shift =
                    vm_access_buffer_contents(buffers, line.args[0].clone(), line.args[2].clone());
                let shift_amount = (*vm_access_buffer(buffers, line.args[1].clone(), line.args[2].clone())).as_u64().unwrap() as usize;
                buf_to_shift.drain(0..shift_amount);
                let zeroes = vec![0; shift_amount];
                buf_to_shift.splice(..0, zeroes.iter().cloned());
                if let Some(x) = buffers.get_mut(&(line.args[0].clone())) {
                    x.contents = buf_to_shift;
                }
                gas_used += 2.0;
            }
            "Jmp" => {
                line_number = line.args[0].parse::<usize>().unwrap() - 1;
                should_increment = false;
                gas_used += 0.8;
            }
            "JmpCond" => {
                if !vm_check_buffer_initialization(buffers, line.args[0].clone()) {
                    vm_throw_local_error(buffers, line.args[2].clone())
                }
                if buffers.get(&line.args[0]).unwrap().as_u64() != Ok(0) {
                    line_number = line.args[1].parse::<usize>().unwrap() - 1;
                    should_increment = false;
                }
                gas_used += 1.0;
            }
            "Call" => {
                stack.push(buffers);
                origins.push(line_number + 1);
                line_number = line.args[0].parse::<usize>().unwrap() - 1;
                should_increment = false;
            }
            "Ret" => unsafe {
                let frame_buffers = stack.pop();
                let return_value_tmp = (*vm_access_buffer(buffers, "00000001".parse().unwrap(), "00000000".parse().unwrap())).clone();
                *buffers = frame_buffers;
                buffers.insert("00000001".into(), return_value_tmp);
                line_number = *origins.last().expect("Could not get last origin");
                origins.remove(origins.len() - 1);
                should_increment = false;
            }
            "Stdout" => {
                if !vm_check_buffer_initialization(buffers, line.args[0].clone()) {
                    vm_throw_local_error(buffers, line.args[1].clone())
                }
                println!("{:?}", vm_access_buffer_contents(buffers, line.args[0].clone(), line.args[1].clone()));
                gas_used += 0.5;
            }
            "PrintStr" => unsafe {
                if !vm_check_buffer_initialization(buffers, line.args[0].clone()) {
                    vm_throw_local_error(buffers, line.args[1].clone())
                }
                let str = std::str::from_utf8(&(*vm_access_buffer(buffers, line.args[0].clone(), line.args[1].clone())).contents).unwrap();
                println!("{}", str);
                gas_used += 0.5;
            }
            "Stderr" => {
                if !vm_check_buffer_initialization(buffers, line.args[0].clone()) {
                    vm_throw_local_error(buffers, line.args[1].clone())
                }
                eprintln!("{:?}", vm_access_buffer_contents(buffers, line.args[0].clone(), line.args[1].clone()));
                gas_used += 0.5;
            }
            "SetCnst" => unsafe {
                if !vm_check_buffer_initialization(buffers, line.args[0].clone()) {
                    vm_throw_local_error(buffers, line.args[2].clone())
                }
                let bfr_ptr: *mut Buffer = vm_access_buffer(buffers, line.args[0].clone(), line.args[2].clone());
                (*bfr_ptr).contents =
                    hex::decode(line.args[1].clone()).expect("Failed to parse raw hex value");
                gas_used += (*bfr_ptr).contents.len() as f64 / 10.0;
                gas_used += 2.0;
            }
            "Tx" => {
                let sender_bytes =
                    vm_access_buffer_contents(buffers, line.args[0].clone(), line.args[3].clone());
                let sender = match std::string::String::from_utf8(sender_bytes) {
                    Ok(v) => v,
                    Err(_) => {
                        vm_throw_local_error(buffers, line.args[3].clone());
                        "".to_owned()
                    }
                };
                let receiver_bytes =
                    vm_access_buffer_contents(buffers, line.args[1].clone(), line.args[3].clone());
                let receiver = match std::string::String::from_utf8(receiver_bytes) {
                    Ok(v) => v,
                    Err(_) => {
                        vm_throw_local_error(buffers, line.args[3].clone());
                        "".to_owned()
                    }
                };
                let amount_bytes =
                    vm_access_buffer_contents(buffers, line.args[2].clone(), line.args[3].clone());
                let amount = match std::string::String::from_utf8(amount_bytes) {
                    Ok(v) => v,
                    Err(_) => {
                        vm_throw_local_error(buffers, line.args[3].clone());
                        "".to_owned()
                    }
                };
                println!("TX {} {} {}", sender, receiver, amount);
                gas_used += 4.0;
            }
            "GetNthBlk" => unsafe {
                // Get a property of the nth block in the chain
                if !vm_check_buffer_initialization(buffers, line.args[0].clone())
                    || !vm_check_buffer_initialization(buffers, line.args[1].clone())
                    || !vm_check_buffer_initialization(buffers, line.args[2].clone())
                {
                    vm_throw_local_error(buffers, line.args[3].clone())
                }
                let block_number = (*vm_access_buffer(buffers, line.args[0].clone(), line.args[3].clone())).as_u64().unwrap() as usize;
                let property_u64 = (*vm_access_buffer(buffers, line.args[1].clone(), line.args[3].clone())).as_u64().unwrap() as usize;
                let property = match property_u64 {
                    0 => "hash".into(),
                    1 => "prev_hash".into(),
                    2 => "transaction_count".into(),
                    _ => {
                        vm_throw_local_error(buffers, line.args[3].clone());
                        "hash".into()
                    }
                };
                let (result, success) =
                    blockutil_interface.get_nth_block_property(block_number as i64, property);
                if !success {
                    vm_throw_local_error(buffers, line.args[3].clone());
                }
                match property_u64 {
                    0 => {
                        if let Some(x) = buffers.get_mut(&(line.args[2].clone())) {
                            x.contents =
                                hex::decode(result).expect("Failed to parse raw hex value");
                        }
                    }
                    1 => {
                        if let Some(x) = buffers.get_mut(&(line.args[2].clone())) {
                            x.contents =
                                hex::decode(result).expect("Failed to parse raw hex value");
                        }
                    }
                    2 => {
                        println!("Transaction count: {}", result);
                        let count_u64 = result.parse::<u64>().unwrap();
                        if let Some(x) = buffers.get_mut(&(line.args[2].clone())) {
                            x.load_u64(count_u64);
                        }
                    }
                    _ => {
                        vm_throw_local_error(buffers, line.args[1].clone());
                    }
                }
                gas_used += 3.0;
            }
            "GetNthTx" => {
                // Get a property of the nth transaction in the nth block in the chain
                if !vm_check_buffer_initialization(buffers, line.args[0].clone())
                    || !vm_check_buffer_initialization(buffers, line.args[1].clone())
                    || !vm_check_buffer_initialization(buffers, line.args[2].clone())
                    || !vm_check_buffer_initialization(buffers, line.args[3].clone())
                {
                    vm_throw_local_error(buffers, line.args[4].clone())
                }
                let block_number = buffers.get(&line.args[0]).unwrap().as_u64().unwrap() as usize;
                let transaction_number =
                    buffers.get(&line.args[1]).unwrap().as_u64().unwrap() as usize;
                let property_u64 = buffers.get(&line.args[2]).unwrap().as_u64().unwrap() as usize;
                let property = match property_u64 {
                    0 => "sender".into(),
                    1 => "receiver".into(),
                    2 => "amount".into(),
                    3 => "body".into(),
                    _ => {
                        vm_throw_local_error(buffers, line.args[4].clone());
                        "sender".into()
                    }
                };
                let (result, success) = blockutil_interface.get_nth_transaction_property(
                    block_number as i64,
                    transaction_number as i64,
                    property,
                );
                if !success {
                    vm_throw_local_error(buffers, line.args[4].clone());
                }
                match property_u64 {
                    0 => {
                        if let Some(x) = buffers.get_mut(&(line.args[3].clone())) {
                            x.contents =
                                hex::decode(result).expect("Failed to parse raw hex value");
                        }
                    }
                    1 => {
                        if let Some(x) = buffers.get_mut(&(line.args[3].clone())) {
                            x.contents =
                                hex::decode(result).expect("Failed to parse raw hex value");
                        }
                    }
                    2 => {
                        let amount_u64 = result.parse::<u64>().unwrap();
                        if let Some(x) = buffers.get_mut(&(line.args[3].clone())) {
                            x.load_u64(amount_u64);
                        }
                    }
                    3 => {
                        if let Some(x) = buffers.get_mut(&(line.args[3].clone())) {
                            x.contents =
                                hex::decode(result).expect("Failed to parse raw hex value");
                        }
                    }
                    _ => {
                        vm_throw_local_error(buffers, line.args[2].clone());
                    }
                }
                gas_used += 3.0;
            }
            "ChainLen" => {
                if !vm_check_buffer_initialization(buffers, line.args[0].clone()) {
                    vm_throw_local_error(buffers, line.args[1].clone());
                }
                let buf = buffers.get_mut(&line.args[0]).unwrap();
                let len = blockutil_interface.get_blockchain_len();
                buf.load_u64(len);
            }
            "UpdateState" => {
                if !vm_check_buffer_initialization(buffers, line.args[0].clone())
                    || !vm_check_buffer_initialization(buffers, line.args[1].clone())
                {
                    vm_throw_local_error(buffers, line.args[2].clone());
                }
                let location = buffers.get(&line.args[0]).unwrap().as_u64().unwrap() as usize;
                let contents_vec_u8 =
                    vm_access_buffer_contents(buffers, line.args[1].clone(), line.args[2].clone());
                let contents_hex = hex::encode(contents_vec_u8.clone());
                println!("State change: {}{}|{}", contract_hash.clone(), location, contents_hex);
                gas_used += 3.0;
                gas_used += 0.6 * contents_vec_u8.len() as f64;
            }
            "UpdateStateExternal" => {
                if !vm_check_buffer_initialization(buffers, line.args[0].clone())
                    || !vm_check_buffer_initialization(buffers, line.args[1].clone())
                {
                    vm_throw_local_error(buffers, line.args[2].clone());
                }
                let location = buffers.get(&line.args[0]).unwrap().as_u64().unwrap() as usize;
                let contents_vec_u8 =
                    vm_access_buffer_contents(buffers, line.args[1].clone(), line.args[2].clone());
                let contents_hex = hex::encode(contents_vec_u8.clone());
                println!("External state change: {}|{}", location, contents_hex);
                gas_used += 3.0;
                gas_used += 0.6 * contents_vec_u8.len() as f64;
            }
            "GetFromState" => {
                if !vm_check_buffer_initialization(buffers, line.args[0].clone())
                    || !vm_check_buffer_initialization(buffers, line.args[1].clone())
                {
                    vm_throw_local_error(buffers, line.args[2].clone());
                }
                let location = std::str::from_utf8(&vm_access_buffer_contents(buffers, line.args[0].clone(), line.args[1].clone())).expect("Could not convert state location to string").to_owned();
                let (contents_vec_u8, success) = blockutil_interface.get_from_state(location.parse().unwrap());
                if !success {
                    vm_throw_local_error(buffers, line.args[2].clone());
                }
                let dst_buffer = buffers.get_mut(&line.args[1]).unwrap();
                dst_buffer.contents = contents_vec_u8;
            }
            "GetFromStateWithPrefix" => {
                if !vm_check_buffer_initialization(buffers, line.args[0].clone())
                    || !vm_check_buffer_initialization(buffers, line.args[1].clone())
                {
                    vm_throw_local_error(buffers, line.args[2].clone());
                }
                let location = std::str::from_utf8(&vm_access_buffer_contents(buffers, line.args[0].clone(), line.args[1].clone())).expect("Could not convert state location to string").to_owned();
                let (contents_vec_u8, success) = blockutil_interface.get_from_state(contract_hash.clone() + &*location);
                if !success {
                    vm_throw_local_error(buffers, line.args[2].clone());
                }
                let dst_buffer = buffers.get_mut(&line.args[1]).unwrap();
                dst_buffer.contents = contents_vec_u8;
            }
            "QueryOracle" => unsafe {
                if !vm_check_buffer_initialization(buffers, line.args[0].clone())
                    || !vm_check_buffer_initialization(buffers, line.args[1].clone())
                    || !vm_check_buffer_initialization(buffers, line.args[2].clone())
                {
                    vm_throw_local_error(buffers, line.args[3].clone());
                }
                let query_type = (*vm_access_buffer(buffers, line.args[0].clone(), line.args[3].clone())).as_u64().unwrap();
                let query_body = vm_access_buffer_contents(buffers, line.args[1].clone(), line.args[3].clone());
                let query_result = blockutil_interface.query_oracle(query_type, query_body);
                if !query_result.1 {
                    vm_throw_local_error(buffers, line.args[3].clone());
                }
                let res_bfr = vm_access_buffer(buffers, line.args[0].clone(), line.args[3].clone());
                (*res_bfr).contents = query_result.0;
            }
            &_ => vm_throw_global_error(buffers),
        }
        if should_increment {
            line_number += 1;
        }
    }
    (0, gas_used)
}
