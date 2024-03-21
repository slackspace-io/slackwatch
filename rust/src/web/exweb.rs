use crate::database;
use crate::database::client::get_latest_scan_id;
use crate::services;
use crate::services::workloads;
use actix_web::{get, web, App, HttpResponse, HttpServer, Responder};

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
    //fetch and update all workloads return if successful or error
    if let Ok(_) = services::workloads::fetch_and_update_all_watched().await {
        HttpResponse::Ok().body("Workloads refreshed")
    } else if let Err(e) = services::workloads::fetch_and_update_all_watched().await {
        HttpResponse::Ok().json(e.to_string())
    } else {
        HttpResponse::Ok().body("Workload not found")
    }
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

#[actix_web::main]
pub(crate) async fn site() -> std::io::Result<()> {
    HttpServer::new(|| {
        App::new()
            .service(index)
            .service(fetch_all_workloads)
            .service(fetch_workload)
            .service(refresh_workloads)
    })
    .bind(("0.0.0.0", 8080))?
    .run()
    .await
}
