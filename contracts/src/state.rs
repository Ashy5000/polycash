use rustc_hash::FxHashMap;

use crate::{buffer::Buffer, blockutil::BlockUtilInterface};

pub trait State {
    fn write(&mut self, location: String, contents: Vec<u8>);
    fn get(&self, location: String) -> Result<Vec<u8>, String>;
}

pub struct CachedState {
    contents: FxHashMap<String, Vec<u8>>
}

impl CachedState {
    fn new() -> Self {
        CachedState { contents: FxHashMap::default() }
    }
}

impl State for CachedState {
    fn write(&mut self, location: String, contents: Vec<u8>) {
        self.contents.insert(location, contents);
    }
    fn get(&self, location: String) -> Result<Vec<u8>, String> {
        if !self.contents.contains_key(&location) {
            Err("Could not find key in cache".parse().unwrap())
        } else {
            Ok(self.contents[&location].clone())
        }
    }
}

pub struct OnchainState {
    blockutil_interface: BlockUtilInterface,
    prefix: String,
}

impl OnchainState {
    fn new(blockutil_interface: BlockUtilInterface, prefix: String) -> Self {
        OnchainState { blockutil_interface, prefix }
    }
}

impl State for OnchainState {
    fn write(&mut self, location: String, contents: Vec<u8>) {
        if location.as_bytes().starts_with(&self.prefix.as_bytes()) {
            println!("State change: {}|{}", location, String::from_utf8(contents).unwrap());
        } else {
            println!("External state change: {}|{}", location, String::from_utf8(contents).unwrap());
        }
    }
    fn get(&self, location: String) -> Result<Vec<u8>, String> {
        let (contents_vec_u8, success) = self.blockutil_interface.get_from_state(location.parse().unwrap());
        if success {
            Ok(contents_vec_u8)
        } else {
            Err("Could not get from onchain state".to_string())
        }
    }
}

pub struct StateManager {
    cached_state: CachedState,
    onchain_state: OnchainState,
}

impl StateManager {
    pub fn new(blockutil_interface: BlockUtilInterface, contract_hash: String) -> Self {
        StateManager {
            cached_state: CachedState::new(),
            onchain_state: OnchainState::new(blockutil_interface, contract_hash),
        }
    }
    pub fn write(&mut self, location: String, contents: Vec<u8>) {
        self.cached_state.write(location, contents);
    }
    pub fn get(&mut self, location: String) -> Result<Vec<u8>, String> {
        if let Ok(res) = self.cached_state.get(location.clone()) {
            return Ok(res);
        } else {
            let onchain_res = self.onchain_state.get(location.clone());
            if onchain_res.is_err() {
                Err("Could not get from onchain state".to_string())
            } else {
                self.cached_state.write(location, onchain_res.clone().unwrap());
                onchain_res
            }
        }
    }
    pub fn flush(&mut self) {
        for (location, contents) in self.cached_state.contents.clone().into_iter() {
            self.onchain_state.write(location, contents);
        }
    }
}
