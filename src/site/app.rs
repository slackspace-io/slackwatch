#![allow(non_snake_case, unused)]

use dioxus::prelude::*;
use dioxus::prelude::ServerFnError;
use wasm_bindgen_futures::spawn_local;
use crate::models;
use crate::models::models::Workload;


#[derive(Routable, Clone)]
#[rustfmt::skip]
enum Route {
        #[route("/")]
        Home {},
        #[route("/refresh-all")]
        RefreshAll {},
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
    let workloads = return_all_workloads();
    log::info!("get_all_workloads: {:?}", workloads);
    Ok(workloads.unwrap())
}
