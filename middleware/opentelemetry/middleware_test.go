package opentelemetry

import (
	"github.com/NotFound1911/mserver"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"log"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestMiddlewareBuilder_Build(t *testing.T) {
	tracer := otel.GetTracerProvider().Tracer(instrumentationName)
	builder := MiddlewareBuilder{
		Tracer: tracer,
	}
	core := mserver.NewCore()
	core.Get("/test/opentelmetry", func(ctx *mserver.Context) error {
		c, span := tracer.Start(ctx.GetRequest().Context(), "first_layer")
		defer span.End()

		secondC, second := tracer.Start(c, "second_layer")
		time.Sleep(time.Second)

		_, third1 := tracer.Start(secondC, "third_layer_1")
		time.Sleep(100 * time.Millisecond)
		third1.End()

		_, third2 := tracer.Start(secondC, "third_layer_2")
		time.Sleep(300 * time.Millisecond)
		third2.End()
		second.End()

		_, first := tracer.Start(ctx.GetRequest().Context(), "first_layer_1")
		defer first.End()
		ctx.SetStatus(http.StatusAccepted).Json(User{
			Name: "test",
		})
		return nil
	})
	core.Use(builder.Build())

	initZipkin(t)

	core.Start(":8081")
}

type User struct {
	Name string
}

func initZipkin(t *testing.T) {
	// 要注意这个端口，和 docker-compose 中的保持一致
	exporter, err := zipkin.New(
		"http://localhost:19411/api/v2/spans",
		zipkin.WithLogger(log.New(os.Stderr, "opentelemetry-demo", log.Ldate|log.Ltime|log.Llongfile)),
	)
	if err != nil {
		t.Fatal(err)
	}

	batcher := sdktrace.NewBatchSpanProcessor(exporter)
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSpanProcessor(batcher),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("opentelemetry-demo"),
		)),
	)
	otel.SetTracerProvider(tp)
}

func initJeager(t *testing.T) {
	url := "http://localhost:14268/api/traces"
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		t.Fatal(err)
	}
	tp := sdktrace.NewTracerProvider(
		// Always be sure to batch in production.
		sdktrace.WithBatcher(exp),
		// Record information about this application in a Resource.
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("opentelemetry-demo"),
			attribute.String("environment", "dev"),
			attribute.Int64("ID", 1),
		)),
	)

	otel.SetTracerProvider(tp)
}
