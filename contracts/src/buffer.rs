pub struct Buffer {
    pub contents: Vec<u8>,
}

impl Buffer {
    pub fn as_u64(&self) -> Result<u64, &'static str> {
        if self.contents.len() != 8 {
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
}