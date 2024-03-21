use crate::database::client::create_table_if_not_exist;
use std::env;
use crate::config::Settings;

mod oldc;
mod database;
mod kubernetes;
mod models;
mod config;
mod notifications;
mod repocheck;
mod services;
mod web;

#[tokio::main]
async fn main() {
    println!("Hello, world!");
    //Logging and env variables
    env::set_var("RUST_LOG", "info");
    env_logger::init();
    dotenv::dotenv().ok();
    let settings = Settings::new();
    log::info!("Starting up");
    log::info!("Loading configuration {:?}", settings);
    // Load Configurations
//    log::info!("Configuration loaded: {:?}", settings);
    //create table if not exist
    create_table_if_not_exist().unwrap();
    //find enabled
    //end_notification().await.unwrap();
    tokio::task::spawn_blocking(|| {
        let _ = web::exweb::site();
    })
    .await
    .expect("Failed to run site")
}
