use rusqlite::Error;
use crate::database;
use crate::kubernetes::client::find_enabled_workloads;
use crate::models::models;
use crate::repocheck::repocheck::get_tags_for_image;

pub async fn fetch_and_update_all_watched() -> Result<(), String> {
    let workloads = find_enabled_workloads().await.map_err(|e| e.to_string())?;
    log::info!("Found {} workloads", workloads.len()); 
    //Update Database
    //stop after 3
    for workload in workloads {
        // Assuming insert_workload is a synchronous operation
        // Consider using a thread pool or similar for actual concurrency
        find_latest_tag_for_image(&workload).await.map_err(|e| e.to_string())?;
        std::thread::spawn(move || {
            database::client::insert_workload(&workload)
        }).join().map_err(|_| "Thread error".to_string())?.expect("TODO: panic message");
    }
    Ok(())
}

pub async fn find_latest_tag_for_image(workload: &models::Workload) -> Result<String, String> {
    let tags = get_tags_for_image(&workload.image).await.map_err(|e| e.to_string())?;
    let latest_tag = tags.first().ok_or("No tags found")?;
    log::info!("Latest tag for image {}: {}", workload.image, latest_tag);
    //update workload with latest tag
    Ok(latest_tag.clone())
}
