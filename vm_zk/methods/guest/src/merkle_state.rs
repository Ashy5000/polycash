use contracts::merkle::{get_from_merkle, MerkleContainer, MerkleNode};
use contracts::state::State;
use crate::lazy_vector::LazyVector;

impl MerkleContainer<Vec<u8>> for LazyVector<MerkleNode<Vec<u8>>> {
    fn get_wrapper(&mut self, index: usize) -> Option<MerkleNode<Vec<u8>>> {
        self.get(index)
    }
}

pub(crate) struct MerkleState {
    pub(crate) contents: LazyVector<MerkleNode<Vec<u8>>>,
    pub(crate) prefix: String
}

impl MerkleState {
    pub(crate) fn new(contents: LazyVector<MerkleNode<Vec<u8>>>, prefix: String) -> Self {
        MerkleState {
            contents,
            prefix
        }
    }
}

impl State for MerkleState {
    fn write(&mut self, location: String, contents: Vec<u8>, out: &mut String) {
        // Same write functionality as OnchainState
        if location.as_bytes().starts_with(&self.prefix.as_bytes()) {
            let string = format!("State change: {}|{}", location, hex::encode(contents));
            println!("{}\n", string);
            out.push_str(&string);
        } else {
            let string = format!("External state change: {}|{}", location, hex::encode(contents));
            println!("{}\n", string);
            out.push_str(&string);
        }
    }

    fn get(&mut self, location: String) -> Result<Vec<u8>, String> {
        Ok(get_from_merkle::<LazyVector<MerkleNode<Vec<u8>>>>(&mut self.contents, location))
    }

    fn dump(&self) -> rustc_hash::FxHashMap<String, Vec<u8>> { panic!("Not implemented."); }
}