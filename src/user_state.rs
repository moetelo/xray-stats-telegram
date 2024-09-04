use std::collections::HashSet;
use std::collections::HashMap;
use std::fs;
use std::path::Path;

use teloxide::types::UserId;

#[derive(Debug, Clone)]
pub struct UserState {
    admins: HashSet<UserId>,
    tg_id_to_xray_email: HashMap<UserId, String>,
}

impl UserState {
    pub fn new(admins_path: &str, users_path: &str) -> Self {
        let admins_file = fs::read_to_string(Path::new(admins_path)).unwrap();

        let mut admins = HashSet::new();
        for line in admins_file.lines() {
            let admin_id: u64 = line.parse().unwrap_or_default();
            admins.insert(UserId(admin_id));
        }

        let users = fs::read_to_string(Path::new(users_path)).unwrap();
        let mut tg_id_to_xray_email = HashMap::new();
        for line in users.lines() {
            let parts: Vec<&str> = line.split(':').collect();
            if parts.len() == 2 {
                let tg_id: u64 = parts[0].parse().unwrap();
                let xray_email = parts[1].to_string();
                tg_id_to_xray_email.insert(UserId(tg_id), xray_email);
            }
        }

        Self {
            admins,
            tg_id_to_xray_email,
        }
    }

    pub fn get_xray_email(&self, chat_id: UserId) -> Option<&String> {
        self.tg_id_to_xray_email.get(&chat_id)
    }

    pub fn is_admin(&self, id: UserId) -> bool {
        self.admins.contains(&id)
    }
}
