# events-demo

AKA Enterprise Tic Tac Toe

*Note*: Since this is for proof of concept, everything is crammed into a single repo.  In a real world
scenario each subdirectory would almost certainly be its own repository.

## Prerequisites

Because Enterprise software is Enterprisey, there are a few things you'll need.

See [the tech choice document](./Tech.md) for reasoning behind each of these.

### Go

Get yourself some [Go](https://golang.org/doc/install), 1.13 or higher.

### Taskfile

[Get Taskfile](https://taskfile.dev/#/installation).

### Kubernetes

You'll need Kubernetes running locally.  This can be with [Minikube](https://kubernetes.io/docs/setup/learning-environment/minikube/)
or in the Docker agent itself, if available.

If you can run `kubectl get all` on a shell and get a non-error response from your local cluster,
you're set!

Note that all Kubernetes files in this repo are configured to use the `events-demo` namespace.  You
may want to [configure kubectl](https://kubernetes.io/docs/tasks/access-application-cluster/configure-access-multiple-clusters/)
to use this namespace to make using other CLI things easier, but this is not strictly necessary.

### Helm

[Get Helm](https://helm.sh/docs/using_helm/) to help install things to Kubernetes in a more controlled way.  This is used in [the infra subdirectory](./infra).

### Telepresence

[Follow the instructions here](https://www.telepresence.io/reference/install) to install Telepresence.

## Getting Started

[Go to the infra subdirectory](./infra) to get your stack running!

