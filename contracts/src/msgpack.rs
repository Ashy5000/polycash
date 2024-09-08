use std::fs;
use std::fs::File;
use std::io::Read;
use rmpv::decode::read_value;
use rustc_hash::FxHashMap;

pub(crate) struct PendingState {
    pub(crate) data: FxHashMap<String, Vec<u8>>,
}
pub(crate) fn decode_pending_state() -> PendingState {
    let mut file = File::open("pending_state.msgpack").unwrap();
    let mut vec: Vec<u8> = Vec::new();
    file.read_to_end(&mut vec).unwrap();
    let mut reader = &*vec;
    let value = read_value(&mut reader).expect("Failed to decode pending state");
    assert!(value.is_map());
    let value_map = value.as_map().expect("Failed to decode pending state");
    let mut res = PendingState{data: FxHashMap::default()};
    for pair in value_map {
        let key = &pair.0;
        assert!(key.is_str());
        let key_string = String::from(key.as_str().unwrap());
        let value = &pair.1;
        assert!(value.is_str());
        let value_bytes = value.as_slice().unwrap();
        res.data.insert(key_string, Vec::from(value_bytes));
    }
    res
}