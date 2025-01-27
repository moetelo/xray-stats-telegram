use crate::{
    commands::{AdminCommand, UserCommand},
    query_date::QueryDate,
    stats_parser::StatsParser,
    user_state::UserState,
};
use crate::{date_util, stats::Stats};
use std::sync::Arc;
use teloxide::{prelude::*, utils::command::BotCommands};

pub async fn answer(
    bot: &Bot,
    user_state: Arc<UserState>,
    stats_parser: Arc<StatsParser>,
    xray_user: String,
    user_id: UserId,
    cmd: UserCommand,
) -> Result<(), teloxide::RequestError> {
    _ = match cmd {
        UserCommand::Help => {
            let help = if user_state.is_admin(user_id) {
                format!(
                    "{}\n\n{}",
                    UserCommand::descriptions(),
                    AdminCommand::descriptions()
                )
            } else {
                UserCommand::descriptions().to_string()
            };

            bot.send_message(user_id, help).await
        }
        UserCommand::Stats(string_date) => match date_util::date_or_today(string_date) {
            Ok(date) => {
                let stats = stats_parser.query_user_by_date(xray_user, &date);
                bot.send_message(user_id, stats.to_string()).await
            }
            Err(_) => handle_invalid_date(bot, user_id).await,
        },
    };

    Ok(())
}

pub async fn answer_admin(
    bot: &Bot,
    stats_parser: Arc<StatsParser>,
    user_id: UserId,
    command: AdminCommand,
) -> Result<(), teloxide::RequestError> {
    _ = match command {
        AdminCommand::All(string_date) => match date_util::date_or_today(string_date) {
            Err(_) => handle_invalid_date(bot, user_id).await,
            Ok(date) => handle_all(bot, stats_parser, user_id, &date).await,
        },
    };

    Ok(())
}

async fn handle_invalid_date(
    bot: &Bot,
    user_id: UserId,
) -> Result<Message, teloxide::RequestError> {
    bot.send_message(user_id, "Invalid date format. Use YYYY-MM-DD or YYYY-MM.")
        .await
}

async fn handle_all(
    bot: &Bot,
    stats_parser: Arc<StatsParser>,
    user_id: UserId,
    date: &QueryDate,
) -> Result<Message, teloxide::RequestError> {
    let all_stats: Vec<Stats> = stats_parser
        .get_all_users()
        .map(|user| stats_parser.query_user_by_date(user, date))
        .collect();

    let (empty_stats_users, non_empty_stats_users): (Vec<&Stats>, Vec<&Stats>) =
        all_stats.iter().partition(|stats| stats.is_empty());

    let mut message = non_empty_stats_users
        .iter()
        .map(|stats| format!("{}\n{}\n\n", stats.user, stats))
        .collect::<String>();

    let empty_stats_usernames: Vec<String> = empty_stats_users
        .into_iter()
        .map(|stats| stats.user.clone())
        .collect();

    if !empty_stats_usernames.is_empty() {
        message.push_str("No data for the following users:\n");
        message.push_str(&empty_stats_usernames.join(", "));
    }

    bot.send_message(user_id, message).await
}
