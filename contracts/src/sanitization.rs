pub fn sanitize_node_console_command(command: &String) -> bool {
    println!("Sanitizing command: {}", command);
    const FORBIDDEN_COMMANDS: [&str; 4] = [";", "send", "savestate", "keygen"];
    for forbidden_command in FORBIDDEN_COMMANDS.iter() {
        if command.contains(forbidden_command) {
            return false;
        }
    }
    true
}
