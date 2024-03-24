use serde_derive::{Deserialize, Serialize};

// Define a structure to deserialize relevant parts of the YAML.
#[derive(Debug, Serialize, Deserialize)]
pub struct Deployment {
    apiVersion: String,
    kind: String,
    // You might need to adjust the structure based on the actual content of your YAML files
    spec: Spec,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct Spec {
    template: Template,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct Template {
    spec: PodSpec,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct PodSpec {
    containers: Vec<Container>,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct Container {
    image: String,
}
