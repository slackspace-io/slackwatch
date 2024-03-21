use config::{Config, ConfigError, Environment, File};
use serde_derive::Deserialize;

#[derive(Debug, Deserialize)]
#[allow(unused)]
pub struct Settings {
    pub notifications: Notifications,
}

#[derive(Debug, Deserialize)]
#[allow(unused)]
pub struct Notifications {
    pub ntfy: Ntfy,
}

#[derive(Debug, Deserialize)]
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
            .add_source(File::with_name("config/config").required(false))
            .add_source(File::with_name(".env.toml").required(false))
            .add_source(Environment::with_prefix("slackwatch"))
            .build()?;
        //print config
        println!("{:?}", s);
        s.try_deserialize()
    }
}
