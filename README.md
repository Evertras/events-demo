# events-demo

AKA Enterprise Tic Tac Toe

*Note*: Since this is for proof of concept, everything is crammed into a single repo.  In a real world
scenario each subdirectory would almost certainly be its own repository.

## Prerequisites

Because Enterprise software is Enterprisey, there are a few things you'll need...

### Go

Get yourself some [Go](https://golang.org/doc/install), 1.13 or higher.

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

[Go to the infra subdirectory](./infra) to get your stack running!

