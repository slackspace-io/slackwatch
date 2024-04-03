use crate::config::Settings;
use crate::services::workloads::fetch_and_update_all_watched;
use config::ConfigError;
use cron::Schedule;
use std::str::FromStr;
use std::time::Duration;
use tokio::time::{sleep_until, Instant as TokioInstant};

#[cfg(feature = "server")]
pub async fn scheduler(schedule: &Schedule) {
    // Example cron schedule: Every minute
    log::info!("Scheduler started");
    //let schedule_str = "0 * * * * *"; // Adjust the cron expression as needed
    //let schedule = Schedule::from_str(schedule_str).expect("Failed to parse cron expression");
    println!("Cron schedule: {}", schedule);
    log::info!("Cron schedule: {}", schedule);
    // Find the next scheduled time
    //print next 5 scheduled times
    let mut i = 0;
    for datetime in schedule.upcoming(chrono::Utc) {
        if i < 5 {
            println!("Next scheduled time: {}", datetime);
            i += 1;
        } else {
            break;
        }
    }
    let now = chrono::Utc::now();
    if let Some(next) = schedule.upcoming(chrono::Utc).next() {
        let duration_until_next = (next - now).to_std().expect("Failed to calculate duration");
        println!("Next scheduled time: {:?}", next);
        // Convert std::time::Instant to tokio::time::Instant
        let tokio_now = TokioInstant::now();
        let tokio_future = tokio_now + duration_until_next;
        // Sleep until the next scheduled time
        sleep_until(tokio_future).await;
        // Execute your function
        refresh_all_workloads().await;
    }
}

#[cfg(feature = "server")]
pub async fn run_scheduler(settings: Settings) {
    //Load Scheduler
    let schedule_str = settings.system.schedule;
    let schedule = &Schedule::from_str(&schedule_str).expect("Failed to parse cron expression");
    loop {
        scheduler(schedule).await;
    }
}

#[cfg(feature = "server")]
async fn refresh_all_workloads() {
    log::info!("Refreshing all workloads");
    fetch_and_update_all_watched()
        .await
        .unwrap_or_else(|e| log::error!("Error refreshing workloads: {}", e));
}
