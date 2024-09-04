mod user_state;
mod stats;
mod stats_parser;
mod date_util;
mod commands;
mod traffic_kind;
mod handlers;

use teloxide::{prelude::*, update_listeners};
use crate::{user_state::UserState, stats_parser::StatsParser, commands::{AdminCommand, UserCommand}};
use std::{fs, sync::Arc};

#[tokio::main]
async fn main() {
    pretty_env_logger::init();
    log::info!("Starting command bot...");

    let traffic_data_dir_file = fs::read_to_string("/usr/local/etc/xray-stats/directory")
        .expect("/usr/local/etc/xray-stats/directory read error. Install https://github.com/moetelo/xray-stats first");

    let traffic_data_dir = traffic_data_dir_file.trim_end();

    let stats_parser = StatsParser::new(traffic_data_dir);

    let user_state = Arc::new(UserState::new(
        "/usr/local/etc/xray-stats-telegram/admins",
        "/usr/local/etc/xray-stats-telegram/users",
    ));

    let user_commands_endpoint = dptree
        ::filter_map(|user_id: UserId, user_state: Arc<UserState>| {
            user_state.get_xray_email(user_id).cloned()
        })
        .filter_command::<UserCommand>()
        .endpoint(handlers::answer);

    let admin_commands_endpoint = dptree
        ::filter(|user_id: UserId, user_state: Arc<UserState>| {
            user_state.is_admin(user_id)
        })
        .filter_command::<AdminCommand>()
        .endpoint(handlers::answer_admin);

    let handler = Update::filter_message()
        .filter_map(|msg: Message| msg.chat.id.as_user())
        .branch(user_commands_endpoint)
        .branch(admin_commands_endpoint);

    let ignore_update = |_upd| Box::pin(async {});
    let bot = Bot::from_env();
    Dispatcher::builder(bot.clone(), handler)
        .dependencies(dptree::deps![
            user_state,
            stats_parser
        ])
        .default_handler(ignore_update)
        .error_handler(LoggingErrorHandler::with_custom_text(
            "An error has occurred in the dispatcher",
        ))
        .enable_ctrlc_handler()
        .build()
        .dispatch_with_listener(
            update_listeners::polling_default(bot.clone()).await,
            LoggingErrorHandler::with_custom_text("An error from the update listener"),
        )
        .await
}
