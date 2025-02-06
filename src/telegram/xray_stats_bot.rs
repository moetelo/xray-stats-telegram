use std::sync::Arc;
use teloxide::{prelude::*, update_listeners};

use crate::{
    stats::stats_parser::StatsParser,
    telegram::commands::{AdminCommand, UserCommand},
    telegram::handlers,
    user_state::UserState,
};

pub struct BotInstance {
    bot: Bot,
    user_state: Arc<UserState>,
    stats_parser: Arc<StatsParser>,
}

impl BotInstance {
    pub fn new(bot: Bot, user_state: Arc<UserState>, stats_parser: Arc<StatsParser>) -> Self {
        Self {
            bot,
            user_state,
            stats_parser,
        }
    }

    pub async fn run(&self) {
        let user_commands_endpoint =
            dptree::filter_map(|user_id: UserId, user_state: Arc<UserState>| {
                user_state.get_xray_email(user_id).cloned()
            })
            .filter_command::<UserCommand>()
            .endpoint(handlers::answer);

        let admin_commands_endpoint =
            dptree::filter(|user_id: UserId, user_state: Arc<UserState>| {
                user_state.is_admin(user_id)
            })
            .filter_command::<AdminCommand>()
            .endpoint(handlers::answer_admin);

        let message_handler = Update::filter_message()
            .branch(user_commands_endpoint)
            .branch(admin_commands_endpoint);

        let callback_query_handler =
            Update::filter_callback_query().endpoint(handlers::answer_callback_query);

        let combined_handler = dptree::entry()
            .filter_map(|msg: Message| msg.chat.id.as_user())
            .branch(message_handler)
            .branch(callback_query_handler);

        let ignore_update = |_upd| Box::pin(async {});
        Dispatcher::builder(self.bot.clone(), combined_handler)
            .dependencies(dptree::deps![
                self.user_state.clone(),
                self.stats_parser.clone()
            ])
            .default_handler(ignore_update)
            .error_handler(LoggingErrorHandler::with_custom_text(
                "An error has occurred in the dispatcher",
            ))
            .enable_ctrlc_handler()
            .build()
            .dispatch_with_listener(
                update_listeners::polling_default(self.bot.clone()).await,
                LoggingErrorHandler::with_custom_text("An error from the update listener"),
            )
            .await;
    }
}
