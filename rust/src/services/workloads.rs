use rusqlite::Error;
use crate::database;
use crate::kubernetes::client::find_enabled_workloads;

pub async fn fetch_and_update_all_watched() -> Result<(), String> {
    let workloads = find_enabled_workloads().await.map_err(|e| e.to_string())?;
    for workload in workloads {
        // Assuming insert_workload is a synchronous operation
        // Consider using a thread pool or similar for actual concurrency
        std::thread::spawn(move || {
            database::client::insert_workload(&workload)
        }).join().map_err(|_| "Thread error".to_string())?.expect("TODO: panic message");
    }
    Ok(())
}


