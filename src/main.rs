use actix_files::Files;
use actix_web::{web, App, HttpServer, HttpResponse, Responder};
use dotenv::dotenv;
use std::env;

mod trading;
mod kite;
mod news;
mod models;


async fn index() -> impl Responder {
    HttpResponse::Ok().content_type("text/html").body(include_str!("../static/index.html"))
}

async fn api_backtest() -> impl Responder {
    // Placeholder: In real logic, run backtest and return results
    HttpResponse::Ok().body("Backtest run (placeholder)")
}

async fn api_news() -> impl Responder {
    match crate::news::fetch_news().await {
        Ok(news) => HttpResponse::Ok().json(news),
        Err(e) => HttpResponse::InternalServerError().body(format!("Error: {}", e)),
    }
}

async fn api_kite_login() -> impl Responder {
    let session = crate::kite::get_kite_session();
    HttpResponse::Ok().json(session)
}

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    dotenv().ok();
    let port = env::var("PORT").unwrap_or_else(|_| "8080".to_string());
    println!("Starting server on http://localhost:{}", port);

    HttpServer::new(|| {
        App::new()
            .route("/", web::get().to(index))
            .route("/api/backtest", web::get().to(api_backtest))
            .route("/api/news", web::get().to(api_news))
            .route("/api/kite-login", web::get().to(api_kite_login))
            .service(Files::new("/static", "static").show_files_listing())
    })
    .bind(("0.0.0.0", port.parse::<u16>().unwrap()))?
    .run()
    .await
}
