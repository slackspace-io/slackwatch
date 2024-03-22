use crate::config::Settings;
use crate::models::models::Workload;
use ntfy::payload::{Action, ActionType};
use ntfy::{Auth, Dispatcher, NtfyError, Payload, Priority};

pub async fn send_notification(workload: &Workload) -> Result<(), NtfyError> {
    //get settings
    let settings = Settings::new().unwrap();
    let token = settings.notifications.ntfy.token;
    let topic = settings.notifications.ntfy.topic;
    let url = settings.notifications.ntfy.url;

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
    Ok(())
}
