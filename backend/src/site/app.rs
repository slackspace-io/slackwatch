#![allow(non_snake_case, unused)]
use crate::models::models::Workload;
use dioxus::prelude::*;

pub fn App() -> Element {
    let mut count = use_signal(|| 0);
    let mut workload: Vec<Workload> = use_signal(|| vec![])();
    log::info!("App started");
    println!("App started");
    let all_workloads = get_all_workloads();

    rsx! {
        h1 { "High-Five counter: {count}" }
        button { onclick: move |_| count += 1, "Up high!" }
        button { onclick: move |_| count -= 1, "Down low!" }
        h2 { "{workload.len()}"},
        h2 { "This is a paragraph" }
        h2 {{ to_owned![all_workloads]}
         for w in all_workloads.map(|w| w.unwrap()){
            h2 { "{w.name}" }
        }
        },
        // call get all workload
        button {
            onclick: move |_| {
                to_owned![workload];
                async move {
                    if let Ok(workloads) = get_all_workloads().await {
                        for w in workloads {
                            println!("workload {:?}", w);
                            workload.push(w);
                        }
                    }
                }
            },
            "Get All Workloads"
        }

        button {
            onclick: move |_| {
                to_owned![count];
                async move {
                    if let Ok(new_count) = double_server(count()).await {
                        count.set(new_count);
                    }
                }
            },
            "Double"
        }


    }
}

#[server]
async fn get_all_workloads() -> Result<Vec<Workload>, ServerFnError> {
    // Perform some expensive computation or access a database on the server
    //    tokio::time::sleep(std::time::Duration::from_secs(1)).await;
    print!("Getting all workloads");
    use crate::database::client::return_all_workloads;
    let workloads = return_all_workloads()?;
    Ok(workloads)
}

#[server]
async fn double_server(number: i32) -> Result<i32, ServerFnError> {
    // Perform some expensive computation or access a database on the server
    //    tokio::time::sleep(std::time::Duration::from_secs(1)).await;
    log::info!("Doubling {number}");
    let result = number * 2;
    println!("server calculated {result}");
    Ok(result)
}
