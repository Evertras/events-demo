package tracing

import (
	"github.com/pkg/errors"

	"github.com/opentracing/opentracing-go"
	jaegerconfig "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"
)

func Init(serviceName string) (opentracing.Tracer, error) {
	cfg, err := jaegerconfig.FromEnv()

	if err != nil {
		return nil, errors.Wrap(err, "failed to create tracer config")
	}

	cfg.ServiceName = "auth-" + serviceName

	tracer, _, err := cfg.NewTracer(
		jaegerconfig.Logger(jaegerlog.StdLogger),
		jaegerconfig.Metrics(metrics.NullFactory),
	)

	if err != nil {
		return nil, errors.Wrap(err, "failed to create tracer")
	}

	return tracer, nil
}

/*
// Quick ref
func initTracing() (io.Closer, error) {
	samplerCfg := &jaegercfg.SamplerConfig{
		Type:              jaeger.SamplerTypeConst,
		Param:             1,
		SamplingServerURL: "http://jaeger:5778/sampling",
	}

	reporterCfg := &jaegercfg.ReporterConfig{
		LogSpans: true,
		LocalAgentHostPort: "jaeger:6831",
	}

	// Sample everything... don't use in production!
	cfg := jaegercfg.Configuration{
		Sampler: samplerCfg,
		Reporter: reporterCfg,
	}

	logger := jaegerlog.StdLogger

	closer, err := cfg.InitGlobalTracer(
		"auth-api",
		jaegercfg.Logger(logger),
		jaegercfg.Metrics(metrics.NullFactory),
	)

	return closer, err
}
*/
