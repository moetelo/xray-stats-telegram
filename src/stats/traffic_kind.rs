pub enum TrafficKind {
    Down,
    Up,
}

impl std::fmt::Display for TrafficKind {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        let s = match self {
            TrafficKind::Down => "down",
            TrafficKind::Up => "up",
        };

        write!(f, "{}", s)
    }
}
