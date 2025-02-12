package otlp

import (
	"context"

	otelpkg "go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type Span interface {
	trace.Span
	EndError(err error, options ...trace.SpanEndOption)
	Error(err error)
}

func Start(ctx context.Context, name, spanName string) (context.Context, Span) {
	ctx, _span := otelpkg.Tracer(name).Start(ctx, spanName)
	return ctx, &span{Span: _span}
}

type span struct {
	Span trace.Span
}



// AddLink implements Span.
func (s *span) AddLink(link trace.Link) {
	panic("unimplemented")
}


func (s *span) End(options ...trace.SpanEndOption) {
	s.Span.End(options...)
}

func (s *span) EndError(err error, options ...trace.SpanEndOption) {
	s.Error(err)
	s.Span.End(options...)
}

func (s *span) AddEvent(name string, options ...trace.EventOption) {
	s.Span.AddEvent(name, options...)
}

func (s *span) IsRecording() bool {
	return s.Span.IsRecording()
}

func (s *span) RecordError(err error, options ...trace.EventOption) {
	s.Span.RecordError(err, options...)
}

func (s *span) SpanContext() trace.SpanContext {
	return s.Span.SpanContext()
}

func (s *span) SetStatus(code codes.Code, description string) {
	s.Span.SetStatus(code, description)
}

func (s *span) SetName(name string) {
	s.Span.SetName(name)
}

func (s *span) SetAttributes(kv ...attribute.KeyValue) {
	s.Span.SetAttributes(kv...)
}

func (s *span) TracerProvider() trace.TracerProvider {
	return s.Span.TracerProvider()
}

func (s *span) Error(err error) {
	if err != nil {
		s.Span.SetStatus(codes.Error, err.Error())
	}
}

// RestoreTraceContext function forms context and span from trace_id and span_id
//
// span_id and trace_id should both be strings in hex format.
//
// Although this function returns both context and span it is required to call Start method to start tracing
// WARNING: if error IS NOT NIL, then context and span are BOTH NIL.
func RestoreTraceContext(traceIdStr, spanIdStr string) (context.Context, trace.Span, error) {
	spanId, err := trace.SpanIDFromHex(spanIdStr)
	if err != nil {
		return nil, nil, err
	}

	traceId, err := trace.TraceIDFromHex(traceIdStr)
	if err != nil {
		return nil, nil, err
	}

	ctx := trace.ContextWithRemoteSpanContext(context.Background(), trace.NewSpanContext(trace.SpanContextConfig{
		TraceID: traceId,
		SpanID:  spanId,
	}))

	return ctx, trace.SpanFromContext(ctx), nil
}
