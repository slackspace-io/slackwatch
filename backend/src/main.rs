use crate::config::Settings;
use crate::database::client::create_table_if_not_exist;
use std::env;
mod config;
mod database;
mod kubernetes;
mod models;
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
    //start scheduler
    //start site
    tokio::task::spawn(services::scheduler::run_scheduler(settings));
    tokio::task::spawn_blocking(|| {
        println!("Site started");
        let _ = web::exweb::site();
    })
    .await
    .expect("Failed to run site")
}
