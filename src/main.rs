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

    let traffic_data_dir_file = fs::read_to_string("/usr/local/etc/xray-stats/directory")
        .expect("/usr/local/etc/xray-stats/directory read error. Install https://github.com/moetelo/xray-stats first");

    let traffic_data_dir = traffic_data_dir_file.trim_end();

    let stats_parser = StatsParser::new(traffic_data_dir);

    let user_state = UserState::new(
        "/usr/local/etc/xray-stats-telegram/admins",
        "/usr/local/etc/xray-stats-telegram/users",
    ).expect("Failed to create UserState, check /usr/local/etc/xray-stats-telegram/admins and /usr/local/etc/xray-stats-telegram/users");

    BotInstance::new(Bot::from_env(), user_state.into(), stats_parser.into())
        .run()
        .await;
}
