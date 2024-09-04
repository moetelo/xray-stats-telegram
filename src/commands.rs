use teloxide::utils::command::BotCommands;

#[derive(BotCommands, Clone)]
#[command(rename_rule = "lowercase", description = "Supported commands:")]
pub enum UserCommand {
    #[command(description = "display this text.")]
    Help,
    #[command(description = "get stats for date. `/stats 2024-09-01`")]
    Stats(String),
}

#[derive(BotCommands, Clone)]
#[command(rename_rule = "lowercase", description = "Admin commands:")]
pub enum AdminCommand {
    #[command(description = "get stats for all users. `/all 2024-09-01`")]
    All(String),
}
