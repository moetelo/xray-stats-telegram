use crate::{stats::Stats, date_util};
use teloxide::{prelude::*, utils::command::BotCommands};
use crate::{user_state::UserState, stats_parser::StatsParser, commands::{AdminCommand, UserCommand}};
use std::sync::Arc;

pub async fn answer(
    user_state: Arc<UserState>,
    stats_parser: StatsParser,
    bot: Bot,
    xray_user: String,
    user_id: UserId,
    cmd: UserCommand,
) -> Result<(), teloxide::RequestError> {
    match cmd {
        UserCommand::Help => {
            let mut help = UserCommand::descriptions().to_string();
            if user_state.is_admin(user_id) {
                help.push_str("\n\n");
                help.push_str(AdminCommand::descriptions().to_string().as_str());
            }

            bot.send_message(user_id, help).await?;
            return Ok(());
        },
        UserCommand::Stats(string_date) => {
            match date_util::date_or_today(string_date) {
                Ok(date) => {
                    let stats = stats_parser.query_user_by_date(xray_user.as_str(), date);
                    bot.send_message(user_id, stats.to_string()).await?;
                    return Ok(());
                }
                Err(_) => {
                    bot.send_message(user_id, "Invalid date format. Use YYYY-MM-DD.").await?;
                    return Ok(());
                }
            }
        }
    };
}

pub async fn answer_admin(
    stats_parser: StatsParser,
    bot: Bot,
    user_id: UserId,
    cmd: AdminCommand,
) -> Result<(), teloxide::RequestError> {
    match cmd {
        AdminCommand::All(string_date) => {
            match date_util::date_or_today(string_date) {
                Ok(date) => {
                    let all_stats = stats_parser.get_all_users().iter()
                        .map(|user| stats_parser.query_user_by_date(user, date))
                        .collect::<Vec<Stats>>();

                    let empty_stats_users = all_stats.iter()
                        .filter(|stats| stats.is_empty())
                        .map(|stats| stats.user.clone())
                        .collect::<Vec<String>>();

                    let non_empty_stats_users = all_stats.iter()
                        .filter(|stats| !stats.is_empty())
                        .collect::<Vec<&Stats>>();

                    let mut message = String::new();
                    for stats in non_empty_stats_users {
                        message.push_str(&stats.user);
                        message.push_str("\n");
                        message.push_str(&stats.to_string());
                        message.push_str("\n\n");
                    }

                    if !empty_stats_users.is_empty() {
                        message.push_str("No data for the following users:\n");
                        message.push_str(&empty_stats_users.join(", "));
                    }

                    bot.send_message(user_id, message).await?
                },
                Err(_) => {
                    bot.send_message(user_id, "Invalid date format. Use YYYY-MM-DD.").await?;
                    return Ok(());
                }
            }
        }
    };

    Ok(())
}
