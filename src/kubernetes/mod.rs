// kubernetes/mod.rs
#[cfg(feature = "server")]
pub mod client;

// Re-exporting Client so it's accessible from the kubernetes module directly.
