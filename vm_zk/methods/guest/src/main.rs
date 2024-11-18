mod lazy_vector;
mod merkle_state;

use contracts::blockutil::{BlockUtilInterface, NodeBlockUtilInterface};
use contracts::merkle::MerkleNode;
use contracts::state::{CachedState, StateManager};
use contracts::vm::{run_vm, VmRunDetails, ZkInfo};
use risc0_zkvm::guest::env;
use smartstring::alias::String;
use crate::lazy_vector::LazyVector;
use crate::merkle_state::MerkleState;

fn main() {
    // Read the input
    let run_details: VmRunDetails = env::read();

    // Extract details
    let contract_contents = run_details.contract_contents.clone();
    let contract_hash = String::from(run_details.contract_hash.clone());
    let gas_limit = run_details.gas_limit.clone();
    let sender = run_details.sender.clone();

    // Setup lazy vector
    let lazy_len = run_details.lazy_len;
    let lazy: LazyVector<MerkleNode> = LazyVector::new(lazy_len);

    // Setup state
    let pending_state = run_details.pending_state.clone();
    let merkle_state = MerkleState::new(lazy, contract_hash.parse().unwrap());
    let mut state_manager = StateManager {
        cached_state: CachedState::new(),
        onchain_state: merkle_state,
        pending_state
    };

    // Setup blockutil interface
    let interface = NodeBlockUtilInterface::new();

    // Run VM
    let (exit_code, gas_used, out) = run_vm(contract_contents.parse().unwrap(), contract_hash, gas_limit, sender, &mut state_manager, interface);

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
