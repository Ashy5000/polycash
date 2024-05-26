// Copyright 2024, Asher Wrobel
/*
This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.
*/
#[derive(Debug, Clone)]
pub struct Buffer {
    pub contents: Vec<u8>,
}

impl Buffer {
    pub fn as_u64(&self) -> Result<u64, &'static str> {
        if self.contents.len() != 8 {
            println!("Invalid length. Actual length: {}", self.contents.len());
            return Err("Invalid length");
        }
        let mut result: u64 = 0;
        let mut shift_amount = 64 - 8;
        for piece in &self.contents {
            result += u64::from(*piece) << shift_amount;
            shift_amount -= 8;
        }
        Ok(result)
    }
    pub fn load_u64(&mut self, x: u64) {
        self.contents = Vec::new();
        let mut shift_amount = 64 - 8;
        while shift_amount >= 0 {
            let mut piece_u64 = x.clone();
            piece_u64 = piece_u64 >> shift_amount;
            piece_u64 = piece_u64 % 256;
            let piece = piece_u64 as u8;
            shift_amount -= 8;
            self.contents.push(piece);
        }
    }
}
