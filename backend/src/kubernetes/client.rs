// kubernetes/client.rs
use crate::models::models::{UpdateStatus, Workload};
use k8s_openapi::api::core::v1::Pod;
use k8s_openapi::chrono;
use kube::{
    api::{Api, ListParams},
    Client as KubeClient, ResourceExt,
};

pub struct Client {
    kube_client: KubeClient,
}

impl Client {
    pub async fn new() -> Result<Self, kube::Error> {
        let kube_client = KubeClient::try_default().await?;
        Ok(Client { kube_client })
    }

   // pub async fn fetch_slackwatch_enabled_containers(&self) -> Result<Vec<Pod>, kube::Error> {
   //     let pods: Api<Pod> = Api::all(self.kube_client.clone());
   //     let lp = ListParams::default().labels("slackwatch.enable=true"); // Adjust based on actual use case
   //     let pod_list = pods.list(&lp).await?;
   //     Ok(pod_list.items)
   // }
   // pub async fn fetch_containers_with_annotation(
   //     &self,
   //     annotation_key: &str,
   // ) -> Result<Vec<Pod>, kube::Error> {
   //     let pods: Api<Pod> = Api::all(self.kube_client.clone());
   //     let lp = ListParams::default().labels(format!("{}=*", annotation_key).as_str()); // Adjust based on actual use case
   //     let pod_list = pods.list(&lp).await?;
   //     Ok(pod_list.items)
   // }

    pub async fn list_pods(&self) -> Result<Vec<Pod>, kube::Error> {
        let pods: Api<Pod> = Api::all(self.kube_client.clone());
        let lp = ListParams::default();
        let pod_list = pods.list(&lp).await?;
        Ok(pod_list.items)
    }
}

pub async fn find_enabled_workloads() -> Result<Vec<Workload>, kube::Error> {
    //get all pods
    let client = Client::new().await?;
    let pods = client.list_pods().await?;
    let mut workloads = Vec::new();
    //count number of pods
    for p in pods {
        if let Some(namespace) = p.namespace() {
            if let Some(annotations) = p.metadata.annotations {
                if let Some(enable) = annotations.get("slackwatch.enable") {
                    //Found a pod with slackwatch.enable annotation
                    if enable == "true" {
                        let exclude_pattern = annotations.get("slackwatch.exclude").cloned();
                        let include_pattern = annotations.get("slackwatch.include").cloned();
                        let git_ops_repo = annotations.get("slackwatch.repo").cloned();
                        for spec in p.spec {
                            for container in spec.containers.clone() {
                                if let Some(name) = Some(container.name) {
                                    if let Some(image) = container.image {
                                        let parts = image.split(":").collect::<Vec<&str>>();
                                        let current_version = parts.get(1).unwrap_or(&"latest");
                                        workloads.push(Workload {
                                            exclude_pattern: exclude_pattern.clone(),
                                            git_ops_repo: git_ops_repo.clone(),
                                            include_pattern: include_pattern.clone(),
                                            update_available: "NotAvailable"
                                                .parse()
                                                .unwrap_or(UpdateStatus::NotAvailable),
                                            image: image.clone(),
                                            name: name.clone(),
                                            namespace: namespace.clone(),
                                            current_version: current_version.to_string(),
                                            last_scanned: chrono::Utc::now().to_rfc3339(),
                                            latest_version: "1.0.0".to_string(),
                                        });
                                    }
                                }
                            }
                        }
                    }
                }
            }
        }
    }
    Ok(workloads)
}
