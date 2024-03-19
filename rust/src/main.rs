use crate::database::client::create_table_if_not_exist;
use crate::kubernetes::client::find_enabled_workloads;
use kube::Client;
use kubernetes::client;
use std::env;
use std::sync::Arc;
use tokio::sync::Mutex;
use warp::Filter;

mod config;
mod database;
mod kubernetes;
mod models;
mod web;



#[tokio::main]
async fn main() {
    println!("Hello, world!");

    //Logging and env variables
    env::set_var("RUST_LOG", "debug");
    env_logger::init();
    dotenv::dotenv().ok();

    // Load Configuration
    let config = match config::AppConfig::new() {
        Ok(cfg) => cfg,
        Err(e) => {
            log::error!("Failed to load configuration: {}", e);
            std::process::exit(1);
        }
    };
    println!("Configuration loaded: {:?}", config);
    log::info!("Configuration loaded: {:?}", config);
    //create table if not exist
    create_table_if_not_exist().unwrap();
    //find enabled
    let workloads = find_enabled_workloads().await.unwrap();
    //print count of workloads
    println!("Workloads found: {}", workloads.len());
    //print first 5 workloads
    // Initialize shared state
    for workload in workloads.iter() {
        println!("Workload: {:?}", workload);
        //insert into db
        database::client::insert_workload(workload).unwrap();
    }
    //get workload from db
    if let workload = database::client::return_workload("frigate".to_string(), "frigate".to_string()) {
        println!("Workload from db: {:?}", workload);

    } else {
        println!("Workload not found");
    }
    //return all
    let workloads = database::client::return_all_workloads().unwrap();
    println!("Workloads from db: {:?}", workloads);
    //run site
    tokio::task::spawn_blocking(|| {
        crate::web::exweb::site();
    }).await.expect("Failed to run site")
}
