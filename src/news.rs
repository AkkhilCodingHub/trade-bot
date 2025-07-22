// News fetching logic (Yahoo Finance or similar via reqwest)
use crate::models::NewsItem;
use reqwest::header::{HeaderMap, HeaderValue};

pub async fn fetch_news() -> Result<Vec<NewsItem>, reqwest::Error> {
    let mut headers = HeaderMap::new();
    // You must set your own API key here or via env
    let rapidapi_key = std::env::var("RAPIDAPI_KEY").unwrap_or_default();
    headers.insert("X-RapidAPI-Key", HeaderValue::from_str(&rapidapi_key).unwrap());
    headers.insert("X-RapidAPI-Host", HeaderValue::from_static("yahoo-finance15.p.rapidapi.com"));

    let client = reqwest::Client::new();
    let res = client
        .get("https://yahoo-finance15.p.rapidapi.com/api/v2/markets/tickers?type=STOCKS&page=1")
        .headers(headers)
        .send()
        .await?;
    let data: serde_json::Value = res.json().await?;
    // Parse news items as needed
    let mut news = Vec::new();
    if let Some(items) = data["body"].as_array() {
        for item in items {
            news.push(NewsItem {
                title: item["name"].as_str().unwrap_or("").to_string(),
                url: "https://finance.yahoo.com".to_string(), // Placeholder
                summary: None,
            });
        }
    }
    Ok(news)
}
