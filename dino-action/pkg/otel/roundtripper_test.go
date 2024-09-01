package otel

import (
	"context"
	"fmt"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
)

func TestRoundtripper(t *testing.T) {
	app := fiber.New(fiber.Config{DisableHeaderNormalizing: true})
	app.Get("/test", func(c *fiber.Ctx) error {
		req := c.Request()
		fmt.Printf("%v", c.Request())
		tr := req.Header.Peek("traceparent")
		return c.SendString(string(tr))
	})
	go app.Listen(":3000")

	ctx, span := Tracer.Start(context.Background(), "test")
	defer span.End()

	rt := NewOtelRoundTripper(ctx, fasthttp.DefaultTransport)
	agent := fiber.Get("http://localhost:3000/test")
	agent.Transport = rt

	_, s, errs := agent.String()
	if len(errs) > 0 {
		t.Fatal(errs)
	}

	assert.NotEmpty(t, s)
	assert.NoError(t, app.Shutdown())
}
