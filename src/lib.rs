use crate::site::site::App;

pub mod site;

#[cfg(feature = "hydrate")]
#[wasm_bindgen::prelude::wasm_bindgen]
pub fn hydrate() {
    use leptos::*;
    log::info!("HYDRATE");
    console_error_panic_hook::set_once();

    mount_to_body(App);
}
