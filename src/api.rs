use warp::{Filter, Rejection, Reply};
use warp::filters::cors::cors;
use warp::http::Method;
use serde_json::json;
use crate::models::models::Workload;
use crate::config::Settings;
use crate::services::workloads::{fetch_and_update_all_watched, update_single_workload};
use crate::gitops::gitops::run_git_operations;
use crate::services::scheduler::next_schedule_time;
use crate::database::client::return_all_workloads;

pub async fn start_api_server() {
    // CORS configuration
    let cors = cors()
        .allow_any_origin()
        .allow_methods(&[Method::GET, Method::POST, Method::OPTIONS])
        .allow_headers(vec!["Content-Type"]);

    // API routes
    let api = warp::path("api");

    // GET /api/workloads - Get all workloads
    let get_workloads = api
        .and(warp::path("workloads"))
        .and(warp::get())
        .and_then(handle_get_workloads);

    // POST /api/workloads/update - Update a workload
    let update_workload = api
        .and(warp::path("workloads"))
        .and(warp::path("update"))
        .and(warp::post())
        .and(warp::body::json())
        .and_then(handle_update_workload);

    // POST /api/workloads/upgrade - Upgrade a workload
    let upgrade_workload = api
        .and(warp::path("workloads"))
        .and(warp::path("upgrade"))
        .and(warp::post())
        .and(warp::body::json())
        .and_then(handle_upgrade_workload);

    // POST /api/workloads/refresh-all - Refresh all workloads
    let refresh_all = api
        .and(warp::path("workloads"))
        .and(warp::path("refresh-all"))
        .and(warp::post())
        .and_then(handle_refresh_all);

    // GET /api/settings - Get settings
    let get_settings = api
        .and(warp::path("settings"))
        .and(warp::get())
        .and_then(handle_get_settings);

    // GET /api/settings/next-schedule-time - Get next schedule time
    let get_next_schedule = api
        .and(warp::path("settings"))
        .and(warp::path("next-schedule-time"))
        .and(warp::get())
        .and_then(handle_get_next_schedule);

    // Serve static files from the frontend/dist directory
    let static_files = warp::fs::dir("frontend/dist");

    // Fallback route for SPA - serve index.html for any other route
    let spa_fallback = warp::any()
        .and(warp::fs::file("frontend/dist/index.html"));

    // Combine all routes
    let routes = get_workloads
        .or(update_workload)
        .or(upgrade_workload)
        .or(refresh_all)
        .or(get_settings)
        .or(get_next_schedule)
        .or(static_files)
        .or(spa_fallback)
        .with(cors);

    // Start the server
    log::info!("Starting API server on 0.0.0.0:8080");
    warp::serve(routes)
        .run(([0, 0, 0, 0], 8080))
        .await;
}

async fn handle_get_workloads() -> Result<impl Reply, Rejection> {
    match return_all_workloads() {
        Ok(workloads) => Ok(warp::reply::json(&workloads)),
        Err(e) => {
            log::error!("Failed to get workloads: {}", e);
            let error = json!({ "error": format!("Failed to get workloads: {}", e) });
            Ok(warp::reply::json(&error))
        }
    }
}

async fn handle_update_workload(workload: Workload) -> Result<impl Reply, Rejection> {
    match update_single_workload(workload).await {
        Ok(_) => Ok(warp::reply::json(&json!({ "status": "success" }))),
        Err(e) => {
            log::error!("Failed to update workload: {}", e);
            let error = json!({ "error": format!("Failed to update workload: {}", e) });
            Ok(warp::reply::json(&error))
        }
    }
}

async fn handle_upgrade_workload(workload: Workload) -> Result<impl Reply, Rejection> {
    match run_git_operations(workload).await {
        Ok(_) => Ok(warp::reply::json(&json!({ "status": "success" }))),
        Err(e) => {
            log::error!("Failed to upgrade workload: {}", e);
            let error = json!({ "error": format!("Failed to upgrade workload: {}", e) });
            Ok(warp::reply::json(&error))
        }
    }
}

async fn handle_refresh_all() -> Result<impl Reply, Rejection> {
    match fetch_and_update_all_watched().await {
        Ok(_) => Ok(warp::reply::json(&json!({ "status": "success" }))),
        Err(e) => {
            log::error!("Failed to refresh all workloads: {}", e);
            let error = json!({ "error": format!("Failed to refresh all workloads: {}", e) });
            Ok(warp::reply::json(&error))
        }
    }
}

async fn handle_get_settings() -> Result<impl Reply, Rejection> {
    match Settings::new() {
        Ok(settings) => Ok(warp::reply::json(&settings)),
        Err(e) => {
            log::error!("Failed to get settings: {}", e);
            let error = json!({ "error": format!("Failed to get settings: {}", e) });
            Ok(warp::reply::json(&error))
        }
    }
}

async fn handle_get_next_schedule() -> Result<impl Reply, Rejection> {
    match Settings::new() {
        Ok(settings) => {
            let schedule_str = &settings.system.schedule;
            let next_schedule = next_schedule_time(&schedule_str).await;
            // Ensure we're returning a string, not an object
            Ok(warp::reply::json(&next_schedule))
        },
        Err(e) => {
            log::error!("Failed to get settings for next schedule: {}", e);
            let error = json!({ "error": format!("Failed to get settings for next schedule: {}", e) });
            Ok(warp::reply::json(&error))
        }
    }
}
