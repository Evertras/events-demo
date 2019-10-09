# Infra

Contains core services and some tasks to bring them up or down locally.

## Installing Local Infra

Check [Taskfile.yml](./Taskfile.yml) for the full commands.

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
# Remove our Traefik ingress
task delete-local-traefik

# Remove our Kafka stack
task delete-local-kafka
```

The order does not matter, and you should feel free to bring things up and down independently at will.

## Kafka connection information

For any clients that need to know a Kafka broker, supply the following:

Kafka brokers: `events-demo-kafka-cp-kafka-headless:9092`
Zookeeper: `events-demo-kafka-cp-zookeeper-headless:2181`

