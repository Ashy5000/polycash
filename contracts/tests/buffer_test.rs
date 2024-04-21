extern crate contracts;

#[cfg(test)]
mod buffer_test {
    #[test]
    fn test_buffer() {
        let buffer = contracts::buffer::Buffer {
            contents: vec![0, 1, 2, 3],
        };
        assert_eq!(buffer.contents[0], 0);
        assert_eq!(buffer.contents[1], 1);
        assert_eq!(buffer.contents[2], 2);
        assert_eq!(buffer.contents[3], 3);
    }
    #[test]
    fn test_buffer_as_u64() {
        let buffer = contracts::buffer::Buffer {
            contents: vec![0, 0, 0, 0, 0, 0, 0, 123],
        };
        let buffer_u64 = buffer.as_u64();
        assert_eq!(buffer_u64, Ok(123));

        let buffer = contracts::buffer::Buffer {
            contents: vec![0, 0, 0, 0, 0, 0, 1, 0],
        };
        let buffer_u64 = buffer.as_u64();
        assert_eq!(buffer_u64, Ok(256));
    }
    #[test]
    fn test_buffer_load_u64() {
        let buffer_u64: u64 = 256;
        let mut buffer = contracts::buffer::Buffer {
            contents: Vec::new(),
        };
        buffer.load_u64(buffer_u64);
        assert_eq!(buffer.contents.len(), 8);
        assert_eq!(buffer.as_u64(), Ok(256));
    }
}
