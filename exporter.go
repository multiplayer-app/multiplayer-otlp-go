package pkg

import (
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"

	otlptracehttp "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
)

func NewExporter(mpKey string) *otlptrace.Exporter {
	headers := map[string]string{
		"Authorization": mpKey,
	}

	client := otlptracehttp.NewClient(
		otlptracehttp.WithEndpointURL(MULTIPLAYER_OTEL_DEFAULT_TRACES_EXPORTER_URL),
		otlptracehttp.WithHeaders(headers))
	exporter := otlptrace.NewUnstarted(client)
	return exporter
}
