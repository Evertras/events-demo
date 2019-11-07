FROM golang:1.13 AS go-builder

RUN apt-get update
RUN apt-get install libssl-dev pkg-config
RUN curl -Lo /tmp/seabolt.tar.gz https://github.com/neo4j-drivers/seabolt/releases/download/v1.7.4/seabolt-1.7.4-Linux-ubuntu-18.04.tar.gz
RUN tar zxvf /tmp/seabolt.tar.gz --strip-components=1 -C /

WORKDIR /go/src/github.com/Evertras/events-demo
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
      go build -a -tags netgo \
        --tags seabolt_static \
        -ldflags '-w -extldflags "-static"' \
        -o /server \
        ./cmd/server/main.go

FROM scratch

COPY --from=go-builder /server server

ENTRYPOINT ["/server"]

