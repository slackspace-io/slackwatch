# Installation

### Helm
Reference files located in the `charts` directory, including a full values.yaml file with all available configuration options.

Helm repo is located at [https://slackspace-io.github.io/slackwatch/helm/](https://slackspace-io.github.io/slackwatch/installation.html#helm)

To install slackwatch using helm, run the following command:

```shell
helm repo add slackwatch https://slackspace-io.github.io/slackwatch/helm/
helm install slackwatch slackwatch/slackwatch-helm
```

Full values.yaml file with all available configuration options displayed below. It can be found in `charts/slackwatch/values-full.yaml`.
```
{{#include ../charts/slackwatch/values-full.yaml}}
```



### Kubernetes
Reference files located in the `k8s` directory.

Slackwatch is deployed as a deployment in the `slackwatch` namespace. The deployment is configured to watch all workloads in the cluster with the annotation `slackwatch.enabled` set to `true`.

The example files located here include the required service account, cluster role, and cluster role binding for slackwatch to watch all workloads in the cluster.

In addition to standard kubernetes files, the required fleet.yaml and kustomization.yaml files are included to deploy slackwatch using Rancher Fleet.


#### Running Locally
If running locally, slackwatch will use the kubeconfig file located at `~/.kube/config` to connect to the cluster. If you are using a different kubeconfig file, you can set the `KUBECONFIG` environment variable to the path of your kubeconfig file.


