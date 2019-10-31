# Tech choice explanation

Here's a few explanations of various tech choices made.  It's important to note there are viable alternatives
for each choice, and none of this should be taken as required for the architecture as a whole to work.

## Kubernetes

https://kubernetes.io/

### Reasoning

The de-facto standard for easily managing large containerized deploys.  It allows us to do cool things like:

* Do rolling updates for services
* Budget resources on various services
* Auto-scale services based on resource usage
* Provide in-cluster load balancing for simple horizontal scaling
* Robust lifecycle management and auto-healing of deployments
* Strong ecosystem/community
* Allows for installation of entire tech stacks with [Helm](https://github.com/helm/helm)

### Alternatives

Docker Swarm is simpler to run, but doesn't have as many features as Kubernetes.  For example,
it lacks autoscaling and the concept of 'deployments', which would make us put in extra work
to have zero downtime rollouts.  However, it's worth consideration for its simplicity; Kubernetes
is not simple to run, and something like EKS can end up costing hundreds of dollars per month.

## Go

### Reasoning

Go is often called "the language of the cloud".  Go is used here for the following reasons:

* It can compile to extremely tiny docker images (<10 MB)
* It has excellent performance running as native code
* It has excellent support for multithreading/concurrency
* Strong type safety
* There's a strong ecosystem for cloud-based systems and clients

### Alternatives

Any language is honestly fine.  Note that while a strength of the architecture is that multiple languages
can be used, in practice a single language will probably win out simply because it will make it easier for
developers to move between pieces more easily.

Considerations for language selection:

* Good clients for other pieces in the architecture
* Good opentracing/metrics libraries for visibility
* Runs well in Docker/Kubernetes (C# may not be great here, for example)
* Developer familiarity

## Kafka

### Reasoning

Kafka is an append-only event log built specifically to scale to ridiculous size.  It allows us to record
events and read them in a distributed fashion via Consumer Groups, which lets us scale cleanly as well
for any service that wants to process events.

Unlike something like RabbitMQ, Kafka will retain all events.  This allows us to much more cleanly
decouple consumers and producers; producers should NOT know who consumers are, and consumers should
NOT know who producers are.  All either should care about is the creation and processing of events.

By default Kafka will drop any messages over 7 days old.  It can be configured to never drop any.  As a
future improvement, we will likely want to use [Log Compaction](https://medium.com/swlh/introduction-to-topic-log-compaction-in-apache-kafka-3e4d4afd2262).
For now, we'll pretend we have really big hard drives.  We do NOT want Kafka to lose data, or our other
services will not be able to rebuild their own data stores!

### Alternatives

[EventStore](https://eventstore.org/) is an interesting alternative but is less popular.

Any database could theoretically be used as an event store, but scaling becomes a concern.

Specifically, RabbitMQ is not considered an alternative here because it loses information as soon as
the message is delivered, and it requires more knowledge of producer/consumer patterns than we want to
be bothered with.

## Avro

https://avro.apache.org/

### Reasoning

Avro is a data serialization scheme similar to Protobuf or Thrift.  While Protobuf is a popular choice
for this sort of thing, Avro is used here because of Kafka and the potential for other Big Data tools
integration.  It seems that Avro is the winner for Big Data due to its semi-dynamic nature, so we will
use Avro with the intent of future proofing ourselves and integrating ourselves cleanly with the Kafka
ecosystem that uses it.

### Alternatives

Protobuf is an obvious alternative.  However, it lacks native integration with Kafka.  If you want to
implement RPCs it's worth considering alongside Avro, but Avro does seem to be a better fit for events
on Kafka.

## Taskfile

https://taskfile.dev

### Reasoning

Taskfile is a Go-based task runner.  We can define relatively simple, human-readable task definitions
and then run them.  Because we are already heavily invested in Go, the installation is very simple
and should work on any developer platform or CI tool.

### Alternatives

Make is a time-honored, battle-hardened task runner.  While it is ubiquitous in Linux, it can be a pain
to install on Windows compared to Taskfile and its syntax can be initially intimidating.

## Telepresence

https://www.telepresence.io

### Reasoning

Kubernetes is a great tool; however, its networking can be incredibly intricate and raises a concern:
does a developer have to deploy their application completely every time they want to test any change
at all, even locally?  If you want to test your application as it would be running in production, then
the answer is unfortunately yes.  Fortunately Telepresence fixes that for us.

Telepresence allows us to run a shell very quickly in Kubernetes which allows for easy development.
Each service should have a task to open a development shell that allows for direct running
via `go run cmd/thing/main.go`, or whatever you would normally run to start the service.  Your service
is now running with full access to other services exactly as it would in production, removing variables
and unexpected bugs that come from trying to develop in one environment while deploying to another.

The advantage here is that we can bring up our entire stack in Kubernetes, then surgically replace
a running instance of a service with our development environment.  This means that we can run our
code in a functionally production setting to minimize "worked on my machine" differences and not
have to deploy to staging environments just to test our changes.

### Alternatives

shrug.jpg

## Jaeger

https://www.jaegertracing.io/

### Reasoning

We want a tracing system in the first place because it's incredibly difficult to debug a distributed system
without knowing how data is actually flowing through it.  Tracing allows us to see the lifecycle of events
and requests.

Jaeger is an OpenTracing compliant tracing system written in Go.  It is part of the Cloud Native Computing
Foundation, which means it's well-supported by most cloud-y things.

### Alternatives

Zipkin is an older alternative.  While it's also popular, Jaeger is better integrated into cloud-based
systems and is built with Go in mind first and foremost.

Elasticsearch APM is an interesting alternative, but requires some hoops to jump through for OpenTracing
compliance and doesn't have as strong of a feature set as Jaeger when it comes to things like sampling.
Elasticsearch APM is a good fit for web servers in javascript, but doesn't seem to be as nice otherwise
just yet.
