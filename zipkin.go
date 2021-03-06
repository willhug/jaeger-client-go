package jaeger

import (
	"github.com/opentracing/opentracing-go"
)

// ZipkinSpanFormat is an OpenTracing carrier format constant
const ZipkinSpanFormat = "zipkin-span-format"

// ExtractableZipkinSpan is a type of Carrier used for integration with Zipkin-aware
// RPC frameworks (like TChannel). It does not support baggage, only trace IDs.
type ExtractableZipkinSpan interface {
	TraceID() uint64
	SpanID() uint64
	ParentID() uint64
	Flags() byte
}

// InjectableZipkinSpan is a type of Carrier used for integration with Zipkin-aware
// RPC frameworks (like TChannel). It does not support baggage, only trace IDs.
type InjectableZipkinSpan interface {
	SetTraceID(traceID uint64)
	SetSpanID(spanID uint64)
	SetParentID(parentID uint64)
	SetFlags(flags byte)
}

type zipkinPropagator struct {
	tracer *tracer
}

func (p *zipkinPropagator) Inject(
	ctx SpanContext,
	abstractCarrier interface{},
) error {
	carrier, ok := abstractCarrier.(InjectableZipkinSpan)
	if !ok {
		return opentracing.ErrInvalidCarrier
	}

	carrier.SetTraceID(ctx.TraceID())
	carrier.SetSpanID(ctx.SpanID())
	carrier.SetParentID(ctx.ParentID())
	carrier.SetFlags(ctx.flags)
	return nil
}

func (p *zipkinPropagator) Extract(abstractCarrier interface{}) (SpanContext, error) {
	carrier, ok := abstractCarrier.(ExtractableZipkinSpan)
	if !ok {
		return emptyContext, opentracing.ErrInvalidCarrier
	}
	if carrier.TraceID() == 0 {
		return emptyContext, opentracing.ErrSpanContextNotFound
	}
	var ctx SpanContext
	ctx.traceID = carrier.TraceID()
	ctx.spanID = carrier.SpanID()
	ctx.parentID = carrier.ParentID()
	ctx.flags = carrier.Flags()
	return ctx, nil
}
