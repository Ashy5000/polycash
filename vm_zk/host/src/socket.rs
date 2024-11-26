use std::io;
use std::io::{Read, Write};
use std::os::unix::net::UnixStream;
use crate::request_handler::handle_request;

pub(crate) struct Socket {
    socket: UnixStream
}

impl Socket {
    pub(crate) fn new() -> Self {
        let socket = UnixStream::connect("/tmp/vm.sock").unwrap();
        Self { socket }
    }

    fn read_exact(&mut self, n: usize) -> io::Result<Vec<u8>> {
        let mut buffer = vec![0u8; n];
        self.socket.read_exact(&mut buffer)?;
        Ok(buffer)
    }

    fn read_message(&mut self) -> io::Result<Vec<u8>> {
        let len_bytes = self.read_exact(4)?;
        let length = u32::from_le_bytes(len_bytes.try_into().unwrap()) as usize;
        println!("{}", length);
        self.read_exact(length)
    }

    pub(crate) fn write_message(&mut self, data: &[u8]) -> io::Result<()> {
        let length = data.len() as u32;
        self.socket.write_all(&length.to_le_bytes())?;
        self.socket.write_all(data)?;
        Ok(())
    }

    pub(crate) fn run(&mut self) -> io::Result<()> {
        loop {
            match self.read_message() {
                Ok(message) => {
                    handle_request(message, self);
                }
                Err(_) => {
                    continue
                }
            }
        }
    }
}