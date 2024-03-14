# SlackWatch: Your Kubernetes Container Guardian 🚀

Welcome to SlackWatch, the Kubernetes-centric solution designed to keep you informed about your container images and updates directly through your favorite notification channel! Whether you're managing a small project or orchestrating a fleet of containers, SlackWatch is here to ensure you never miss an update.

## Key Features

- **Container Monitoring**: Effortlessly monitor the containers in your Kubernetes cluster.
- **Version Comparison**: Automatically compare your current container images against external repositories for newer versions.
- **Notifications**: Receive timely notifications when it's time to update, so you're always ahead of the game.
- **Kubernetes Integration**: Crafted to run within a Kubernetes environment, SlackWatch fits right into your existing workflow.

## Quick Start Guide

### Prerequisites

- Kubernetes cluster
- Node.js and npm for the frontend

### Setting Up the Backend

1. **Clone the repository** to get started.
2. **Navigate to the backend directory**: `cd backend`.
3. **Build the Go application**: Run `go build -o slackwatch-backend`.
4. **Run the backend**: Execute `./slackwatch-backend`. This starts the backend server, which begins monitoring your Kubernetes cluster.

### Setting Up the Frontend

1. **Navigate to the frontend directory**: From the root of the repository, `cd frontend`.
2. **Install dependencies**: Run `npm install` to get all the necessary packages.
3. **Start the development server**: Execute `npm run dev`.
4. **Open your browser** and go to [http://localhost:3000](http://localhost:3000) to view the SlackWatch dashboard.

## Configuration

SlackWatch is designed to be highly configurable to fit your specific needs. You can adjust settings such as Kubernetes cluster details, notification channels, and more in the `config.yaml` file located in the backend directory.

```yaml
kubernetes:
  pollingInterval: 30
  useInClusterConfig: true
notifications:
  ntfy:
    url: "http://ntfy.sh"
    topic: "slackwatch-updates"
```

## Contributing

Got ideas on how to make SlackWatch even better? Contributions are welcome! Whether it's adding new features, fixing bugs, or improving the documentation, your input is valuable.

## Deployment

SlackWatch is Kubernetes-native and can be deployed directly within your cluster. Check out the Kubernetes manifests in the `k8s` directory for deployment configurations.

## Stay in the Loop

We're constantly working on improving SlackWatch, adding new features, and refining the user experience. Keep an eye on the project for updates, and feel free to reach out with feedback or suggestions.

Join the SlackWatch community today and take the first step towards effortless container version management! 
