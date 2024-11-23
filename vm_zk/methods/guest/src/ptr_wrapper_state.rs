use contracts::state::State;
use rustc_hash::FxHashMap;

pub(crate) struct PtrWrapperState<T> {
    pub(crate) ptr: *mut T
}


impl<T> PtrWrapperState<T> {
    pub(crate) fn new(ptr: *mut T) -> Self {
        Self {
            ptr
        }
    }
}

impl<T: State> State for PtrWrapperState<T> {
    fn write(&mut self, location: String, contents: Vec<u8>, out: &mut String) {
        unsafe { (*self.ptr).write(location, contents, out) }
    }

    fn get(&mut self, location: String) -> Result<Vec<u8>, String> {
        unsafe { (*self.ptr).get(location) }
    }

    fn dump(&self) -> FxHashMap<String, Vec<u8>> {
        unimplemented!()
    }
}