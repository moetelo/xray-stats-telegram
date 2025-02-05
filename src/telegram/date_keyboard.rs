use teloxide::types::{InlineKeyboardButton, InlineKeyboardMarkup};

use crate::stats::QueryDate;

pub fn make(date: &QueryDate) -> InlineKeyboardMarkup {
    InlineKeyboardMarkup::new([[
        InlineKeyboardButton::callback("⬅️", date.prev().to_string()),
        InlineKeyboardButton::callback("🔄", date.to_string()),
        InlineKeyboardButton::callback("➡️", date.next().to_string()),
    ]])
}
