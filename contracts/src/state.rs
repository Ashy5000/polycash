use crate::blockutil::{BlockUtilInterface, NodeBlockUtilInterface};
use crate::msgpack::PendingState;
use rustc_hash::FxHashMap;

pub trait State {
    fn write(&mut self, location: String, contents: Vec<u8>, out: &mut String);
    fn get(&mut self, location: String) -> Result<Vec<u8>, String>;
    fn dump(&self) -> FxHashMap<String, Vec<u8>>;
}

pub struct CachedState {
    contents: FxHashMap<String, Vec<u8>>,
}

impl CachedState {
    pub fn new() -> Self {
        CachedState {
            contents: FxHashMap::default(),
        }
    }
}

impl State for CachedState {
    fn write(&mut self, location: String, contents: Vec<u8>, _out: &mut String) {
        self.contents.insert(location, contents);
    }
    fn get(&mut self, location: String) -> Result<Vec<u8>, String> {
        if !self.contents.contains_key(&location) {
            Err("Could not find key in cache".parse().unwrap())
        } else {
            Ok(self.contents[&location].clone())
        }
    }

    fn dump(&self) -> FxHashMap<String, Vec<u8>> {
        self.contents.clone()
    }
}

pub struct OnchainState {
    blockutil_interface: NodeBlockUtilInterface,
    prefix: String,
}

impl OnchainState {
    fn new(blockutil_interface: NodeBlockUtilInterface, prefix: String) -> Self {
        OnchainState {
            blockutil_interface,
            prefix,
        }
    }
}

impl State for OnchainState {
    fn write(&mut self, location: String, contents: Vec<u8>, out: &mut String) {
        if location.as_bytes().starts_with(&self.prefix.as_bytes()) {
            let string = format!(
                "State change: {}|{}",
                location,
                String::from_utf8(contents).unwrap()
            );
            println!("{}\n", string);
            out.push_str(&string);
        } else {
            let string = format!(
                "External state change: {}|{}",
                location,
                String::from_utf8(contents).unwrap()
            );
            println!("{}\n", string);
            out.push_str(&string);
        }
    }
    fn get(&mut self, location: String) -> Result<Vec<u8>, String> {
        let (contents_vec_u8, success) = self
            .blockutil_interface
            .get_from_state(location.parse().unwrap());
        if success {
            Ok(contents_vec_u8)
        } else {
            Err("Could not get from onchain state".to_string())
        }
    }

    fn dump(&self) -> FxHashMap<String, Vec<u8>> {
        panic!("Not implemented.");
    }
}

pub struct StateManager<A: State, B: State, C: State> {
    pub cached_state: A,
    pub onchain_state: B,
    pub pending_state: C,
}

impl<A: State, B: State, C: State> StateManager<A, B, C> {
    pub fn new(
        blockutil_interface: NodeBlockUtilInterface,
        contract_hash: String,
        pending_state: PendingState,
    ) -> StateManager<CachedState, OnchainState, PendingState> {
        StateManager {
            cached_state: CachedState::new(),
            onchain_state: OnchainState::new(blockutil_interface, contract_hash),
            pending_state,
        }
    }
    pub fn write(&mut self, location: String, contents: Vec<u8>, out: &mut String) {
        self.cached_state.write(location, contents, out);
    }
    pub fn get_sync(&mut self, location: String) -> Result<Vec<u8>, String> {
        if let Ok(res) = self.pending_state.get(location.clone()) {
            Ok(res.to_vec())
        } else {
            self.get(location)
        }
    }
    pub fn get(&mut self, location: String) -> Result<Vec<u8>, String> {
        if let Ok(res) = self.cached_state.get(location.clone()) {
            Ok(res)
        } else {
            let onchain_res = self.onchain_state.get(location.clone());
            if onchain_res.is_err() {
                Err("Could not get from onchain state".to_string())
            } else {
                self.cached_state
                    .write(location, onchain_res.clone().unwrap(), &mut String::new());
                onchain_res
            }
        }
    }
    pub fn flush(&mut self, out: &mut String) {
        for (location, contents) in self.cached_state.dump().into_iter() {
            self.onchain_state.write(location, contents, out);
        }
    }
}
