mod date;
mod stats;
mod telegram;
mod user_state;

use crate::telegram::BotInstance;
use std::fs;
use teloxide::prelude::Bot;

use crate::{stats::StatsParser, user_state::UserState};

#[tokio::main]
async fn main() {
    pretty_env_logger::init();
    log::info!("Starting command bot...");

    // TODO: make /usr/local/etc/xray-stats/directory a symlink
    let traffic_data_dir_file = fs::read_to_string("/usr/local/etc/xray-stats/directory")
        .expect("/usr/local/etc/xray-stats/directory should point to a valid traffic data directory. Install https://github.com/moetelo/xray-stats first");

    let traffic_data_dir = traffic_data_dir_file.trim_end();

    let stats_parser = StatsParser::new(traffic_data_dir);

    let user_state = UserState::new(
        "/usr/local/etc/xray-stats-telegram/admins",
        "/usr/local/etc/xray-stats-telegram/users",
    ).expect("/usr/local/etc/xray-stats-telegram/admins and /usr/local/etc/xray-stats-telegram/users should be present");

    BotInstance::new(Bot::from_env(), user_state.into(), stats_parser.into())
        .run()
        .await;
}
