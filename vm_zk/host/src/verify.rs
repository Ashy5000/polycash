use contracts::vm::ZkInfo;
use methods::POLYCASH_ZK_GUEST_ID;
use risc0_zkvm::Receipt;
pub(crate) fn verify(receipt: Receipt, expected_merkle_root: String, expected_input_hash: String, expected_transition_hash: String) -> bool {
    let res = receipt.verify(POLYCASH_ZK_GUEST_ID);
    let output: ZkInfo = receipt.journal.decode().unwrap();
    println!("Expected: {}, Actual: {}", expected_transition_hash, output.state_transition_root);
    res.is_ok() && output.merkle_root == expected_merkle_root && output.input_hash == expected_input_hash && output.state_transition_root == expected_transition_hash
}
