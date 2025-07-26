#![allow(non_snake_case, unused)]

use std::env;
use log::info;
use crate::config::Settings;

mod config;
mod database;
mod gitops;
mod kubernetes;
mod models;
mod notifications;
mod repocheck;
mod services;
mod api;

#[tokio::main]
async fn main() {
    println!("Hello, world!");
    env::set_var("RUST_LOG", "info");
    env_logger::init();

    let settings = Settings::new().unwrap_or_else(|err| {
        log::error!("Failed to load settings: {}", err);
        panic!("Failed to load settings: {}", err);
    });

    log::info!("Starting up");
    log::info!("Loading configuration {:?}", settings);

    use crate::database::client::create_table_if_not_exist;
    create_table_if_not_exist().unwrap();

    // Start the scheduler in a separate task
    tokio::task::spawn(services::scheduler::run_scheduler(settings.clone()));

    // Start the API server
    log::info!("Starting API server");
    api::start_api_server().await;
}
