use contracts::blockutil::BlockUtilInterface;

#[derive(Clone)]
pub(crate) struct ZkBlockutilInterface {
    blockchain_len: u64
}

impl ZkBlockutilInterface {
    pub fn new(blockchain_len: u64) -> Self {
        Self {
            blockchain_len
        }
    }
}

impl BlockUtilInterface for ZkBlockutilInterface {
    fn read_contract(&self, location: u64) -> Result<smartstring::alias::String, smartstring::alias::String> {
        todo!()
    }

    fn get_from_state(&self, _property: smartstring::alias::String) -> (Vec<u8>, bool) {
        // Implemented through an alternative state manager, not through the blockutil interface.
        // Never used, and thus left unimplemented.
        unimplemented!()
    }

    fn get_blockchain_len(&self) -> u64 {
        self.blockchain_len
    }

    fn query_oracle(&self, query_type: u64, query_body: Vec<u8>) -> (Vec<u8>, bool) {
        todo!()
    }
}