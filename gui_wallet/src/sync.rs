use regex::Regex;
use std::fs;
use std::process::Command;
use std::str::FromStr;

pub(crate) fn sync() -> f64 {
    let key_string = fs::read_to_string("../key.json").expect("Key file missing!");
    let re = Regex::new(r"Y(.):(?<y>.*)}").unwrap();
    let Some(caps) = re.captures(&key_string) else {
        println!("no match!");
        return -1.0;
    };
    let output = Command::new("env")
        .arg("-C")
        .arg("..")
        .arg("./builds/node/node_linux-amd64")
        .arg("-command")
        .arg("sync;balance ".to_owned() + &caps["y"])
        .output()
        .expect("Failed to run node executable");
    let output_string =
        String::from_utf8(output.stdout).expect("Failed to convert output to string");
    let re = Regex::new(r": (?<balance>[0-9]\.[0-9]*)").unwrap();
    let Some(caps) = re.captures(&output_string) else {
        println!("no match!");
        return -1.0;
    };
    let balance = f64::from_str(&caps["balance"]).expect("Invalid balance format");
    balance
}
