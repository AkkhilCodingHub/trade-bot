// Trading logic module (momentum/value strategy)
use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct StockData {
    pub date: String,
    pub open: f64,
    pub high: f64,
    pub low: f64,
    pub close: f64,
    pub volume: f64,
    pub pe: Option<f64>,
    pub pe_avg: Option<f64>,
}

pub fn calculate_momentum(data: &[StockData], period: usize) -> Option<f64> {
    if data.len() < period + 1 {
        return None;
    }
    let latest = data.last()?.close;
    let previous = data[data.len() - period - 1].close;
    Some(latest - previous)
}

pub fn should_buy(momentum: f64, pe: f64, pe_avg: f64) -> bool {
    momentum > 0.0 && pe < pe_avg
}

pub fn should_sell(momentum: f64, pe: f64, pe_avg: f64) -> bool {
    momentum < 0.0 || pe > pe_avg
}
