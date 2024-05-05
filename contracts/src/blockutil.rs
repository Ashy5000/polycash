pub struct BlockUtilInterface {
    node_executable_path: String,
}

impl BlockUtilInterface {
    pub fn new() -> Self {
        // Read the executable path from the executable path file
        let node_executable_path = std::fs::read_to_string("node_executable_path.txt")
            .expect("Could not read node_executable_path.txt");
        // Remove the newline character
        // This is necessary because the path is read from a file
        let node_executable_path = node_executable_path.trim().to_string();
        Self {
            node_executable_path,
        }
    }
    pub fn get_nth_block_property(&self, n: i64, property: String) -> String {
        // Run the node script to get a property of the nth block
        let output = std::process::Command::new(self.node_executable_path.clone())
            .arg("--command")
            .arg(format!("sync;getNthBlock {} {}", n, property))
            .output();
        let output = output.expect("Failed to execute node script");
        let output = String::from_utf8(output.stdout).expect("Failed to convert output to string");
        // Remove the newline character
        let output = output.trim().to_string();
        // Split the output by the newline character
        let output = output.split("\n").collect::<Vec<&str>>();
        // Get the last element of the output
        // This is neccessary because the previous lines contain logs
        let output = output[output.len() - 1].to_string();
        output
    }
}
