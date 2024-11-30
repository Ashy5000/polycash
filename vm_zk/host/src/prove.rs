use crate::lazy_vector::HostVector;
use contracts::merkle::MerkleNode;
use contracts::vm::{VmRunDetails, ZkInfo};
use methods::POLYCASH_ZK_GUEST_ELF;
use risc0_zkvm::{default_prover, ExecutorEnv, Receipt};
use crate::socket::Socket;

pub(crate) fn prove(
    run_details: VmRunDetails,
    host_vector: HostVector<MerkleNode<Vec<u8>>>,
    socket: &mut Socket
) -> Receipt {
    let env = ExecutorEnv::builder()
        .write(&run_details)
        .unwrap()
        .io_callback(
            zk_common::lazy_vector::SYS_VECTOR_ORACLE,
            host_vector.handle_guest_request(),
        )
        .build()
        .unwrap();

    // Obtain the default prover.
    let prover = default_prover();

    // Proof information by proving the specified ELF binary.
    // This struct contains the receipt along with statistics about execution of the guest
    let prove_info = prover.prove(env, POLYCASH_ZK_GUEST_ELF).unwrap();

    // Extract the receipt
    let receipt = prove_info.receipt;

    // Get output from journal
    let output: ZkInfo = receipt.journal.decode().unwrap();
    
    // Return receipt
    receipt
}
