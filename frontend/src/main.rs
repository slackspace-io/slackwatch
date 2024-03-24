#[tokio::main]
async fn main() {
    tokio::task::spawn_blocking(|| {
        println!("Site started");
        let _ = frontend::server::start_site();
        //        let _ = web::exweb::site();
    })
    .await
    .expect("Failed to run site")
    // tokio::task::spawn(frontend::server::start_site());
}

#[cfg(not(feature = "ssr"))]
pub fn main() {
    // no client-side main function
    // unless we want this to work with e.g., Trunk for pure client-side testing
    // see lib.rs for hydration function instead
}
