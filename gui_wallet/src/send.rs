use std::process::Command;

pub(crate) fn send(amount: f64, address: String) {
    Command::new("env")
        .arg("-C")
        .arg("..")
        .arg("./builds/node/node_linux-amd64")
        .arg("-command")
        .arg(format!("send {address} {amount}"))
        .output()
        .expect("Failed to run node executable");
    println!("Transaction broadcasted successfully")
}