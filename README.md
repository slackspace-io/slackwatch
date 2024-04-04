# SlackWatch: Automated Kubernetes Image Updates

SlackWatch simplifies Kubernetes image management by keeping your deployments up-to-date.

Features

* **Image Version Monitoring:** Detects new container image versions (using semver tags) and sends notifications.
* **Scheduled Updates:** Runs checks on a configurable schedule.
* **Registry Authentication:** Integrates with standard container registry authentication.
* **GitOps Integration:** Triggers automated upgrades by committing new tags to your GitOps repository.
* **Filtering:** Uses include/exclude regex patterns for targeted updates.
* **Notifications:** Notifies you via ntfy about updates and successful commits.
* **Web UI:** Provides a dashboard for image status and manual upgrade triggers.

What is SlackWatch?

SlackWatch is a Kubernetes tool designed to streamline the process of keeping your container images up-to-date. It's ideal if you:

* Rely on version tags for image management (rather than `:latest`)
* Want a customizable and automated solution for image updates
* Prefer a GitOps-friendly workflow

## Example of a new version of ghostfolio

### Web View showing the new version.

<img alt="slackwatch-update.png" src="https://raw.githubusercontent.com/slackspace-io/slackwatch/main/.github/assets/img/slackwatch-update.png"/>

### Mobile notification via ntfy

<img alt="Screenshot_20240404-172936.png" src="https://raw.githubusercontent.com/slackspace-io/slackwatch/main/.github/assets/img/slackwatch_mobile_notification_ntfy.png"/>

## Trigger automated upgrade by commiting new tag to your gitops repo. 

### Commit completed

<img alt="slackwatch_commit_example.png" src="https://raw.githubusercontent.com/slackspace-io/slackwatch/main/.github/assets/img/slackwatch_commit_example.png"/>

### Notified of successful commit

<img alt="Screenshot_20240404-173514.png" src="https://raw.githubusercontent.com/slackspace-io/slackwatch/main/.github/assets/img/slackwatch_ntfy_commit_notification.png"/>


## Quick Start Guide

## Deployment (Into Kubernetes)

SlackWatch is Kubernetes-native and is intended to be deployed directly within your cluster. 
Check out the Kubernetes manifests in the `k8s` directory for deployment configurations. These are working files for Rancher Fleet, but should be easily adoptable to other methods.




### Prerequisites for running locally or developing

- Kubernetes cluster accessible from your local machine
- rust installed
- dioxus-cli installed, dx binary in your path.

### Deploying Slackwatch locally
- Deploy using references k8s files under k8s directory
- Run docker locally for testing
- Run `dx serve --platform fullstack` to start locally. Will use local kube-config. 

### Setting Up Slackwatch for development

1. **Clone the repository** to get started.
2. **Navigate to the git directory**: 
3. **Install dependencies**: Run `dx serve --platform fullstack` 

## Configuration

SlackWatch is designed to be highly configurable to fit your specific needs. You can adjust settings such as Kubernetes cluster details, notification channels, and more in the `config.yaml` file located in the backend directory.

```toml
[system]
schedule = "0 0 9-22/2 * * *"
data_dir = "/app/slackwatch/data"

[notifications.ntfy]
url = "http://localhost:9090"
topic = "slackwatch"
priorty = 1
reminder = "24h"
#The ntfy token should be an env variable named SLACKWATCH_NOTIFICATIONS.NTFY.TOKEN in k8s


[[gitops]]
##should match annotation on pod. This acts as a key to which gitops repo to use.
name = "fleet-slack-house"
repository_url = "https://github.com/slackspace-io/slackwatch.git"
branch = "main"
##k8s secret env variable which has your git repo's token
access_token_env_name = "SLACKWATCH_TOKEN"
commit_message = "Updated by slackwatch"
commit_name = "slackwatch"
commit_email = "slackwatch@slackspace.io"



```

## Contributing

Got ideas on how to make SlackWatch even better? Contributions are welcome! Whether it's adding new features, fixing bugs, or improving the documentation, your input is valuable.

## Stay in the Loop

We're constantly working on improving SlackWatch, adding new features, and refining the user experience. Keep an eye on the project for updates, and feel free to reach out with feedback or suggestions.

Join the SlackWatch community today and take the first step towards effortless container version management! 
