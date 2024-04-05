# Installation

### Kubernetes
Reference files located in the `k8s` directory.

Slackwatch is deployed as a deployment in the `slackwatch` namespace. The deployment is configured to watch all workloads in the cluster with the annotation `slackwatch.enabled` set to `true`.

The example files located here include the required service account, cluster role, and cluster role binding for slackwatch to watch all workloads in the cluster.

In addition to standard kubernetes files, the required fleet.yaml and kustomization.yaml files are included to deploy slackwatch using Rancher Fleet.


#### Running Locally
If running locally, slackwatch will use the kubeconfig file located at `~/.kube/config` to connect to the cluster. If you are using a different kubeconfig file, you can set the `KUBECONFIG` environment variable to the path of your kubeconfig file.


