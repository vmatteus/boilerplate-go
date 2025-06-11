package telemetry

import (
	"context"
	"log"
	"os"

	"github.com/your-org/boilerplate-go/internal/config"
	"go.opentelemetry.io/contrib/instrumentation/host"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	traceOtel "go.opentelemetry.io/otel/trace"
)

func InitTelemetry(ctx context.Context, appConfig *config.Config) func() {

	os.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", appConfig.Apm.Url)
	os.Setenv("OTEL_RESOURCE_ATTRIBUTES", appConfig.Apm.Attributes)
	os.Setenv("OTEL_EXPORTER_OTLP_HEADERS", appConfig.Apm.Headers)

	res, err := resource.New(
		ctx,
		resource.WithProcess(),
		resource.WithOS(),
		resource.WithTelemetrySDK(),
		resource.WithContainer(),
		resource.WithHost(),
		resource.WithSchemaURL(semconv.SchemaURL),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(appConfig.Application.Name),
			semconv.ServiceVersionKey.String("1.0"),
			semconv.DeploymentEnvironmentKey.String(appConfig.Application.Environment),
		),
	)
	if err != nil {
		log.Fatalf("OpenTelemetry resource failed: %v", err)
		return func() {}
	}

	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	traceExporter, err := otlptracehttp.New(ctx)
	if err != nil {
		log.Fatalf("OpenTelemetry trace exporter failed: %v", err)
		return func() {}
	}

	tp := trace.NewTracerProvider(
		trace.WithSyncer(traceExporter),
		trace.WithResource(res),
	)
	otel.SetTracerProvider(tp)

	metricExporter, err := otlpmetrichttp.New(ctx)
	if err != nil {
		log.Fatalf("OpenTelemetry metric exporter failed: %v", err)
		return func() {}
	}
	meterProvider := metric.NewMeterProvider(
		metric.WithReader(metric.NewPeriodicReader(metricExporter)),
		metric.WithResource(res),
	)
	otel.SetMeterProvider(meterProvider)

	if err := runtime.Start(runtime.WithMinimumReadMemStatsInterval(1)); err != nil {
		log.Printf("Erro ao iniciar métricas runtime: %v", err)
	}

	if err := host.Start(); err != nil {
		log.Printf("Erro ao iniciar métricas do host: %v", err)
	}

	return func() {}
}

func GetTracer() traceOtel.Tracer {
	return otel.Tracer("banking-router")
}
