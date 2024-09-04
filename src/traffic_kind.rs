pub enum TrafficKind {
    Down,
    Up,
}

impl TrafficKind {
    pub fn as_str(&self) -> &'static str {
        match self {
            TrafficKind::Down => "down",
            TrafficKind::Up => "up",
        }
    }
}
