mod read_contract;

fn main() {
    let contract_contents = read_contract::read_contract();
    println!("Contract contents:\n{}", contract_contents);
}
