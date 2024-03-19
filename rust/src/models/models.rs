use serde::Serialize;

//Data model for workload
#[derive(Debug, Serialize)]
pub struct Workload {
    pub name: String,
    pub exclude_pattern: Option<String>,
    pub git_ops_repo: Option<String>,
    pub include_pattern: Option<String>,
    pub update_available: UpdateStatus,
    pub image: String,
    pub last_scanned: String,
    pub namespace: String,
    pub current_version: String,
    pub latest_version: String,
}

#[derive(strum_macros::Display, strum_macros::EnumString, Debug, Serialize)]
pub enum UpdateStatus {
    Available,
    NotAvailable,
}
