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