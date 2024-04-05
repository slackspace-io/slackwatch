## Introduction

Slackwatch is meant to notify you of new versions of your workloads in your Kubernetes cluster. 

It does this by comparing the tags of your workloads to the tags of the images in your container registry. If a new version is found, it will send a notification to a specified endpoint.

Slackwatch does not monitor the manifests of containers(currently). Which means on changes in semantic versioning of the container image, slackwatch will notify you.

Additionally, workloads using `:latest` cannot be watched by slackwatch.

### To get Started
Follow the [installation guide](installation.md). Review the configuration options in the [configuration guide](configuration.md), and set annotations on your workloads according to the [workload annotations](workload_annotations.md) documentation.
