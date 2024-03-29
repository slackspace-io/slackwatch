#![allow(non_snake_case, unused)]

use dioxus::prelude::*;
use dioxus::prelude::ServerFnError;
use wasm_bindgen_futures::spawn_local;
use crate::models;
use crate::models::models::Workload;

#[server] // For the server-side rendering version
async fn get_all_workloads() -> Result<String, ServerFnError> {
    use crate::database::client::return_all_workloads;
    let workloads = return_all_workloads();
    Ok(workloads.unwrap().iter().map(|w| w.name.clone()).collect::<Vec<String>>().join(", "))

//    Ok(workloads.unwrap().iter().map(|w| w.name.clone()).collect::<Vec<String>>().join(", "))
}

pub fn App() -> Element {
    let mut workloads = use_resource(get_all_workloads);
    let mut all = use_resource(get_all);
    println!("App started");
    rsx! {All{}}
//    rsx! {
//        "server data is {workloads():?}"
//        div {}
//        "server data is {all():?}"
//        div { {all().map(|w| rsx! { div {"{w:?}"}})}}
//
//
//}
}

#[component]
fn All() -> Element {
    let workloads = use_server_future(get_all)?;
    match workloads() {
        Some(Ok(workloads)) => {
            rsx! {
                div {"server data is {workloads:?}" }
                div { {workloads.iter().map(|w| rsx! { div {"{w.name}"} })}
                    }
            }
        },
        Some(Err(err)) => {
            rsx! {
                div { "Error: {err}" }
            }
        },
        _ => {
            rsx! {
                div { "Loading..." }
            }

        }
    }
}


#[server]
async fn get_all() -> Result<Vec<Workload>, ServerFnError> {
    use crate::database::client::return_all_workloads;
    let workloads = return_all_workloads();
    Ok(workloads.unwrap())
}
