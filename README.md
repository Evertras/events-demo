# events-demo

AKA Enterprise Tic Tac Toe

## Prerequisites

Because Enterprise software is Enterprisey, there are a few things you'll need...

### Taskfile

[Get Taskfile](https://taskfile.dev/#/installation) to help run tasks more easily.  We're using
this instead of Make because Make can be a pain to install on Windows and its syntax
can be arcane for the uninitiated.  Taskfile is much clearer and is a simple cross-platform install.

### Kubernetes

You'll need Kubernetes running locally.  This can be with [Minikube](https://kubernetes.io/docs/setup/learning-environment/minikube/)
or in the Docker agent itself, if available.

If you can run `kubectl get all` on a shell and get a non-error response from your local cluster,
you're set!

### Helm

[Get Helm](https://helm.sh/docs/using_helm/) to help install things to Kubernetes in a more controlled way.  This is used in [the infra subdirectory](./infra).

### Telepresence

[Follow the instructions here](https://www.telepresence.io/reference/install) to install Telepresence.

Telepresence allows us to run a shell very quickly in Kubernetes which allows for easy development.
Each service should have a task to open a development shell that allows for direct running
via `go run cmd/thing/main.go`, or whatever you would normally run to start the service.

The advantage here is that we can bring up our entire stack in Kubernetes, then surgically replace
a running instance of a service with our development environment.  This means that we can run our
code in a functionally production setting to minimize "worked on my machine" differences and not
have to deploy to staging environments just to test our changes.

## Getting Started

Go to [the infra subdirectory](./infra) and run:

```bash
# Installs Traefik to act as an ingress for your Kubernetes cluster
task install-local-traefik

# This will set up a local Kafka stack in your Kubernetes cluster
task install-local-kafka
```

You can visit [the Traefik dashboard](http://dashboard.localhost/dashboard)
to see what services are available to the outside world.  As we add other services, they will
use the `<service>.localhost` convention.

[See this article](https://medium.com/@geraldcroes/kubernetes-traefik-101-when-simplicity-matters-957eeede2cf8) for a great introduction on using Traefik with Kubernetes.

## Tearing Down

Each `install-local-*` command has a `delete-local-*` counterpart.

```bash
task delete-local-kafka
task delete-local-traefik
```

Note that any `install` or `delete` can be run in any order at any time; none should
depend on the other to exist.

## Kafka connection information

For any clients that need to know a Kafka broker, supply the following:

Kafka brokers: `events-demo-kafka-cp-kafka-headless:9092`
Zookeeper: `events-demo-kafka-cp-zookeeper-headless:2181`

