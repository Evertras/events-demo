module github.com/Evertras/events-demo/auth

go 1.13

require (
	github.com/Evertras/events-demo/shared/stream v0.0.0-00010101000000-000000000000
	github.com/actgardner/gogen-avro v6.2.0+incompatible
	github.com/badoux/checkmail v0.0.0-20181210160741-9661bd69e9ad
	github.com/bsm/redislock v0.4.0
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/go-redis/redis v6.15.5+incompatible
	github.com/google/uuid v1.1.1
	github.com/lib/pq v1.2.0
	github.com/opentracing/opentracing-go v1.1.0
	github.com/pkg/errors v0.8.1
	github.com/segmentio/kafka-go v0.3.4
	github.com/uber/jaeger-client-go v2.19.0+incompatible
	github.com/uber/jaeger-lib v2.2.0+incompatible
	golang.org/x/crypto v0.0.0-20191011191535-87dc89f01550
)

replace github.com/Evertras/events-demo/shared/stream => ./lib/shared/stream
