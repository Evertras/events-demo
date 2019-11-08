module github.com/Evertras/events-demo/friends

go 1.13

replace github.com/Evertras/events-demo/shared/stream => ./lib/shared/stream

require (
	github.com/Evertras/events-demo/shared/stream v0.0.0-00010101000000-000000000000
	github.com/actgardner/gogen-avro v6.2.0+incompatible
	github.com/google/uuid v1.1.1
	github.com/neo4j-drivers/gobolt v1.7.4 // indirect
	github.com/neo4j/neo4j-go-driver v1.7.4
	github.com/opentracing/opentracing-go v1.1.0
	github.com/pkg/errors v0.8.1
	github.com/uber/jaeger-client-go v2.20.0+incompatible
	github.com/uber/jaeger-lib v2.2.0+incompatible
	go.uber.org/atomic v1.5.0 // indirect
)
