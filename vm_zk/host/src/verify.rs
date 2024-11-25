use methods::POLYCASH_ZK_GUEST_ID;
use risc0_zkvm::Receipt;
pub(crate) fn verify(receipt: Receipt) -> bool {
    let res = receipt.verify(POLYCASH_ZK_GUEST_ID);
    res.is_ok()
}
