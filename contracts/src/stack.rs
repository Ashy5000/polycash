use rustc_hash::FxHashMap;
use crate::buffer::Buffer;
use smartstring::alias::String;

// Copyright 2024, Asher Wrobel
/*
This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.
*/
pub struct StackFrame {
    pub buffers: FxHashMap<String, Buffer>
}

pub struct Stack {
    pub frames: Vec<StackFrame>
}

impl Stack {
    pub fn push(&mut self, buffers_param: &FxHashMap<String, Buffer>) {
        let frame = StackFrame {
            buffers: buffers_param.clone()
        };
        self.frames.push(frame);
    }
    pub fn pop(&mut self) -> FxHashMap<String, Buffer> {
        self.frames.pop().unwrap().buffers
    }
}