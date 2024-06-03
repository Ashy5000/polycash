use std::process::Command;
// Copyright 2024, Asher Wrobel
/*
This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.
*/
use crate::sanitization::sanitize_node_console_command;

pub struct BlockUtilInterface {
    node_executable_path: String,
}

impl BlockUtilInterface {
    pub fn new() -> Self {
        // Read the executable path from the executable path file
        let node_executable_path = std::fs::read_to_string("node_executable_path.txt")
            .expect("Could not read node_executable_path.txt");
        // Remove the newline character
        // This is necessary because the path is read from a file
        let node_executable_path = node_executable_path.trim().to_string();
        Self {
            node_executable_path,
        }
    }
    pub fn get_nth_block_property(&self, n: i64, property: String) -> (String, bool) {
        // Run the node script to get a property of the nth block
        let command = format!("sync;getNthBlock {} {}", n, property);
        if !sanitize_node_console_command(&command) {
            println!("Forbidden command");
            return ("".to_string(), false);
        }
        let output = Command::new(self.node_executable_path.clone())
            .arg("--command")
            .arg(command)
            .output();
        let output = output.expect("Failed to execute node script");
        let output = String::from_utf8(output.stdout).expect("Failed to convert output to string");
        // Remove the newline character
        let output = output.trim().to_string();
        // Split the output by the newline character
        let output = output.split("\n").collect::<Vec<&str>>();
        // Get the last element of the output
        // This is neccessary because the previous lines contain logs
        let output = output[output.len() - 1].to_string();
        (output, true)
    }
    pub fn get_nth_transaction_property(
        &self,
        block_pos: i64,
        transaction_pos: i64,
        property: String,
    ) -> (String, bool) {
        // Run the node script to get a property of the nth transaction
        let command = format!(
            "sync;getNthTransaction {} {} {}",
            block_pos, transaction_pos, property
        );
        if !sanitize_node_console_command(&command) {
            println!("Forbidden command");
            return ("".to_string(), false);
        }
        let output = Command::new(self.node_executable_path.clone())
            .arg("--command")
            .arg(command)
            .output();
        let output =
            String::from_utf8(output.unwrap().stdout).expect("Failed to convert output to string");
        // Remove the newline character
        let output = output.trim().to_string();
        // Split the output by the newline character
        let output = output.split("\n").collect::<Vec<&str>>();
        // Get the last element of the output
        // This is neccessary because the previous lines contain logs
        let output = output[output.len() - 1].to_string();
        (output, true)
    }
    pub fn get_from_state(&self, property: u64) -> (Vec<u8>, bool) {
        let command = format!("sync;getFromState {}", property);
        if !sanitize_node_console_command(&command) {
            println!("Forbidden command");
            return (vec![], false);
        }
        let output = Command::new(self.node_executable_path.clone())
            .arg("--command")
            .arg(command)
            .output();
        let output =
            String::from_utf8(output.unwrap().stdout).expect("Failed to convert output to string");
        let output = output.trim().to_string();
        let output = output.split("\n").collect::<Vec<&str>>();
        let output = output[output.len() - 1].to_string();
        println!("output: {}", output);
        let data_hex: String = output.chars().skip(6).collect();
        println!("data_hex: {}", data_hex);
        let data = hex::decode(data_hex).expect("Decoding failed");
        (data, true)
    }
    pub fn get_blockchain_len(&self) -> u64 {
        let command = "sync;getBlockchainLen";
        let output = Command::new(self.node_executable_path.clone())
            .arg("--command")
            .arg(command)
            .output();
        let output =
            String::from_utf8(output.unwrap().stdout).expect("Failed to convert output to string");
        let output = output.trim().to_string();
        let output = output.split("\n").collect::<Vec<&str>>();
        let output = output[output.len() - 1].to_string();
        let len = output.parse::<u64>().unwrap(); // Long Live the Turbofish.
        len
    }
}
