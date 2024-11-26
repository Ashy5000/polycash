mod lazy_vector;
mod prove;
mod verify;
mod request_handler;
mod socket;

use crate::socket::Socket;

fn main() {
    let mut socket = Socket::new();
    socket.run().unwrap();
}
