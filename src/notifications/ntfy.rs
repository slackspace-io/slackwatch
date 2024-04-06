use dioxus::html::set;
use dioxus::prelude::client;
use dioxus::prelude::server_fn::response::Res;
use crate::config::{Ntfy, Settings};
use crate::models::models::Workload;
use ntfy::payload::{Action, ActionType};
use ntfy::{Auth, Dispatcher, NtfyError, Payload, Priority};
use wasm_bindgen_futures::js_sys::Atomics::load;
use wasm_bindgen_futures::js_sys::Set;


pub async fn notify_commit(workload: &Workload) -> Result<(), NtfyError> {
    //get settings
    match load_settings() {
        Ok(settings) => {
            let url = settings.url;
            let topic = settings.topic;
            let token = settings.token;


            let dispatcher = Dispatcher::builder(&url)
                .credentials(Auth::new("", &token)) // Add optional credentials
                .build()?; // Build dispatcher

            //let action = Action::new(
            //    ActionType::Http,
            //    "Acknowledge",
            //    Url::parse(&url)?,
            //);

            //make message for payload about new container update
            let message = format!(
                "Deployment {} has been updated to version {}",
                workload.name, workload.latest_version
            );

            let payload = Payload::new(&topic)
                .message(message) // Add optional message
                .title(&workload.name) // Add optiona title
                .tags(["Update"]) // Add optional tags
                .priority(Priority::Default) // Edit priority
                //.actions([action]) // Add optional actions
                //.click(Url::parse("https://example.com")?) // Add optional clickable url
                //.attach(Url::parse("https://example.com/file.jpg")?) // Add optional url attachment
                //.delay(Local::now() + Duration::minutes(1)) // Add optional delay
                .markdown(true); // Use markdown

            match dispatcher.send(&payload).await {
                Ok(_) => log::info!("Payload sent successfully."),
                Err(e) => log::error!("Failed to send payload: {}", e),
            }
            log::info!("Notification sent");
            Ok(())// Proceed with using settings
        },
        Err(e) => {
            log::info!("Failed to load settings: {}", e);
            Ok(())
            // Handle the error, e.g., by returning or panicking
        }
    }


}

fn load_settings() ->Result<Ntfy, String> {
    //get settings
    let settings = Settings::new().unwrap_or_else(|err| {
        log::error!("Failed to load settings: {}", err);
        panic!("Failed to load settings: {}", err);
    });
    if let Some(notifications) = settings.notifications {
        if let Some(ntfy_config) = notifications.ntfy {
            return Ok(ntfy_config.clone());
        } else {
            return Err("No Ntfy Config Found".to_string());
        }
    } else {
        return Err("No Notifications Config Found".to_string());
    }

}

pub async fn send_notification(workload: &Workload) -> Result<(), NtfyError> {
    //get settings
    match load_settings() {
        Ok(settings) => {
            let url = settings.url;
            let topic = settings.topic;
            let token = settings.token;

            let dispatcher = Dispatcher::builder(&url)
                .credentials(Auth::new("", &token)) // Add optional credentials
                .build()?; // Build dispatcher

            //let action = Action::new(
            //    ActionType::Http,
            //    "Acknowledge",
            //    Url::parse(&url)?,
            //);

            //make message for payload about new container update
            let message = format!(
                "Update Available: {} From {} to {}",
                workload.name, workload.current_version, workload.latest_version
            );

            let payload = Payload::new(&topic)
                .message(message) // Add optional message
                .title(&workload.name) // Add optiona title
                .tags(["Update"]) // Add optional tags
                .priority(Priority::High) // Edit priority
                //.actions([action]) // Add optional actions
                //.click(Url::parse("https://example.com")?) // Add optional clickable url
                //.attach(Url::parse("https://example.com/file.jpg")?) // Add optional url attachment
                //.delay(Local::now() + Duration::minutes(1)) // Add optional delay
                .markdown(true); // Use markdown

            match dispatcher.send(&payload).await {
                Ok(_) => log::info!("Payload sent successfully."),
                Err(e) => log::error!("Failed to send payload: {}", e),
            }
            log::info!("Notification sent");
            Ok(()) // Proceed with using settings
        },
        Err(e) => {
            log::info!("Failed to load settings: {}", e);
            Ok(()) // Handle the error, e.g., by returning or panicking
            // Handle the error, e.g., by returning or panicking
        }
    }


}
