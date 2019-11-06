package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/pkg/errors"
	jaegerconfig "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	span, _ := startSpan("hello", r)
	defer span.Finish()

	w.Write([]byte("Hello\n"))
	w.Write([]byte(r.URL.String() + "\n"))
	w.Write([]byte(r.URL.Hostname() + "\n"))
	w.Write([]byte(r.URL.RawPath + "\n"))

	log.Println("helloHandler end")
}

func headerEchoHandler(w http.ResponseWriter, r *http.Request) {
	span, _ := startSpan("header echo", r)
	defer span.Finish()

	for header := range r.Header {
		w.Write([]byte(header + ": " + r.Header.Get(header) + "\n"))
	}

	log.Println("headerEchoHandler end")
}

func main() {
	addr := "0.0.0.0:13030"

	if err := initTracing(); err != nil {
		log.Fatal(err)
	}

	router := http.NewServeMux()

	router.HandleFunc("/", helloHandler)
	router.HandleFunc("/headers", headerEchoHandler)

	server := &http.Server{
		Addr:         addr,
		WriteTimeout: time.Second * 5,
		ReadTimeout:  time.Second * 5,
		Handler:      router,
	}

	log.Println("Serving", addr)

	log.Fatal(server.ListenAndServe())
}

func initTracing() error {
	cfg, err := jaegerconfig.FromEnv()

	if err != nil {
		return errors.Wrap(err, "failed to create tracer config")
	}

	cfg.ServiceName = "sample"

	tracer, _, err := cfg.NewTracer(
		jaegerconfig.Logger(jaegerlog.StdLogger),
		jaegerconfig.Metrics(metrics.NullFactory),
	)

	if err != nil {
		return errors.Wrap(err, "failed to create tracer")
	}

	opentracing.SetGlobalTracer(tracer)

	return nil
}

func startSpan(operationName string, r *http.Request) (opentracing.Span, context.Context) {
	spanCtx, _ := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
	span := opentracing.StartSpan(operationName, ext.RPCServerOption(spanCtx))

	span.SetTag("component", "server")

	ctx := opentracing.ContextWithSpan(r.Context(), span)

	return span, ctx
}
