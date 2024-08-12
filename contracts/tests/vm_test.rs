extern crate contracts;

#[cfg(test)]
mod vm_test {
    use rustc_hash::FxHashMap;
    use smartstring::alias::String;
    use contracts::buffer::Buffer;

    #[test]
    fn test_vm_access_buffer_contents() {
        let mut buffers: FxHashMap<String, Buffer> = FxHashMap::default();
        buffers.insert(
            "00000000".parse().unwrap(),
            Buffer { contents: Vec::new() }
        );
        let contents = contracts::vm::vm_access_buffer_contents(&mut buffers, "12341234".parse().unwrap(), "00000000".parse().unwrap());
        assert_eq!(contents, vec![]);
        if let Some(x) = buffers.get("00000000") {
            assert_eq!(x.contents, vec![1]);
        } else {
            panic!("No buffer found at location `00000000`");
        }
        buffers.insert(
            "12341234".parse().unwrap(),
            Buffer { contents: vec![0, 1, 2, 3] }
        );
        let contents = contracts::vm::vm_access_buffer_contents(&mut buffers, "12341234".parse().unwrap(), "00000000".parse().unwrap());
        assert_eq!(contents, vec![0, 1, 2, 3]);
    }

    #[test]
    fn test_vm_access_buffer() {
        let mut buffers: FxHashMap<String, Buffer> = FxHashMap::default();
        buffers.insert(
            "00000000".parse().unwrap(),
            Buffer { contents: Vec::new() }
        );
        let buffer = contracts::vm::vm_access_buffer(&mut buffers, "12341234".parse().unwrap(), "00000000".parse().unwrap());
        assert_eq!(buffer, contracts::vm::VM_NIL as *mut Buffer);
        if let Some(x) = buffers.get("00000000") {
            assert_eq!(x.contents, vec![1]);
        } else {
            panic!("No buffer found at location `00000000`");
        }
        buffers.insert(
            "12341234".parse().unwrap(),
            Buffer { contents: vec![0, 1, 2, 3] }
        );
        let buffer = contracts::vm::vm_access_buffer(&mut buffers, "12341234".parse().unwrap(), "00000000".parse().unwrap());
        unsafe {
            assert_eq!((*buffer).contents, vec![0, 1, 2, 3]);
        }
    }
    
    #[test]
    fn test_vm_check_buffer_initialization() {
        let mut buffers: FxHashMap<String, Buffer> = FxHashMap::default();
        let initialization = contracts::vm::vm_check_buffer_initialization(&mut buffers, "00000000".parse().unwrap());
        assert!(!initialization);
        buffers.insert(
            "12341234".parse().unwrap(),
            Buffer { contents: Vec::new() }
        );
        let initialization = contracts::vm::vm_check_buffer_initialization(&mut buffers, "12341234".parse().unwrap());
        assert!(initialization);
    }

    #[test]
    fn test_vm_throw_global_error() {
        let mut buffers: FxHashMap<String, Buffer> = FxHashMap::default();
        buffers.insert(
            "00000000".parse().unwrap(),
            Buffer { contents: Vec::new() }
        );
        contracts::vm::vm_throw_global_error(&mut buffers);
        let mut found_buffer = false;
        if let Some(x) = buffers.get("00000000") {
            found_buffer = true;
            assert_eq!(x.contents.len(), 1);
            assert_eq!(x.contents[0], 1);
        }
        assert!(found_buffer);
    }

    #[test]
    fn test_vm_throw_local_error() {
        let mut buffers: FxHashMap<String, Buffer> = FxHashMap::default();
        buffers.insert(
            "12341234".parse().unwrap(),
            Buffer { contents: Vec::new() }
        );
        contracts::vm::vm_throw_local_error(&mut buffers, "12341234".parse().unwrap());
        let mut found_buffer = true;
        if let Some(x) = buffers.get("12341234") {
            found_buffer = true;
            assert_eq!(x.contents.len(), 1);
            assert_eq!(x.contents[0], 1);
        }
        assert!(found_buffer);
    }

    #[test]
    fn test_vm_exit() {
        let mut buffers: FxHashMap<String, Buffer> = FxHashMap::default();
        buffers.insert(
            "00000000".parse().unwrap(),
            Buffer { contents: Vec::new() }
        );
        let contents = "\
            Exit 1234
        ";
        let mut tree = contracts::syntax_tree::build_syntax_tree();
        tree.create(contents.parse().unwrap());
        let blockutil_interface = contracts::blockutil::BlockUtilInterface::new();
        let contract_hash = "00000000000000000000000000000000";
        let (exit_code, gas_used) = contracts::vm::run_vm(tree, &mut buffers, blockutil_interface, contract_hash.parse().unwrap());
        assert_eq!(exit_code, 1234);
        assert_eq!(gas_used, 1.0);
    }
}
