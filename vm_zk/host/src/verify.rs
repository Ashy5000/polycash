use contracts::vm::ZkInfo;
use methods::POLYCASH_ZK_GUEST_ID;
use risc0_zkvm::Receipt;
pub(crate) fn verify(receipt: Receipt, expected_merkle_root: String) -> bool {
    let res = receipt.verify(POLYCASH_ZK_GUEST_ID);
    let output: ZkInfo = receipt.journal.decode().unwrap();
    println!("Output: {} Input: {}", output.merkle_root, expected_merkle_root);
    res.is_ok() && output.merkle_root == expected_merkle_root
}
