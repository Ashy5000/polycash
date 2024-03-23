use std::fs;
use std::process::Command;
use regex::Regex;

#[derive(Debug)]
pub(crate) struct PublicKeyError {
    message: String,
}

pub(crate) fn get_public_key() -> Result<String, PublicKeyError> {
    let key_string_result = fs::read_to_string("../key.json");
    let key_string = match key_string_result {
        Ok(key_string) => key_string,
        Err(_) => {
            return Err(PublicKeyError {
                message: "Key not found".to_string(),
            })
        }
    };
    let re = Regex::new(r"Y(.):(?<y>.*),").unwrap();
    let Some(caps) = re.captures(&key_string) else {
        println!("no match!");
        return Ok("".to_string())
    };
    Ok(caps["y"].to_string())
}

pub(crate) fn generate_key() {
    let status = Command::new("env")
        .arg("-C")
        .arg("..")
        .arg("./builds/node/node_linux-amd64")
        .arg("-command")
        .arg("keygen")
        .status()
        .expect("Failed to run node executable");
    if status.success() {
        println!("Key generated successfully!");
    } else {
        println!("Failed to generate key");
    }
}