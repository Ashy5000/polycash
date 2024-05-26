// Copyright 2024, Asher Wrobel
/*
This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.
*/
#[derive(Debug, Clone)]
pub struct Line {
    pub command: String,
    pub args: Vec<String>,
}

pub(crate) fn build_line() -> Line {
    Line {
        command: "".to_owned(),
        args: Vec::new(),
    }
}

#[derive(Debug)]
pub struct SyntaxTree {
    pub lines: Vec<Line>
}

impl SyntaxTree {
    pub fn create(&mut self, contract_contents: String) {
        let asm_lines_iter = contract_contents.split("\n");
        for asm_line in asm_lines_iter {
            let parts_iter = asm_line.split(" ");
            let mut parts = Vec::new();
            for mut part in parts_iter {
                if part == ";".to_owned() {
                    break
                }
                let part_chars = &part.chars();
                let first_two_chars: String = part_chars.to_owned().into_iter().take(2).collect();
                let part_string;
                if first_two_chars == "0x".to_owned() {
                    part_string = part_chars.to_owned().into_iter().skip(2).take(part.len() - 2).collect::<String>().to_owned();
                    part = part_string.as_str();
                }
                parts.push(part.to_owned())
            }
            if parts.len() == 0 {
                continue
            }
            let mut line = build_line();
            line.command = parts[0].to_owned();
            let args = parts.split_off(1);
            line.args = args;
            self.lines.push(line);
        }
    }
}

pub fn build_syntax_tree() -> SyntaxTree {
    SyntaxTree {
        lines: Vec::new(),
    }
}