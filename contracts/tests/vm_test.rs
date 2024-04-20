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
}