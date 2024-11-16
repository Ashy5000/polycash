use risc0_zkvm::guest::env;
use serde::{Serialize, Deserialize};

#[derive(Serialize, Deserialize)]
pub struct LazyVector<T> {
    len: usize,
    #[serde(skip)]
    elements: Vec<Option<T>>
}

impl<T: Serialize + Clone + for<'a> Deserialize<'a>> LazyVector<T> {
    pub fn new(len: usize) -> Self {
        Self {
            len,
            elements: vec![None; len]
        }
    }

    pub fn len(&self) -> usize {
        self.len
    }

    pub fn get(&mut self, index: usize) -> Option<T> {
        if index >= self.len {
            return None;
        }

        if self.elements[index].is_none() {
            let element: T = bincode::deserialize(env::send_recv_slice::<usize, u8>(zk_common::lazy_vector::SYS_VECTOR_ORACLE, &[index])).unwrap();
            self.elements[index] = Some(element);
        }
        
        self.elements[index].clone()
    }
}