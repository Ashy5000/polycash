extern crate contracts;

#[cfg(test)]
mod vm_test {
    use std::collections::HashMap;

    use contracts::buffer::Buffer;

    #[test]
    fn test_vm_check_buffer_initialization() {
        let mut buffers: HashMap<String, Buffer> = HashMap::new();
        let initialization = contracts::vm::vm_check_buffer_initialization(&mut buffers, "00000000".to_owned());
        assert!(!initialization);
        buffers.insert(
            "12341234".to_owned(),
            Buffer { contents: Vec::new() }
        );
        let initialization = contracts::vm::vm_check_buffer_initialization(&mut buffers, "12341234".to_owned());
        assert!(initialization);
    }

    #[test]
    fn test_vm_throw_global_error() {
        let mut buffers: HashMap<String, Buffer> = HashMap::new();
        buffers.insert(
            "00000000".to_owned(),
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
        let mut buffers: HashMap<String, Buffer> = HashMap::new();
        buffers.insert(
            "12341234".to_owned(),
            Buffer { contents: Vec::new() }
        );
        contracts::vm::vm_throw_local_error(&mut buffers, "12341234".to_owned())
    }
}