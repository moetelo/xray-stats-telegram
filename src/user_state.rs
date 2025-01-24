use std::collections::HashMap;
use std::collections::HashSet;
use std::fs;
use std::path::Path;

use teloxide::types::UserId;

#[derive(Debug, Clone)]
pub struct UserState {
    admins: HashSet<UserId>,
    tg_id_to_xray_email: HashMap<UserId, String>,
}

impl UserState {
    pub fn new(admins_path: &str, users_path: &str) -> Result<Self, Box<dyn std::error::Error>> {
        let admins_file = fs::read_to_string(Path::new(admins_path))?;
        let admins = HashSet::from_iter(admins_file.lines().map(|line| {
            let admin_id: u64 = line.parse().expect("Expected id, found non-integer");
            UserId(admin_id)
        }));

        let users_file = fs::read_to_string(Path::new(users_path))?;
        let tg_id_to_xray_email =
            HashMap::from_iter(users_file.lines().map(Self::parse_users_line));

        Ok(Self {
            admins,
            tg_id_to_xray_email,
        })
    }

    fn parse_users_line(line: &str) -> (UserId, String) {
        let mut parts_iter = line.split(':');
        let tg_id: u64 = parts_iter
            .next()
            .expect("expected tg_id:xray_email")
            .parse()
            .expect("expected integer");

        let xray_email = parts_iter
            .next()
            .expect("expected tg_id:xray_email")
            .to_string();

        let _ = parts_iter.next().is_none_or(|_| {
            panic!("expected tg_id:xray_email, found more than 1 colon (:)");
        });

        (UserId(tg_id), xray_email)
    }

    pub fn get_xray_email(&self, chat_id: UserId) -> Option<&String> {
        self.tg_id_to_xray_email.get(&chat_id)
    }

    pub fn is_admin(&self, id: UserId) -> bool {
        self.admins.contains(&id)
    }
}
