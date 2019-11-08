module github.com/Evertras/events-demo/friends

go 1.13

require (
	github.com/Evertras/events-demo/shared/stream v0.0.0-00010101000000-000000000000
	github.com/actgardner/gogen-avro v6.2.0+incompatible
	github.com/codahale/hdrhistogram v0.0.0-20161010025455-3a0bb77429bd // indirect
	github.com/golang/mock v1.3.1 // indirect
	github.com/google/uuid v1.1.1
	github.com/neo4j-drivers/gobolt v1.7.4 // indirect
	github.com/neo4j/neo4j-go-driver v1.7.4
	github.com/onsi/ginkgo v1.10.3 // indirect
	github.com/onsi/gomega v1.7.1 // indirect
	github.com/opentracing/opentracing-go v1.1.0
	github.com/pkg/errors v0.8.1
	github.com/uber/jaeger-client-go v2.19.0+incompatible
	github.com/uber/jaeger-lib v2.2.0+incompatible
)

replace github.com/Evertras/events-demo/shared/stream => ./lib/shared/stream
