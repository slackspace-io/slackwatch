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
    env::set_var("RUST_LOG", "info");
    env_logger::init();
     let settings = Settings::new().unwrap_or_else(|err| {
         log::error!("Failed to load settings: {}", err);
         panic!("Failed to load settings: {}", err);
     });
     log::info!("Starting up");
     log::info!("Loading configuration {:?}", settings);
    #[cfg(feature = "server")]
    use crate::database::client::create_table_if_not_exist;
    create_table_if_not_exist().unwrap();
    tokio::task::spawn(services::scheduler::run_scheduler(settings.clone()));
    let mut config = dioxus::fullstack::Config::new();

    #[cfg(feature = "server")]
    {
        config = config.addr(std::net::SocketAddr::from(([0, 0, 0, 0], 8080)));
    }

    let site = std::thread::spawn(|| LaunchBuilder::new().with_cfg(config).launch(App));
    log::info!("Started logger");
    println!("Started print");
    site.join().unwrap();
}

#[cfg(not(feature = "server"))]
fn main() {
    println!("Hello, world from non async main!");
    launch(App);

}

