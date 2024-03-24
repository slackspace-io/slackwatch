mod database;
mod models;
mod site;

use leptos::mount_to_body;

//#[cfg(feature = "hydrate")]
#[wasm_bindgen::prelude::wasm_bindgen]
pub fn hydrate() {
    use leptos::*;
    use site::app::*;

    console_error_panic_hook::set_once();

    mount_to_body(App);
}
