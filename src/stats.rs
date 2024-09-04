use std::fmt;

#[derive(Debug)]
pub struct Stats {
    pub user: String,
    pub down: u64,
    pub up: u64,
}

impl Stats {
    pub fn is_empty(&self) -> bool {
        self.down == 0 && self.up == 0
    }
}

impl fmt::Display for Stats {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        write!(f, "↓ {} (mb) ↑ {} (mb)", self.down / 1024 / 1024, self.up / 1024 / 1024)
    }
}
