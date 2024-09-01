package otel

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var Tracer trace.Tracer

func NewTracer(cfg OtelConfig) {
	if !cfg.Trace.Enabled {
		Tracer = otel.Tracer("")
		return
	}

	Tracer = otel.Tracer(cfg.Trace.TracerName)
}
