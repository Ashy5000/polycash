use serde::{Deserialize, Serialize};
use rustc_hash::FxHashMap;
use sha2::{Digest, Sha256};

#[derive(Clone, Serialize, Deserialize)]
pub struct MerkleNode<T> {
    pub hash: String,
    pub value: T,
    pub children: FxHashMap<char, usize>
}

fn hash_node<T: Clone + AsRef<[u8]> + >(tree: &mut Vec<MerkleNode<T>>, node_index: usize) {
    let mut hasher = Sha256::new();
    hasher.update(tree[node_index].value.clone());
    for (_, index) in tree[node_index].children.iter() {
        hasher.update(tree[index.clone()].hash.clone());
    }
    let hash = hasher.finalize();
    tree[node_index].hash = hex::encode(hash);
}

fn hash_tree<T: Clone + AsRef<[u8]>>(tree: &mut Vec<MerkleNode<T>>, index: usize) {
    if tree[index].children.len() == 0 {
        // Base case
        hash_node(tree, index);
    } else {
        // Recursive case
        for (_, child_index) in tree[index].children.clone().iter() {
            hash_tree(tree, *child_index);
            hash_node(tree, index);
        }
    }
}

pub fn merklize<T: Default + Clone + AsRef<[u8]>>(state: FxHashMap<String, T>) -> Vec<MerkleNode<T>> {
    let mut tree: Vec<MerkleNode<T>> = vec![MerkleNode{
        hash: String::new(),
        value: T::default(),
        children: FxHashMap::default()
    }];
    for (loc, val) in state.iter() {
        let mut active: usize = 0;
        for c in loc.chars() {
            if tree[active].children.contains_key(&c) {
                active = tree[active].children[&c];
                continue;
            }
            tree.push(MerkleNode{
                hash: String::new(),
                value: T::default(),
                children: FxHashMap::default()
            });
            let new_index = tree.len() - 1;
            tree[active].children.insert(c, new_index);
            active = new_index;
        }
        tree[active].value = val.clone();
    }
    hash_tree(&mut tree, 0);
    tree
}

pub trait MerkleContainer<T> {
    fn get_wrapper(&mut self, index: usize) -> Option<MerkleNode<T>>;
}

pub fn get_from_merkle<A : MerkleContainer<B> + Clone, B : Clone>(tree: &mut A, loc: String) -> B {
    let mut active: usize = 0;
    for c in loc.chars() {
        active = tree.get_wrapper(active).unwrap().children[&c];
    }
    tree.get_wrapper(active).unwrap().value.clone()
}