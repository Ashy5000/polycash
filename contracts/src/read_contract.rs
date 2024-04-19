use std::{env, fs};

pub(crate) fn read_contract() -> String {
    let args: Vec<String> = env::args().collect();
    if args.len() != 2 {
        panic!("Expected exactly 1 argument, got {}", args.len())
    }
    let contract_path = &args[1];
    let contents = fs::read_to_string(contract_path)
        .expect("Should have been able to read smart contract file");
    contents
}