#![allow(non_snake_case, unused)]

use dioxus::prelude::*;
use dioxus::prelude::ServerFnError;
use wasm_bindgen_futures::spawn_local;
use crate::models;

#[server] // For the server-side rendering version
async fn get_all_workloads() -> Result<String, ServerFnError> {
    use crate::database::client::return_all_workloads;
    let workloads = return_all_workloads();
    Ok(workloads.unwrap().iter().map(|w| w.name.clone()).collect::<Vec<String>>().join(", "))
}

pub fn App() -> Element {
    let mut workloads = use_resource(get_all_workloads);
    println!("App started");

    rsx! {
        "server data is {workloads():?}"
}
}
