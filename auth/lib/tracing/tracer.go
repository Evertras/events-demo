package tracing

import (
	"github.com/pkg/errors"

	"github.com/opentracing/opentracing-go"
	jaegerconfig "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"
)

func Init(serviceName string) error {
	cfg, err := jaegerconfig.FromEnv()

	if err != nil {
		return errors.Wrap(err, "failed to create tracer config")
	}

	cfg.ServiceName = serviceName

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
