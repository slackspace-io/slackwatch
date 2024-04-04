# SlackWatch: Your Kubernetes Container Guardian 🚀

### Features
- Notify you of new versions of your container images by comparing semver tags.
- Check for updates on a configured schedule.
- Authenticate with your container registry assuming standard docker registry auth behaviours.
- Trigger automated upgrade by commiting new tag to your gitops repo, if you have configured a gitops repo.
- Use include/exclude regex patterns to filter out images you are not interested in, these are set as annotations on the deployment.
- Notify you of new versions, as well sucessful commits via ntfy.
- A Web UI to view the status of your images, and trigger automated upgrades.

## What is SlackWatch?
Slackwatch is under heavy development. It is not yet ready for production use. I personally am using it, including the automated upgrade feature, but I am not recommending it for others yet.

I started slackwatch to solve my own problem of keeping track of container versions in my kubernetes cluster at home. I was usually notified by noticing a new release on github, but this did not mean the image was actually available.

I did try to use other tools, but they didn't quite fit my needs. I was primarily interested in comparing tag versions, since I do not run :latest on majority of my workloads. I also then liked the idea of triggering a commit to my gitops repo to trigger an upgrade.

I initially wrote this in go+nextjs, but decided this was a good opportunity to begin learning rust.

I am not saying rust is the best choice for this project, but I am enjoying learning it.

The stack now, and going forward is planned to be rust + using dioxus as the web framework.

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
