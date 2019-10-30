package server

import (
	"net/http"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/Evertras/events-demo/auth/lib/auth"
	"github.com/Evertras/events-demo/auth/lib/tracing"
)

type Server interface {
	ListenAndServe() error
}

type server struct {
	tracer     opentracing.Tracer
	httpServer *http.Server
}

func New(addr string, auth auth.Auth) (Server, error) {
	tracer, err := tracing.Init("http")

	if err != nil {
		return nil, errors.Wrap(err, "failed to init tracer")
	}

	router := http.NewServeMux()

	s := &server{
		httpServer: &http.Server{
			Addr:         addr,
			WriteTimeout: time.Second * 5,
			ReadTimeout:  time.Second * 5,
			Handler:      router,
		},
		tracer: tracer,
	}

	router.HandleFunc("/check", checkAuthHandler(tracer))
	router.HandleFunc("/login", loginHandler(tracer, auth))
	router.HandleFunc("/register", registerHandler(tracer, auth))

	return s, nil
}

func (s *server) ListenAndServe() error {
	return s.httpServer.ListenAndServe()
}
