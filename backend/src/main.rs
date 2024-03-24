use crate::config::Settings;
#[cfg(feature = "ssr")]
use crate::database::client::create_table_if_not_exist;
use crate::gitops::gitops::run_git_operations;
use k8s_openapi::merge_strategies::list::set;
use std::env;
use warp::host::exact;

mod config;
mod database;
mod gitops;
mod kubernetes;
mod models;
mod notifications;
mod repocheck;
mod services;
mod site;

#[tokio::main]
async fn main() {
    println!("Hello, world!");
    //Logging and env variables
    env::set_var("RUST_LOG", "info");
    env_logger::init();
    dotenv::dotenv().ok();
    // Load Configurations
    let settings = Settings::new().unwrap_or_else(|err| {
        log::error!("Failed to load settings: {}", err);
        panic!("Failed to load settings: {}", err);
    });
    log::info!("Starting up");
    log::info!("Loading configuration {:?}", settings);
    // Load Configurations
    //    log::info!("Configuration loaded: {:?}", settings);
    //create table if not exist
    #[cfg(feature = "ssr")]
    create_table_if_not_exist().unwrap();
    //tokio::task::spawn(services::scheduler::run_scheduler(settings.clone()));
    tokio::task::spawn_blocking(|| {
        println!("Site started");
        let _ = site::server::run_site();
    })
    .await
    .expect("Failed to run site")
}
