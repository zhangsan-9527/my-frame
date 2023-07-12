//go:build e2e

package opentelemetry

import (
	"frame/web"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.15.0"
	"log"
	"os"
	"testing"
	"time"
)

func TestMiddlewareBuilder_Build(t *testing.T) {
	tracer := otel.GetTracerProvider().Tracer(instrumentationName)
	builder := MiddlewareBuilder{
		Tracer: tracer,
	}
	server := web.NewHTTPServer(web.ServerWithMiddleware(builder.Build()))

	server.Get("/user", func(ctx *web.Context) {
		c, span := tracer.Start(ctx.Req.Context(), "first_layer")
		defer span.End()

		secondC, second := tracer.Start(c, "second_layer")
		time.Sleep(time.Second)
		_, third1 := tracer.Start(secondC, "third_layer1")
		time.Sleep(100 * time.Millisecond)
		third1.End()
		_, third2 := tracer.Start(secondC, "third_layer2")
		time.Sleep(100 * time.Millisecond)
		third2.End()
		second.End()

		_, first := tracer.Start(ctx.Req.Context(), "first_layer1")
		defer first.End()

		ctx.RespJSON(202, User{
			Name: "zhangsan",
			Age:  18,
		})

	})

	initZipkin(t)

	server.Start(":8081")

}

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func initZipkin(t *testing.T) {
	exporyer, err := zipkin.New(
		"http://120.46.191.186:9411/api/v2/spans",
		zipkin.WithLogger(log.New(os.Stderr, "opentelemetry-demo", log.Ldate)),
	)
	if err != nil {
		t.Fatal(err)
	}
	batcher := sdktrace.NewBatchSpanProcessor(exporyer)
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
		// Always be syre to batch in production
		sdktrace.WithBatcher(exp),
		// Record information about this application in a Resource
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("opentelemetry-demo"),
			attribute.String("environment", "dev"),
			attribute.Int64("ID", 1),
		)),
	)
	otel.SetTracerProvider(tp)
}
