#![allow(non_snake_case, unused)]

use std::env;
//#[cfg(feature = "server")]
//use crate::config::Settings;
//#[cfg(feature = "server")]
//use crate::database::client::create_table_if_not_exist;
//#[cfg(feature = "server")]
//use crate::gitops::gitops::run_git_operations;
//use dioxus::core_macro::rsx;
//use dioxus::dioxus_core::Element;
//use dioxus::hooks::use_signal;
//use dioxus::prelude::launch;
use crate::site::app::App;
use dioxus::prelude::*;
use log::info;
use crate::config::Settings;

//#[cfg(feature = "server")]
//use k8s_openapi::merge_strategies::list::set;
//use std::env;
//
mod config;
mod database;
mod gitops;
mod kubernetes;
mod models;
mod notifications;
mod repocheck;
mod services;
mod site;


#[cfg(feature = "server")]
#[tokio::main]
async fn  main() {
    println!("Hello, world!");
    //Logging and env variables
    env::set_var("RUST_LOG", "info");
    env_logger::init();
    //  dotenv::dotenv().ok();
    // Load Configurations
     let settings = Settings::new().unwrap_or_else(|err| {
         log::error!("Failed to load settings: {}", err);
         panic!("Failed to load settings: {}", err);
     });
    // log::info!("Starting up");
    // log::info!("Loading configuration {:?}", settings);
    // Load Configurations
    //    log::info!("Configuration loaded: {:?}", settings);
    //create table if not exist
    //create_table_if_not_exist().unwrap();
    //find enabled
    //end_notification().await.unwrap();
    //start scheduler
    //start site
    //run_git_operations(settings.clone()).unwrap_or_else(|err| {
    //    log::error!("Failed to run git operations: {}", err);
    //    panic!("Failed to run git operations: {}", err);
    //});
    //launch(app);
    //working tokio stuff
    tokio::task::spawn(services::scheduler::run_scheduler(settings.clone()));

    let site = std::thread::spawn(|| launch(App));
    log::info!("Started logger");
    println!("Started print");
    // join
//    tokio::task::spawn(launch(App));
    //tokio::task::spawn_blocking(|| launch(App));
    //let _ = web::exweb::site();
    site.join().unwrap();
}

#[cfg(not(feature = "server"))]
fn main() {
    println!("Hello, world from non async main!");
    launch(App);

}

