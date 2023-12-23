package config

import (
	"context"
	"io"
	"os"

	gcpexporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"

	libdomain "github.com/kujilabo/redstart/lib/domain"
)

type OTLPConfig struct {
	Endpoint string `yaml:"endpoint" validate:"required"`
	Insecure bool   `yaml:"insecure"`
}

type TraceConfig struct {
	Exporter string      `yaml:"exporter" validate:"required"`
	OTLP     *OTLPConfig `yaml:"otlp"`
}

func initTracerExporter(ctx context.Context, traceConfig *TraceConfig) (sdktrace.SpanExporter, error) {
	switch traceConfig.Exporter {
	case "otlp":
		options := make([]otlptracehttp.Option, 0)
		options = append(options, otlptracehttp.WithEndpoint(traceConfig.OTLP.Endpoint))
		if traceConfig.OTLP.Insecure {
			options = append(options, otlptracehttp.WithInsecure())
		}
		return otlptracehttp.New(ctx, options...)
	case "gcp":
		projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
		return gcpexporter.New(gcpexporter.WithProjectID(projectID))
	case "stdout":
		return stdouttrace.New(
			stdouttrace.WithPrettyPrint(),
			stdouttrace.WithWriter(os.Stderr),
		)
	case "none":
		return stdouttrace.New(
			stdouttrace.WithPrettyPrint(),
			stdouttrace.WithWriter(io.Discard),
		)
	default:
		return nil, libdomain.ErrInvalidArgument
	}
}

func InitTracerProvider(ctx context.Context, appName string, traceConfig *TraceConfig) (*sdktrace.TracerProvider, error) {
	exp, err := initTracerExporter(ctx, traceConfig)
	if err != nil {
		return nil, err
	}
	tp := sdktrace.NewTracerProvider(
		// Always be sure to batch in production.
		sdktrace.WithBatcher(exp),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		// Record information about this application in a Resource.
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(appName),
		)),
	)

	return tp, nil
}
