// Copyright 2024, Asher Wrobel
/*
This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.
*/
pub fn sanitize_node_console_command(command: &String) -> bool {
    println!("Sanitizing command: {}", command);
    const FORBIDDEN_COMMANDS: [&str; 3] = [";send", ";savestate", ";keygen"];
    for forbidden_command in FORBIDDEN_COMMANDS.iter() {
        if command.contains(forbidden_command) {
            return false;
        }
    }
    true
}
