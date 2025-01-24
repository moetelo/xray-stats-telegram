use crate::stats::Stats;
use crate::traffic_kind::TrafficKind;
use std::fs::{self};
use std::path::{Path, PathBuf};

use chrono::NaiveDate;

#[derive(Debug, Clone)]
pub struct StatsParser {
    traffic_data_directory: PathBuf,
}

impl StatsParser {
    pub fn new<P: AsRef<Path>>(traffic_data_directory: P) -> Self {
        Self {
            traffic_data_directory: traffic_data_directory.as_ref().to_path_buf(),
        }
    }

    pub fn query_user_by_date(&self, user: String, date: NaiveDate) -> Stats {
        let down = self.get_user_traffic(&user, TrafficKind::Down, date);
        let up = self.get_user_traffic(&user, TrafficKind::Up, date);

        Stats {
            user,
            down,
            up,
        }
    }

    fn get_user_traffic(&self, user: &String, traffic_kind: TrafficKind, date: NaiveDate) -> u64 {
        let iso_date = date.format("%Y-%m-%d").to_string();

        let path = self
            .traffic_data_directory
            .join(user)
            .join(traffic_kind.as_str())
            .join(iso_date);

        fs::read_to_string(path)
            .unwrap_or_default()
            .lines()
            .filter_map(|line| line.parse::<u64>().ok())
            .sum()
    }

    pub fn get_all_users(&self) -> impl Iterator<Item = String> {
        self.traffic_data_directory
            .read_dir()
            .expect("Failed to read traffic data directory")
            .filter_map(|entry| entry.ok())
            .filter(|entry| {
                entry
                    .file_type()
                    .map(|file_type| file_type.is_dir())
                    .unwrap_or(false)
            })
            .map(|entry| entry.file_name().into_string().unwrap())
    }
}
