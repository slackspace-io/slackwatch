use crate::database;
use crate::database::client::get_latest_scan_id;
use crate::kubernetes::client::find_enabled_workloads;
use crate::models::models::{UpdateStatus, Workload};
use crate::notifications::ntfy::send_notification;
use crate::repocheck::repocheck::get_tags_for_image;
use regex::Regex;
use semver::Version;

pub async fn fetch_and_update_all_watched() -> Result<(), String> {
    let workloads = find_enabled_workloads().await.map_err(|e| e.to_string())?;
    log::info!("Found {} workloads", workloads.len());
    //Update Database
    //stop after 3
    let scan_id = get_latest_scan_id().unwrap_or(0) + 1;
    for workload in workloads {
        find_latest_tag_for_image(&workload)
            .await
            .map_err(|e| e.to_string())?;
        let workload = parse_tags(&workload).await.map_err(|e| e.to_string())?;
        if workload.update_available.to_string() == "Available" {
            send_notification(&workload)
                .await
                .unwrap_or_else(|e| log::error!("Error sending notification: {}", e));
        }

        std::thread::spawn(move || database::client::insert_workload(&workload, scan_id))
            .join()
            .map_err(|_| "Thread error".to_string())?
            .expect("TODO: panic message");
    }
    Ok(())
}

pub async fn find_latest_tag_for_image(workload: &Workload) -> Result<String, String> {
    let tags = get_tags_for_image(&workload.image)
        .await
        .map_err(|e| e.to_string())?;
    let latest_tag = tags.first().ok_or("No tags found")?;
    log::info!("Latest tag for image {}: {}", workload.image, latest_tag);
    //update workload with latest tag
    Ok(latest_tag.clone())
}

pub async fn test_call() {
    let workloads = find_enabled_workloads().await.unwrap();
    for workload in workloads.iter().take(1) {
        //let workload = workload.clone();
        let workload = parse_tags(&workload).await.unwrap();
        log::info!("Workload: {:?}", workload)
    }
}

fn strip_tag_lettings(tag: &str) -> String {
    tag.chars().skip_while(|c| !c.is_digit(10)).collect()
}

pub async fn parse_tags(workload: &Workload) -> Result<Workload, Box<dyn std::error::Error>> {
    let mut tags = get_tags_for_image(&workload.image).await?;
    tags.sort();

    // Include Pattern Handling
    if let Some(include_pattern_str) = &workload.include_pattern {
        log::info!("Include pattern defined, using only include");
        log::info!("Include pattern: {}", include_pattern_str);

        // Build regex from patterns, assuming comma-separated
        let include_patterns = include_pattern_str
            .split(",")
            .map(|pattern| Regex::new(pattern).unwrap()) // Compile each regex
            .collect::<Vec<Regex>>();

        tags = tags
            .into_iter()
            .filter(|tag| include_patterns.iter().any(|regex| regex.is_match(tag)))
            .collect();

        log::info!("Filtered tags: {:?}", tags);
    } else if let Some(exclude_pattern_str) = &workload.exclude_pattern {
        log::info!("Exclude pattern defined, using only exclude");
        log::info!("Exclude pattern: {}", exclude_pattern_str);

        // Build regex from patterns, assuming comma-separated
        let exclude_patterns = exclude_pattern_str
            .split(",")
            .map(|pattern| Regex::new(pattern).unwrap()) // Compile each regex
            .collect::<Vec<Regex>>();

        tags = tags
            .into_iter()
            .filter(|tag| exclude_patterns.iter().all(|regex| !regex.is_match(tag)))
            .collect();

        log::info!("Filtered tags: {:?}", tags);
    }
    let current_version = Version::parse(&strip_tag_lettings(&workload.current_version))
        .unwrap_or_else(|_| Version::new(0, 0, 0));

    // Perform SemVer comparison with each tag:
    let mut latest_version = String::new();
    let mut update_available = UpdateStatus::NotAvailable;
    for tag in tags {
        if let Ok(tag_version) = Version::parse(&strip_tag_lettings(&tag)) {
            if tag_version > current_version {
                // tag_version is greater than current_version
                // Do something with this tag
                println!("Tag {} is newer than current version", tag);
                latest_version = tag.clone();
            }
        } else {
            // Handle the case where the tag is not a valid SemVer format
            println!("Tag {} is not a valid SemVer", tag);
            println!(
                "Tag {} is not a valid SemVer - stripped",
                &strip_tag_lettings(&tag)
            );
        }
    }
    if !latest_version.is_empty() {
        log::info!("Latest version for {}: {}", workload.image, latest_version);
        update_available = UpdateStatus::Available;
    }
    Ok(Workload {
        name: workload.name.clone(),
        exclude_pattern: workload.exclude_pattern.clone(),
        git_ops_repo: workload.git_ops_repo.clone(),
        include_pattern: workload.include_pattern.clone(),
        namespace: workload.namespace.clone(),
        current_version: workload.current_version.clone(),
        image: workload.image.clone(),
        update_available,
        last_scanned: workload.last_scanned.clone(),
        latest_version: latest_version.clone(),
    })
}
