package otel

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
	"go.opentelemetry.io/otel/propagation"
)

type OtelRoundTripper struct {
	ctx              context.Context
	wrappedTransport fasthttp.RoundTripper
}

var _ fasthttp.RoundTripper = (*OtelRoundTripper)(nil)

func NewOtelRoundTripper(ctx context.Context, wrappedTransport fasthttp.RoundTripper) *OtelRoundTripper {
	return &OtelRoundTripper{ctx: ctx, wrappedTransport: wrappedTransport}
}

func DistributedTraceSupport(ctx context.Context, agent *fiber.Agent) {
	if agent.Transport == nil {
		agent.Transport = fasthttp.DefaultTransport
	}

	agent.Transport = NewOtelRoundTripper(ctx, agent.Transport)
}

func (o *OtelRoundTripper) RoundTrip(client *fasthttp.HostClient, req *fasthttp.Request, resp *fasthttp.Response) (bool, error) {
	tc := propagation.TraceContext{}
	mc := propagation.MapCarrier{}

	tc.Inject(o.ctx, mc)
	req.Header.DisableNormalizing() // Disable header normalization to keep the injected headers as is
	// fmt.Printf("%v", mc)
	if _, ok := mc["traceparent"]; ok && mc["traceparent"] != "" {
		req.Header.Set("traceparent", mc.Get("traceparent"))
	}
	if _, ok := mc["tracestate"]; ok {
		req.Header.Set("tracestate", mc.Get("tracestate"))
	}

	return o.wrappedTransport.RoundTrip(client, req, resp)
}
