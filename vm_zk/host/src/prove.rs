use contracts::vm::{VmRunDetails, ZkInfo};
use risc0_zkvm::{default_prover, ExecutorEnv, Receipt};
use methods::{
    POLYCASH_ZK_GUEST_ELF
};
use crate::lazy_vector::HostVector;

pub(crate) fn prove(run_details: VmRunDetails, host_vector: HostVector<i32>) -> Receipt {
    let env = ExecutorEnv::builder()
        .write(&run_details)
        .unwrap()
        .io_callback(zk_common::lazy_vector::SYS_VECTOR_ORACLE, host_vector.handle_guest_request())
        .build()
        .unwrap();

    // Obtain the default prover.
    let prover = default_prover();

    // Proof information by proving the specified ELF binary.
    // This struct contains the receipt along with statistics about execution of the guest
    let prove_info = prover
        .prove(env, POLYCASH_ZK_GUEST_ELF)
        .unwrap();

    // Extract the receipt
    let receipt = prove_info.receipt;

    // Get output from journal
    let output: ZkInfo = receipt.journal.decode().unwrap();

    // Verify input
    let input = output.input;
    assert_eq!(run_details.contract_contents, input.contract_contents);
    assert_eq!(run_details.contract_hash, input.contract_hash);
    assert_eq!(run_details.gas_limit, input.gas_limit);
    assert_eq!(run_details.sender, input.sender);
    assert_eq!(run_details.pending_state.data, input.pending_state.data);

    println!("Gas used: {}", output.gas_used);

    // Return receipt
    receipt
}