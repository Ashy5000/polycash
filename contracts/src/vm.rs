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
    syntax_tree::{SyntaxTree, Line, build_syntax_tree, self}, state::StateManager,
};
use crate::stack::Stack;
use smartstring::alias::String;
use ring::digest;

pub const VM_NIL: *const c_void = 0x0000 as *const c_void;

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

pub struct VmExitDetails {
    exit_code: i64,
    gas_used: f64,
}

pub struct VmInstructionResult {
    exit_details: Option<VmExitDetails>,
    next_pc: usize,
}

pub fn vm_execute_instruction(
    line: Line,
    buffers: &mut FxHashMap<String, Buffer>,
    blockutil_interface: BlockUtilInterface,
    contract_hash: String,
    pc: usize,
    gas_used: &mut f64,
    stack: &mut Stack,
    state_manager: &mut StateManager,
    gas_limit: f64,
) -> VmInstructionResult {
    match line.command.as_str() {
        "Exit" => {
            *gas_used += 1.0;
            return VmInstructionResult {
                exit_details: Some(VmExitDetails {
                    exit_code: line.args[0].parse::<i64>().unwrap(),
                    gas_used: *gas_used,
                }),
                next_pc: 0,
            };
        }
        "ExitBfr" => {
            *gas_used += 1.0;
            if !vm_check_buffer_initialization(buffers, line.args[0].clone()) {
                vm_throw_local_error(buffers, line.args[1].clone())
            }
            let exit_code = buffers
                .get(&line.args[0].clone())
                .expect("Unable to read exit code buffer")
                .as_u64()
                .expect("Could not convert exit code to u64");
            VmInstructionResult {
                exit_details: Some(VmExitDetails{
                    exit_code: exit_code as i64,
                    gas_used: *gas_used
                }),
                next_pc: 0,
            }
        }
        "InitBfr" => {
            buffers.insert(
                line.args[0].clone(),
                Buffer {
                    contents: Vec::new(),
                },
            );
            *gas_used += 2.0;
            *gas_used += line.args[0].len() as f64 / 10.0;
            VmInstructionResult {
                exit_details: None,
                next_pc: pc + 1,
            }
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
                *gas_used += 2.0;
            }
            let src_contents: Vec<u8> =
                vm_access_buffer_contents(buffers, line.args[0].clone(), line.args[2].clone());
            if let Some(dst) = buffers.get_mut(&(line.args[1].clone())) {
                dst.contents = src_contents.clone();
                *gas_used += src_contents.len() as f64 / 10.0;
            }
            *gas_used += 2.0;
            VmInstructionResult {
                exit_details: None,
                next_pc: pc + 1,
            }
        }
        "FreeBfr" => {
            if !vm_check_buffer_initialization(buffers, line.args[0].clone()) {
                vm_throw_local_error(buffers, line.args[1].clone())
            }
            buffers
                .remove(&line.args[0].clone())
                .expect("Failed to free memory");
            *gas_used += 2.0;
            VmInstructionResult {
                exit_details: None,
                next_pc: pc + 1,
            }
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
            *gas_used += 1.0;
            VmInstructionResult {
                exit_details: None,
                next_pc: pc + 1,
            }
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
            VmInstructionResult {
                exit_details: None,
                next_pc: pc + 1,
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
            *gas_used += 0.5;
            VmInstructionResult {
                exit_details: None,
                next_pc: pc + 1,
            }
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
            *gas_used += 0.5;
            VmInstructionResult {
                exit_details: None,
                next_pc: pc + 1,
            }
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
            *gas_used += 1.0;
            VmInstructionResult {
                exit_details: None,
                next_pc: pc + 1,
            }
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
            *gas_used += 1.0;
            VmInstructionResult {
                exit_details: None,
                next_pc: pc + 1,
            }
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
            *gas_used += 1.5;
            VmInstructionResult {
                exit_details: None,
                next_pc: pc + 1,
            }
        }
        "Mod" => {
            execute_math_operation(
                Modulo {},
                buffers,
                line.args[0].clone(),
                line.args[1].clone(),
                line.args[2].clone(),
                line.args[3].clone(),
            );
            *gas_used += 1.0;
            VmInstructionResult {
                exit_details: None,
                next_pc: pc + 1,
            }
        },
        "Eq" => {
            execute_math_operation(
                Eq {},
                buffers,
                line.args[0].clone(),
                line.args[1].clone(),
                line.args[2].clone(),
                line.args[3].clone(),
            );
            *gas_used += 0.5;
            VmInstructionResult {
                exit_details: None,
                next_pc: pc + 1,
            }
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
            *gas_used += 0.5;
            VmInstructionResult {
                exit_details: None,
                next_pc: pc + 1,
            }
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
            *gas_used += 1.0;
            VmInstructionResult {
                exit_details: None,
                next_pc: pc + 1,
            }
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
            *gas_used += 1.0;
            VmInstructionResult {
                exit_details: None,
                next_pc: pc + 1,
            }
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
            *gas_used += 0.5;
            VmInstructionResult {
                exit_details: None,
                next_pc: pc + 1,
            }
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
                *gas_used += x.contents.len() as f64 / 10.0;
            }
            *gas_used += 2.0;
            VmInstructionResult {
                exit_details: None,
                next_pc: pc + 1,
            }
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
                *gas_used += x.contents.len() as f64 / 10.0;
            }
            *gas_used += 2.0;
            VmInstructionResult {
                exit_details: None,
                next_pc: pc + 1,
            }
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
            *gas_used += 2.0;
            VmInstructionResult {
                exit_details: None,
                next_pc: pc + 1,
            }
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
            *gas_used += 2.0;
            VmInstructionResult {
                exit_details: None,
                next_pc: pc + 1,
            }
        }
        "Jmp" => {
            *gas_used += 0.8;
            VmInstructionResult {
                exit_details: None,
                next_pc: line.args[0].parse::<usize>().unwrap() - 1,
            }
        }
        "JmpCond" => {
            if !vm_check_buffer_initialization(buffers, line.args[0].clone()) {
                vm_throw_local_error(buffers, line.args[2].clone())
            }
            *gas_used += 1.0;
            if buffers.get(&line.args[0]).unwrap().as_u64() != Ok(0) {
                VmInstructionResult {
                    exit_details: None,
                    next_pc: line.args[1].parse::<usize>().unwrap() - 1,
                }
            } else {
                VmInstructionResult {
                    exit_details: None,
                    next_pc: pc + 1
                }
            }
        }
        "Call" => {
            stack.push(buffers, pc + 1);
            VmInstructionResult {
                exit_details: None,
                next_pc: line.args[0].parse::<usize>().unwrap() - 1,
            }
        }
        "Ret" => unsafe {
            let frame = stack.pop();
            let return_value_tmp = (*vm_access_buffer(buffers, "00000001".parse().unwrap(), "00000000".parse().unwrap())).clone();
            *buffers = frame.buffers;
            buffers.insert("00000001".into(), return_value_tmp);
            VmInstructionResult {
                exit_details: None,
                next_pc: frame.origin,
            }
        }
        "Stdout" => {
            if !vm_check_buffer_initialization(buffers, line.args[0].clone()) {
                vm_throw_local_error(buffers, line.args[1].clone())
            }
            println!("{:?}", vm_access_buffer_contents(buffers, line.args[0].clone(), line.args[1].clone()));
            *gas_used += 0.5;
            VmInstructionResult {
                exit_details: None,
                next_pc: pc + 1,
            }
        }
        "PrintStr" => unsafe {
            if !vm_check_buffer_initialization(buffers, line.args[0].clone()) {
                vm_throw_local_error(buffers, line.args[1].clone())
            }
            let str = std::str::from_utf8(&(*vm_access_buffer(buffers, line.args[0].clone(), line.args[1].clone())).contents).unwrap();
            println!("{}", str);
            *gas_used += 0.5;
            VmInstructionResult {
                exit_details: None,
                next_pc: pc + 1,
            }
        }
        "Stderr" => {
            if !vm_check_buffer_initialization(buffers, line.args[0].clone()) {
                vm_throw_local_error(buffers, line.args[1].clone())
            }
            eprintln!("{:?}", vm_access_buffer_contents(buffers, line.args[0].clone(), line.args[1].clone()));
            *gas_used += 0.5;
            VmInstructionResult {
                exit_details: None,
                next_pc: pc + 1,
            }
        }
        "SetCnst" => unsafe {
            if !vm_check_buffer_initialization(buffers, line.args[0].clone()) {
                vm_throw_local_error(buffers, line.args[2].clone())
            }
            let bfr_ptr: *mut Buffer = vm_access_buffer(buffers, line.args[0].clone(), line.args[2].clone());
            (*bfr_ptr).contents =
                hex::decode(line.args[1].clone()).expect("Failed to parse raw hex value");
            *gas_used += (*bfr_ptr).contents.len() as f64 / 10.0;
            *gas_used += 2.0;
            VmInstructionResult {
                exit_details: None,
                next_pc: pc + 1,
            }
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
            *gas_used += 4.0;
            VmInstructionResult {
                exit_details: None,
                next_pc: pc + 1,
            }
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
            *gas_used += 3.0;
            VmInstructionResult {
                exit_details: None,
                next_pc: pc + 1,
            }
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
            *gas_used += 3.0;
            VmInstructionResult {
                exit_details: None,
                next_pc: pc + 1,
            }
        }
        "ChainLen" => {
            if !vm_check_buffer_initialization(buffers, line.args[0].clone()) {
                vm_throw_local_error(buffers, line.args[1].clone());
            }
            let buf = buffers.get_mut(&line.args[0]).unwrap();
            let len = blockutil_interface.get_blockchain_len();
            buf.load_u64(len);
            *gas_used += 2.0;
            VmInstructionResult {
                exit_details: None,
                next_pc: pc + 1,
            }
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
            let full_location = format!("{}{}", contract_hash.clone(), location);
            state_manager.write(full_location, contents_vec_u8.clone());
            *gas_used += 3.0;
            *gas_used += 0.6 * contents_vec_u8.len() as f64;
            VmInstructionResult {
                exit_details: None,
                next_pc: pc + 1,
            }
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
            state_manager.write(format!("{}", location), contents_vec_u8.clone());
            *gas_used += 3.0;
            *gas_used += 0.6 * contents_vec_u8.len() as f64;
            VmInstructionResult {
                exit_details: None,
                next_pc: pc + 1,
            }
        }
        "GetFromState" => unsafe {
            if !vm_check_buffer_initialization(buffers, line.args[0].clone())
                || !vm_check_buffer_initialization(buffers, line.args[1].clone())
            {
                vm_throw_local_error(buffers, line.args[2].clone());
            }
            let location = (*vm_access_buffer(buffers, line.args[0].clone(), line.args[1].clone())).as_u64().unwrap() as usize;
            let location = format!("{}{}", contract_hash.clone(), location);
            let contents_vec_u8 = state_manager.get(location).unwrap();
            let dst_buffer = buffers.get_mut(&line.args[1]).unwrap();
            dst_buffer.contents = contents_vec_u8;
            *gas_used += 2.0;
            VmInstructionResult {
                exit_details: None,
                next_pc: pc + 1,
            }
        }
        "GetFromStateExternal" => {
            if !vm_check_buffer_initialization(buffers, line.args[0].clone())
                || !vm_check_buffer_initialization(buffers, line.args[1].clone())
            {
                vm_throw_local_error(buffers, line.args[2].clone());
            }
            let location_vec_u8 = &vm_access_buffer_contents(buffers, line.args[0].clone(), line.args[1].clone());
            let location = hex::encode(location_vec_u8);
            let contents_vec_u8: Vec<u8> = state_manager.get(location).unwrap();
            let dst_buffer = buffers.get_mut(&line.args[1]).unwrap();
            dst_buffer.contents = contents_vec_u8;
            *gas_used += 2.0;
            VmInstructionResult {
                exit_details: None,
                next_pc: pc + 1,
            }
        }
        "QueryOracle" => unsafe {
            *gas_used += 10.0;
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
            *gas_used += 10.0;
            VmInstructionResult {
                exit_details: None,
                next_pc: pc + 1,
            }
        }
        "Invoke" => unsafe {
            *gas_used += 3.0;
            let location = line.args[0].parse::<u64>().expect("Failed to parse location");
            let contents = blockutil_interface.read_contract(location).expect("Failed to read invoked contract");
            let mut tree = syntax_tree::build_syntax_tree();
            tree.create(contents.clone());
            let child_hash = hex::encode(digest::digest(&digest::SHA256, contents.as_ref()).as_ref());
            stack.push(buffers, pc + 1);
            let mut child_pc = 0;
            buffers.clear();
            vm_simulate(
                tree,
                buffers,
                stack,
                state_manager,
                gas_used,
                blockutil_interface,
                child_hash.into(),
                gas_limit,
                &mut child_pc
            );
            let frame = stack.pop();
            let return_value_tmp = vm_access_buffer(buffers, "00000001".parse().unwrap(), "00000000".parse().unwrap());
            *buffers = frame.buffers;
            if return_value_tmp != VM_NIL as *mut Buffer {
                buffers.insert("00000001".into(), (*return_value_tmp).clone());
            }
            VmInstructionResult{
                exit_details: None, next_pc: pc + 1
            }
        }
        &_ => {
            vm_throw_global_error(buffers);
            VmInstructionResult{
                exit_details: None, next_pc: pc + 1
            }
        },
    }
}

pub fn vm_simulate(
    syntax_tree: SyntaxTree,
    buffers: &mut FxHashMap<String, Buffer>,
    stack: &mut Stack,
    state_manager: &mut StateManager,
    gas_used: &mut f64,
    blockutil_interface: BlockUtilInterface,
    contract_hash: String,
    gas_limit: f64,
    pc: &mut usize,
) -> (i64, f64) {
    while *pc < syntax_tree.lines.len() {
        if *gas_used > gas_limit {
            return (2, gas_limit);
        }
        let line = &syntax_tree.lines[*pc];
        let res = vm_execute_instruction(line.clone(), buffers, blockutil_interface.clone(), contract_hash.clone(), *pc, gas_used, stack, state_manager, gas_limit);
        if let Some(exit_details) = res.exit_details {
            state_manager.flush();
            return (exit_details.exit_code, exit_details.gas_used);
        }
        *pc = res.next_pc;
    }
    state_manager.flush();
    (0, *gas_used)
}

pub fn run_vm(contract_contents: String, contract_hash: String, gas_limit: f64) -> (i64, f64) {
    let mut tree = build_syntax_tree();
    tree.create(contract_contents);
    let mut buffers: FxHashMap<String, Buffer> = FxHashMap::default();
    buffers.insert(
        "00000000".parse().unwrap(),
        Buffer { contents: vec![] }
    );
    let interface = BlockUtilInterface::new();
    let mut stack = Stack{frames: vec![]};
    let mut pc: usize = 0;
    let mut gas_used = 0.0;
    let mut state_manager = StateManager::new(interface.clone(), contract_hash.to_string());
    vm_simulate(tree, &mut buffers, &mut stack, &mut state_manager, &mut gas_used, interface, contract_hash, gas_limit, &mut pc)
}
