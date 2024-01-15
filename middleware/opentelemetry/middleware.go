package opentelemetry

import (
	"github.com/NotFound1911/mserver"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

const instrumentationName = "github.com/NotFound1911/mserver/middlewaress/opentelemetry"

type MiddlewareBuilder struct {
	Tracer trace.Tracer
}

func (m MiddlewareBuilder) Build() mserver.Middleware {
	if m.Tracer == nil {
		m.Tracer = otel.GetTracerProvider().Tracer(instrumentationName)
	}
	return func(next mserver.HandleFunc) mserver.HandleFunc {
		return func(ctx *mserver.Context) error {
			reqCtx := ctx.GetRequest().Context()
			// 尝试和客户端的 trace 结合在一起
			reqCtx = otel.GetTextMapPropagator().Extract(reqCtx, propagation.HeaderCarrier(ctx.GetRequest().Header))

			reqCtx, span := m.Tracer.Start(reqCtx, "unknown")
			defer span.End()

			span.SetAttributes(attribute.String("http.method", ctx.GetRequest().Method))
			span.SetAttributes(attribute.String("http.url", ctx.GetRequest().URL.String()))
			span.SetAttributes(attribute.String("http.scheme", ctx.GetRequest().URL.Scheme))
			span.SetAttributes(attribute.String("http.host", ctx.GetRequest().Host))
			// 还可以继续添加

			ctx.SetRequest(ctx.GetRequest().WithContext(reqCtx))
			next(ctx)

			// 执行完next 才可能有值
			span.SetName(ctx.MatchedRoute)

			// 把响应加上去
			span.SetAttributes(attribute.Int("http.status", ctx.GetRespStatusCode()))
			return nil
		}
	}
}
