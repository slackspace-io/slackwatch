#![allow(non_snake_case, unused)]

use dioxus::prelude::*;
use dioxus::prelude::ServerFnError;
use serde_derive::{Deserialize, Serialize};
use wasm_bindgen_futures::spawn_local;
use crate::config::{GitopsConfig, Notifications, Settings, System};
use crate::models;
use crate::models::models::Workload;


#[derive(Routable, Clone)]
#[rustfmt::skip]
enum Route {
        #[route("/")]
        Home {},
        #[route("/refresh-all")]
        RefreshAll {},
        #[route("/settings")]
        SettingsPage {},
    }

#[derive(Debug, Deserialize, Clone)]
#[allow(unused)]
pub struct AppSettings {
    pub settings: Settings,
}


#[derive(PartialEq, Clone,Props)]
struct WorkloadCardProps {
    workload: Workload,
}


#[server] // For the server-side rendering version
async fn get_all_workloads() -> Result<String, ServerFnError> {
    use crate::database::client::return_all_workloads;
    let workloads = return_all_workloads();
    Ok(workloads.unwrap().iter().map(|w| w.name.clone()).collect::<Vec<String>>().join(", "))
}



#[component]
fn SettingsCard(props: AppSettings) -> Element {
    rsx! {
        div {
            class: "settings-section",
        }
    }
}

#[component]
fn SettingsPage() -> Element {
    let settings_context = use_context::<Signal<AppSettings>>();
    let settings = settings_context.read();
    rsx! {
        //div {
        //  NextScheduledTimeCard {}
        //
        //},
        div {
            class: "settings-page",
            div {
                class: "settings-section",
                div { class: "settings-section-header", "System Settings" },
                div { class: "settings-item",
                    span { class: "settings-item-key", "Schedule: " },
                    span { class: "settings-item-value", "{settings.settings.system.schedule}" }
                },
                div { class: "settings-item",
                    span { class: "settings-item-key", "Data Directory: " },
                    span { class: "settings-item-value", "{settings.settings.system.data_dir}" }
                },
                div { class: "settings-item",
                    span { class: "settings-item-key", "Run at Startup: " },
                    span { class: "settings-item-value", "{settings.settings.system.run_at_startup}" }
                }
            },
            div {
                class: "settings-section",
                div { class: "settings-section-header", "Gitops Settings" },
                    for gitops in settings.clone().settings.gitops.unwrap().iter() {
                        div { class: "settings-item",
                            span { class: "settings-item-key", "Name: " }
                            span { class: "settings-item-value", "{gitops.name}" }
                        }
                        div { class: "settings-item",
                            span { class: "settings-item-key", "Repository URL: " }
                            span { class: "settings-item-value", "{gitops.repository_url}" }
                        }
                }

            }

        },
    }
}




#[component]
fn Home() -> Element {
    let workloads = use_server_future(get_all)?;
    match workloads() {
        Some(Ok(workloads)) => {
            if workloads.len() == 0 {
                rsx! {
                    div { "No workloads found" } ,
                    br {},
                    br {},
                    a { href: "/refresh-all", "Click to Refresh All" }

                }
            } else {
            rsx! {
                div { class: "workloads-page",
                    NextScheduledTimeCard {},
                    for w in workloads.iter() {
                    WorkloadCard{workload: w.clone()}
                        }
            }
            }
        }
        },
        Some(Err(err)) => {
            rsx! { div { "Error: {err}" } }
        },
        _ => {
            rsx! { div { "Loading..." } }
        }

    }
}

#[component]
fn DebugWorkloadCard(props: WorkloadCardProps) -> Element {
    rsx! {
        div {
            class: if props.workload.update_available == models::models::UpdateStatus::Available {
                "workload-card-update-available"
            } else {
                "workload-card"
            },
            div { class: "workload-name", "{props.workload.name}" },
            div { class: "workload-
            namespace", "Namespace: {props.workload.namespace}" },
            div { class: "workload-version", "Current Tag {props.workload.current_version}" },
            div { class: "workload-image", "Image: {props.workload.image}" },
            div { class: "workload-last-scanned", "Last Scanned: {props.workload.last_scanned}" },
            if props.workload.update_available == models::models::UpdateStatus::Available {
                div { class: "workload-latest-version", "Latest Version Available: {props.workload.latest_version}" }
                br {}
            }
        }
    }
}



#[component]
fn WorkloadCard(props: WorkloadCardProps) -> Element {
    let data = use_signal(|| {props.workload.clone()});
    rsx! {
        div {
            class: if props.workload.update_available == models::models::UpdateStatus::Available {
                "workload-card-update-available"
            } else {
                "workload-card"
            },
            div { class: "workload-name", "{props.workload.name}" },
            button {onclick: move |_| {
                to_owned![data, props.workload];
                async move {
                    if let Ok(_) = update_workload(data()).await {
                    }
                }
            },
              class: "workload-update-single",  "Refresh"},
            div { class: "workload-namespace", "Namespace: {props.workload.namespace}" },
            div { class: "workload-version", "Current Tag {props.workload.current_version}" },
            div { class: "workload-image", "Image: {props.workload.image}" },
            div { class: "workload-last-scanned", "Last Scanned: {props.workload.last_scanned}" },
            if props.workload.update_available == models::models::UpdateStatus::Available {
                div { class: "workload-latest-version", "Latest Version Available: {props.workload.latest_version}" }
                br {}
                button { onclick: move |_| {
                    async move {
                        if let Ok(_) = upgrade_workload(data()).await {
                        }
                    }
                },
                    class: "upgrade-button", "Upgrade"}
            }
        }
    }
}

pub fn App() -> Element {
    println!("App started");
    use_context_provider(|| Signal::new(AppSettings { settings: Settings::new().unwrap_or_else(
        |err| {
            log::error!("Failed to load settings: {}", err);
            panic!("Failed to load settings: {}", err);
        }
    ) }));
    //load config
    rsx! { Router::<Route> {} }
}


#[component]
fn RefreshAll() -> Element {
    let refresh = use_server_future(refresh_all)?;
    match refresh() {
        Some(Ok(_)) => {
            rsx! {
                div { "Refreshed" },
                br {},
                br {},
                a { href: "/", "Go back to Home" }
            }
        },
        Some(Err(err)) => {
            rsx! { div { "Error: {err}" } }
        }
        _ => rsx! { div { "Loading..." } }
    }
}


#[component]
fn NotFound() -> Element {
    rsx! { div { "Not Found" } }
}


#[component]
fn All() -> Element {
    let workloads = use_server_future(get_all)?;
    match workloads() {
        Some(Ok(workloads)) => {
            rsx! {
                div { class: "workloads-container",
                    div {
                        {workloads.iter().map(|w|
                            rsx!{
                            div {
                            class: if w.update_available == models::models::UpdateStatus::Available {
                                "workload-card-update-available"
                            } else {
                                "workload-card"
                            },
                            div { class: "workload-name", "{w.name}" },
                            div { class: "workload-namespace", "Namespace: {w.namespace}" },
                            div { class: "workload-version", "Current Tag {w.current_version}" },
                            div { class: "workload-image", "Image: {w.image}" },
                            div { class: "workload-last-scanned", "Last Scanned: {w.last_scanned}" },
                                        div { class: "workload-name", "{w.latest_version}" },
                            if w.update_available == models::models::UpdateStatus::Available {
                                div { class: "workload-update-available", "Update Available" }
                            }
                            }
                            },

                        )},
                    }
                }
            }
        },
        Some(Err(err)) => {
            rsx! { div { "Error: {err}" } }
        },
        _ => {
            rsx! { div { "Loading..." } }

        }
    }
}


// ... rest of the code ...

#[component]
fn NextScheduledTimeCard() -> Element {
    let settings_context = use_context::<Signal<AppSettings>>().clone();
    let mut next_schedule = use_server_future(move || async move {
        let settings = settings_context.read();
        get_next_schedule_time(settings.settings.clone()).await
    })?;
    match next_schedule() {
        Some(Ok(next_schedule)) => {
            rsx! {
                div { class: "next-scheduled-time",
                    div { class: "system-info", "System Info" },
                    div { "Next Run: {next_schedule}" }
                }
            }
        },
        Some(Err(err)) => {
            rsx! { div { "Error: {err}" } }
        },
        None => {
            rsx! { div { "Loading..." } }
        }
        _ => {
            rsx! { div { "Loading..." } }
        }
    }
}

#[server]
async fn get_next_schedule_time(settings: Settings) -> Result<String, ServerFnError> {
    use crate::services::scheduler::next_schedule_time;
    let schedule_str = &settings.system.schedule;
    let next_schedule = next_schedule_time(&schedule_str).await;
    log::info!("get_next_schedule_time: {:?}", next_schedule);
    if next_schedule.contains("No upcoming schedule") {
        Result::Err(ServerFnError::new(&next_schedule))
    } else {
        Result::Ok(next_schedule)
    }
}


#[server]
async fn upgrade_workload(workload: Workload) -> Result<(), ServerFnError> {
    log::info!("upgrade_workload: {:?}", workload);
    println!("upgrade_workload");
       use crate::gitops::gitops::run_git_operations;
    let result = run_git_operations(workload).await;
    Ok(result.unwrap())
}

#[server]
async fn update_workload(workload: Workload) -> Result<(), ServerFnError> {
    log::info!("update_workload: {:?}", workload);
    println!("update_workload");
    use crate::services::workloads::update_single_workload;
    match update_single_workload(workload).await {
        Ok(result) => Ok(result),
        Err(e) => {
            log::error!("Failed to update workload: {}", e);
            Err(ServerFnError::new(&format!("Failed to update workload: {}", e)))
        }
    }
}


#[server]
async fn refresh_all() -> Result<(), ServerFnError> {
    use crate::services::workloads::fetch_and_update_all_watched;
    let result = fetch_and_update_all_watched().await;
    Ok(result.unwrap())
}


#[server]
async fn get_all() -> Result<Vec<Workload>, ServerFnError> {
    use crate::database::client::return_all_workloads;
    let settings_context = consume_context::<Signal<AppSettings>>();
    log::info!("settings_context: {:?}", settings_context);
    let workloads = return_all_workloads();
    log::info!("get_all_workloads: {:?}", workloads);
    Ok(workloads.unwrap())
}
