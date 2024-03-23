use crate::database;
use crate::database::client::get_latest_scan_id;
use crate::gitops::gitops::run_git_operations;
use crate::models::models;
use crate::models::models::{ApiResponse, Workload};
use crate::services;
use crate::services::workloads;
use actix_web::{get, web, App, HttpRequest, HttpResponse, HttpServer, Responder};
use std::collections::HashMap;

#[get("/")]
async fn index() -> impl Responder {
    let _ = workloads::test_call().await;
    let scan = get_latest_scan_id();
    log::info!("Scan id: {:?}", scan);
    log::info!("Hello world");
    HttpResponse::Ok().body("Hello world!")
}

#[get("/api/workloads/all")]
async fn fetch_all_workloads() -> impl Responder {
    //get workload from db
    if let Ok(workload) = database::client::return_all_workloads() {
        HttpResponse::Ok().json(workload)
    } else if let Err(e) = database::client::return_all_workloads() {
        HttpResponse::Ok().json(e.to_string())
    } else {
        HttpResponse::Ok().body("Workload not found")
    }
}

#[get("/api/workloads/refresh")]
async fn refresh_workloads() -> impl Responder {
    tokio::spawn(async move {
        // This is the background task.
        // Since we're not waiting on it, we won't hold up the HTTP response.
        match services::workloads::fetch_and_update_all_watched().await {
            Ok(_) => println!("Workloads refreshed successfully."),
            Err(e) => eprintln!("Failed to refresh workloads: {}", e),
        }
    });
    HttpResponse::Ok().json(ApiResponse {
        status: "success".to_string(),
        message: "Workloads refreshed".to_string(),
    })
}

#[get("/api/workloads/{name}/{namespace}")]
async fn fetch_workload(path: web::Path<(String, String)>) -> impl Responder {
    let (name, namespace) = path.into_inner();
    //get workload from db
    if let Ok(workload) = database::client::return_workload(name.clone(), namespace.clone()) {
        log::info!("Workload found: {:?}", workload);
        HttpResponse::Ok().json(workload)
    } else if let Err(e) = database::client::return_workload(name, namespace) {
        log::error!("Error: {}", e);
        HttpResponse::Ok().json(e.to_string())
    } else {
        HttpResponse::Ok().body("Workload not found")
    }
}

#[get("api/workloads/update")]
async fn update_workload(req: HttpRequest) -> impl Responder {
    let query_params =
        web::Query::<HashMap<String, String>>::from_query(req.query_string()).unwrap();
    //create workload struct
    let name = query_params.get("name").unwrap();
    let namespace = query_params.get("namespace").unwrap();
    let image = query_params.get("image").unwrap();
    let current_version = query_params.get("current_version").unwrap();
    let latest_version = query_params.get("latest_version").unwrap();
    let git_ops_repo = query_params.get("git_ops_repo").unwrap();
    let git_directory = query_params.get("git_directory").unwrap();
    let workload = Workload {
        git_directory: Some(git_directory.clone()),
        name: name.clone(),
        exclude_pattern: None,
        git_ops_repo: Some(git_ops_repo.clone()),
        include_pattern: None,
        update_available: models::UpdateStatus::Available,
        image: image.clone(),
        last_scanned: "2021-08-01".to_string(),
        namespace: namespace.clone(),
        current_version: current_version.clone(),
        latest_version: latest_version.clone(),
    };

    log::info!("name = {}, namespace = {}, image = {}, current_version = {}, latest_version = {}, git_ops_repo = {}", name, namespace, image, current_version, latest_version, git_ops_repo);
    log::info!("Workload: {:?}", workload);
    run_git_operations(workload).unwrap();
    HttpResponse::Ok().body("Workload updated")
}

#[actix_web::main]
pub async fn site() -> std::io::Result<()> {
    log::info!("Starting web server");
    HttpServer::new(|| {
        App::new()
            .service(index)
            .service(fetch_all_workloads)
            .service(fetch_workload)
            .service(refresh_workloads)
            .service(update_workload)
    })
    .bind(("0.0.0.0", 8080))?
    .run()
    .await
}
