// KiteConnect API logic (using reqwest for HTTP)
use serde::{Deserialize, Serialize};
use std::env;

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct KiteSession {
    pub api_key: String,
    pub api_secret: String,
    pub login_url: String,
}

pub fn get_kite_session() -> KiteSession {
    let api_key = env::var("KITE_API_KEY").unwrap_or_default();
    let api_secret = env::var("KITE_API_SECRET").unwrap_or_default();
    let login_url = format!("https://kite.zerodha.com/connect/login?v=3&api_key={}", api_key);
    KiteSession {
        api_key,
        api_secret,
        login_url,
    }
}
