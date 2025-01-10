package main

import (
	"context"
	"log/slog"
	"os"

	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/sdk/log"
)

func main() {
	ctx := context.Background()
	insecure := os.Getenv("OTEL_EXPORTER_OTLP_INSECURE") == "true"
	useHttp := os.Getenv("OTEL_EXPORTER_OTLP_PROTOCOL") == "http"

	var processor *log.BatchProcessor

	if useHttp {
		options := []otlploghttp.Option{}

		if insecure {
			options = append(options, otlploghttp.WithInsecure())
		}

		exp, err := otlploghttp.New(ctx, options...)
		if err != nil {
			panic(err)
		}

		processor = log.NewBatchProcessor(exp)
	} else {
		options := []otlploggrpc.Option{}

		if insecure {
			options = append(options, otlploggrpc.WithInsecure())
		}

		exp, err := otlploggrpc.New(ctx, options...)
		if err != nil {
			panic(err)
		}

		processor = log.NewBatchProcessor(exp)
	}
	provider := log.NewLoggerProvider(log.WithProcessor(processor))
	handler := otelslog.NewHandler("otelslog", otelslog.WithLoggerProvider(provider))
	slogger := slog.New(handler)
	slogger.Info("Hello, World!")

	select {}
}
