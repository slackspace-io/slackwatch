//// src/config.rs
//use serde::{Deserialize, Serialize};
//use std::env;
//use std::path::PathBuf;
//
//#[derive(Debug, Serialize, Deserialize)]
//pub struct AppConfig {
//    pub database_path: PathBuf,
//    pub refresh_schedule: String,
//    pub refresh_interval: u64,
//    pub k8s_config_path: PathBuf,
//}
//
//impl AppConfig {
//    pub fn new() -> anyhow::Result<Self> {
//        dotenv::dotenv().ok(); // Load .env file if it exists
//        log::info!("Loading configuration");
//
//        let database_path = env::var("DATABASE_PATH")
//            .unwrap_or_else(|_| "./data.db".into())
//            .into();
//
//        let refresh_schedule = env::var("REFRESH_SCHEDULE").unwrap_or_else(|_| "0 0 * * *".into()); // Default to daily at midnight
//
//        Ok(AppConfig {
//            k8s_config_path: env::var("K8S_CONFIG_PATH")
//                .unwrap_or_else(|_| "~/.kube/config".into())
//                .into(),
//            database_path,
//            refresh_schedule,
//            refresh_interval: env::var("REFRESH_INTERVAL")
//                .unwrap_or_else(|_| "86400".into())
//                .parse::<u64>()?,
//        })
//    }
//}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_app_config_new() {
        // Setup
        let _ = dotenv::dotenv(); // Ensure .env is loaded for this test if it exists

        // This test assumes specific values are set in .env or defaults are acceptable
        let expected_db_path: PathBuf = env::var("DATABASE_PATH")
            .unwrap_or_else(|_| "./data.db".into())
            .into();
        let expected_schedule = env::var("REFRESH_SCHEDULE").unwrap_or_else(|_| "0 0 * * *".into());

        // Exercise
        let config = AppConfig::new().expect("Failed to create AppConfig");

        // Verify
        assert_eq!(config.database_path, expected_db_path);
        assert_eq!(config.refresh_schedule, expected_schedule);
    }
}
