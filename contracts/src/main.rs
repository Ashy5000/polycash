mod read_contract;
mod syntax_tree;
mod buffer;

fn main() {
    let contract_contents = read_contract::read_contract();
    println!("{}", contract_contents);
    let mut tree = syntax_tree::build_syntax_tree();
    tree.create(contract_contents);
    println!("{:?}", tree);
}
