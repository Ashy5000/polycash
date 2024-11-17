mod lazy_vector;

use contracts::merkle::MerkleNode;
use contracts::vm::{run_vm, VmRunDetails, ZkInfo};
use risc0_zkvm::guest::env;
use smartstring::alias::String;
use crate::lazy_vector::LazyVector;

fn main() {
    // Read the input
    let run_details: VmRunDetails = env::read();

    // Extract details
    let contract_contents = run_details.contract_contents.clone();
    let contract_hash = String::from(run_details.contract_hash.clone());
    let gas_limit = run_details.gas_limit.clone();
    let sender = run_details.sender.clone();
    let pending_state = run_details.pending_state.clone();
    
    // Setup lazy vector
    let mut lazy: LazyVector<MerkleNode> = LazyVector::new(3);
    println!("Hash: {}", lazy.get(0).unwrap().hash);

    // Run VM
    let (exit_code, gas_used, out) = run_vm(contract_contents.parse().unwrap(), contract_hash, gas_limit, sender, pending_state);

    // Format output
    let output = ZkInfo {
        exit_code,
        gas_used,
        input: run_details.clone(),
        out: std::string::String::from(out)
    };

    // Write public output to the journal
    env::commit(&output);
}
