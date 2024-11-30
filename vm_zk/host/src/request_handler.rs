use std::fs;
use contracts::blockutil::{BlockUtilInterface, NodeBlockUtilInterface};
use contracts::merkle::merklize;
use contracts::read_contract::read_contract;
use contracts::vm::ZkInfo;
use risc0_zkvm::Receipt;
use rustc_hash::FxHashMap;
use crate::lazy_vector::HostVector;
use crate::prove::prove;
use crate::socket::Socket;
use crate::verify::verify;

pub(crate) fn handle_request(data: Vec<u8>, socket: &mut Socket) {
    let s = String::from_utf8(data).unwrap();
    let args: Vec<&str> = s.split(" ").collect();
    if args[0] == "V" {
        // Verify
        let receipt_file = fs::read(args[1]).unwrap();
        let receipt: Receipt = rmp_serde::from_slice(&*receipt_file).unwrap();
        let expected_merkle_root = args[2].to_owned();
        assert!(verify(receipt, expected_merkle_root));
        println!("Verification success!");
        socket.write_message("Verification success!".as_ref()).unwrap();
        return;
    }

    let contracts_file = fs::read_to_string(args[0]).unwrap();
    let contract_contents_str = contracts_file.split("%").collect::<Vec<&str>>(); // % marks separation between contracts
    let mut contract_contents = Vec::new();
    for contract in contract_contents_str {
        contract_contents.push(std::string::String::from(contract));
    }
    let contract_hashes_str = args[1].split("%").collect::<Vec<&str>>();
    let mut contract_hashes = Vec::new();
    for hash in contract_hashes_str {
        contract_hashes.push(std::string::String::from(hash));
    }
    let gas_limits_str = args[2].split("%").collect::<Vec<&str>>();
    let mut gas_limits = Vec::new();
    for limit in gas_limits_str {
        gas_limits.push(limit.parse::<f64>().unwrap() as i64);
    }
    let senders_str: Vec<&str> = args[3].split("%").collect::<Vec<&str>>();
    let mut senders: Vec<Vec<u8>> = Vec::new();
    for sender in senders_str {
        senders.push(sender.into());
    }

    // Initialize merkle tree
    let mut data: FxHashMap<String, Vec<u8>> = FxHashMap::default();
    let merkle_file = fs::read_to_string(args[4]).unwrap();
    if merkle_file.len() != 0 {
        let merkle_pairs: Vec<&str> = merkle_file.split("*").collect();
        for pair in merkle_pairs {
            let segments: Vec<&str> = pair.split(">").collect();
            let key = String::from(segments[0]);
            let value = hex::decode(segments[1].trim()).unwrap();
            data.insert(key, value);
        }
    }
    let tree = merklize(data);
    let lazy_len = tree.len();
    let host_vector = HostVector::new(tree);

    // Create node blockutil for data fetching
    let node_blockutil = NodeBlockUtilInterface::new();

    // Fetch data from node
    let blockchain_len = node_blockutil.get_blockchain_len();

    let run_details = contracts::vm::VmRunDetails {
        contract_contents,
        contract_hash: contract_hashes,
        gas_limits,
        senders,
        lazy_len,
        blockchain_len,
    };

    let receipt = prove(run_details, host_vector, socket);
    let receipt_serialized = rmp_serde::to_vec(&receipt).unwrap();

    let out_file = args[5];
    fs::write(out_file, receipt_serialized).unwrap();
    
    let zk_info: ZkInfo = receipt.journal.decode().unwrap();
    socket.write_message(zk_info.out.as_ref()).unwrap()
}