// Common data models
use serde::{Serialize, Deserialize};

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct TradeDecision {
    pub buy: bool,
    pub sell: bool,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct NewsItem {
    pub title: String,
    pub url: String,
    pub summary: Option<String>,
}
