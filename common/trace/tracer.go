package trace

import (
	"context"
	"skytakeout/global"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

func InitTracer() func(context.Context) error {
	exp, err := otlptracehttp.New(context.Background(),
		otlptracehttp.WithEndpoint(global.Config.Jaeger.EndPoint),
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		panic(err)
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(global.Config.Jaeger.ServiceName),
		)),
	)
	otel.SetTracerProvider(tp)
	return tp.Shutdown
}
