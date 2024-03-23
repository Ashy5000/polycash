use std::fs;
use regex::Regex;

pub(crate) fn get_public_key() -> String {
    let key_string = fs::read_to_string("../key.json").expect("Key file missing!");
    let re = Regex::new(r"Y(.):(?<y>.*),").unwrap();
    let Some(caps) = re.captures(&key_string) else {
        println!("no match!");
        return "".to_string()
    };
    caps["y"].to_string()
}