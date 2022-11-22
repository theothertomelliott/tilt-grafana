package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"strconv"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// name is the Tracer name used to identify this instrumentation library.
const name = "fib"

// App is a Fibonacci computation application.
type App struct {
	r io.Reader
	s *zap.SugaredLogger
}

// NewApp returns a new App.
func NewApp(r io.Reader, l *log.Logger) *App {
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()
	return &App{r: r, s: sugar}
}

// Run starts polling users for Fibonacci number requests and writes results.
func (a *App) Run(ctx context.Context) error {
	for {
		// Each execution of the run loop, we should get a new "root" span and context.
		newCtx, span := otel.Tracer(name).Start(ctx, "Run")

		n, err := a.Poll(newCtx)
		if err != nil {
			span.End()
			return err
		}

		a.Write(newCtx, n)
		span.End()
	}
}

// Poll asks a user for input and returns the request.
func (a *App) Poll(ctx context.Context) (uint, error) {
	_, span := otel.Tracer(name).Start(ctx, "Poll")
	defer span.End()

	a.s.Infow(
		"What Fibonacci number would you like to know: ",
		"traceID", span.SpanContext().TraceID(),
	)

	var n uint
	_, err := fmt.Fscanf(a.r, "%d\n", &n)

	// Store n as a string to not overflow an int64.
	nStr := strconv.FormatUint(uint64(n), 10)
	span.SetAttributes(attribute.String("request.n", nStr))

	return n, err
}

// Write writes the n-th Fibonacci number back to the user.
func (a *App) Write(ctx context.Context, n uint) {
	var span trace.Span
	ctx, span = otel.Tracer(name).Start(ctx, "Write")
	defer span.End()

	f, err := func(ctx context.Context) (uint64, error) {
		_, span := otel.Tracer(name).Start(ctx, "Fibonacci")
		defer span.End()
		return Fibonacci(n)
	}(ctx)
	if err != nil {
		a.s.Infow(
			fmt.Sprintf("Fibonacci(%d): %v\n", n, err),
			"traceID", span.SpanContext().TraceID(),
		)
	} else {
		a.s.Infow(
			fmt.Sprintf("Fibonacci(%d) = %d\n", n, f),
			"traceID", span.SpanContext().TraceID(),
		)
	}
}
