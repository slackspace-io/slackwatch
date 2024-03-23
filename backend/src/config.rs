use config::{Config, ConfigError, Environment, File};
use serde_derive::Deserialize;

#[derive(Debug, Deserialize, Clone)]
#[allow(unused)]
pub struct Settings {
    pub notifications: Notifications,
    pub system: System,
    pub gitops: Vec<GitopsConfig>,
}

#[derive(Debug, Deserialize, Clone)]
#[allow(unused)]
pub struct System {
    pub schedule: String,
}

#[derive(Debug, Deserialize, Clone)]
#[allow(unused)]
pub struct GitopsConfig {
    pub name: String,
    pub repository_url: String,
    pub branch: String,
    pub commit_name: String,
    pub commit_email: String,
    pub access_token_env_name: String,
    pub commit_message: String,
}

#[derive(Debug, Deserialize, Clone)]
#[allow(unused)]
pub struct Notifications {
    pub ntfy: Ntfy,
}

#[derive(Debug, Deserialize, Clone)]
#[allow(unused)]
pub struct Ntfy {
    pub url: String,
    pub topic: String,
    pub reminder: String,
    pub token: String,
}

impl Settings {
    pub fn new() -> Result<Self, ConfigError> {
        dotenv::dotenv().ok();
        let s = Config::builder()
            .add_source(File::with_name("/app/config/config").required(false))
            .add_source(File::with_name(".env.toml").required(false))
            .add_source(Environment::with_prefix("slackwatch"))
            .build()?;
        //print config
        println!("{:?}", s);
        s.try_deserialize()
    }
    //add clone
}
