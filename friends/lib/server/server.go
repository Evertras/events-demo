package server

import (
	"context"
	"net/http"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

type Server interface {
	ListenAndServe() error
}

type server struct {
	httpServer *http.Server
}

func New(addr string) Server {
	router := http.NewServeMux()

	router.HandleFunc("/invite", inviteHandler())

	s := &server{
		httpServer: &http.Server{
			Addr:         addr,
			WriteTimeout: time.Second * 5,
			ReadTimeout:  time.Second * 5,
			Handler:      router,
		},
	}

	return s
}

func (s *server) ListenAndServe() error {
	return s.httpServer.ListenAndServe()
}

func startSpan(operationName string, r *http.Request) (opentracing.Span, context.Context) {
	spanCtx, _ := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
	span := opentracing.StartSpan(operationName, ext.RPCServerOption(spanCtx))

	span.SetTag("component", "server")

	ctx := opentracing.ContextWithSpan(r.Context(), span)

	return span, ctx
}
