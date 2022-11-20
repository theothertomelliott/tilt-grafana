package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
)

func main() {
	l := log.New(os.Stdout, "", 0)

	exp, err := newExporter(context.Background())
	if err != nil {
		l.Fatal(err)
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exp),
		trace.WithResource(newResource()),
	)
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			l.Fatal(err)
		}
	}()
	otel.SetTracerProvider(tp)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	errCh := make(chan error)

	// Randomly generate input numbers
	reader, writer := io.Pipe()
	go func() {
		for true {
			writer.Write([]byte(fmt.Sprintf("%v\n", rand.Intn(100))))
			time.Sleep(time.Second)
		}
	}()

	app := NewApp(reader, l)
	go func() {
		errCh <- app.Run(context.Background())
	}()

	select {
	case <-sigCh:
		l.Println("\ngoodbye")
		return
	case err := <-errCh:
		if err != nil {
			l.Fatal(err)
		}
	}
}

const (
	ENV_OTEL_JAEGER_ENDPOINT    = "OTEL_JAEGER_ENDPOINT"
	ENV_OTEL_OTLP_HTTP_ENDPOINT = "OTEL_OTLP_HTTP_ENDPOINT"
)

// newExporter returns an opentelemetry exporter configured according to environment variables
func newExporter(ctx context.Context) (trace.SpanExporter, error) {
	if os.Getenv(ENV_OTEL_OTLP_HTTP_ENDPOINT) != "" {
		url := os.Getenv(ENV_OTEL_OTLP_HTTP_ENDPOINT)
		client := otlptracehttp.NewClient(
			otlptracehttp.WithEndpoint(url),
			// TODO: Make this optional
			otlptracehttp.WithInsecure(),
		)
		return otlptrace.New(ctx, client)
	}

	if os.Getenv(ENV_OTEL_JAEGER_ENDPOINT) != "" {
		url := os.Getenv(ENV_OTEL_JAEGER_ENDPOINT)
		return jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	}

	// Default to stdout logging
	w := os.Stdout
	return stdouttrace.New(
		stdouttrace.WithWriter(w),
		// Use human-readable output.
		stdouttrace.WithPrettyPrint(),
	)
}

// newResource returns a resource describing this application.
func newResource() *resource.Resource {
	r, _ := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("fib"),
			semconv.ServiceVersionKey.String("v0.1.0"),
			attribute.String("environment", "demo"),
		),
	)
	return r
}
