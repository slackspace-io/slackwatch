use crate::config::Settings;
use config::ConfigError;
use cron::Schedule;
use std::str::FromStr;
use std::time::Duration;
use tokio::time::{sleep_until, Instant as TokioInstant};

pub async fn scheduler(schedule: &Schedule) {
    // Example cron schedule: Every minute
    log::info!("Scheduler started");
    //let schedule_str = "0 * * * * *"; // Adjust the cron expression as needed
    //let schedule = Schedule::from_str(schedule_str).expect("Failed to parse cron expression");
    println!("Cron schedule: {}", schedule);
    log::info!("Cron schedule: {}", schedule);
    // Find the next scheduled time
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
        your_scheduled_function().await;
    }
}

pub async fn run_scheduler(settings: Result<Settings, ConfigError>) {
    //Load Scheduler
    let settings = settings.unwrap_or_else(|err| {
        log::error!("Failed to load settings: {}", err);
        panic!("Failed to load settings: {}", err);
    });
    let schedule_str = settings.system.schedule;
    let schedule = &Schedule::from_str(&schedule_str).expect("Failed to parse cron expression");
    loop {
        scheduler(schedule).await;
    }
}

async fn your_scheduled_function() {
    println!("Function is executed according to schedule.");
}
