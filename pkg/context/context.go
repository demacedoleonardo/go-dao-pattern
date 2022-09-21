package context

import (
	"context"
	"fmt"
	"net/http"

	"github.com/pedidosya/aragorn-service/pkg/metrics"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

const (
	appVersion = "App-Version"
	origin     = "Origin"
	tagStatus  = "status"
	tagReason  = "reason"
)

var (
	appName = "aragorn.service"
)

type (
	Tags struct { 
		tags map[string]string
	}

	Context struct {
		r    http.Request
		tags Tags
		spam ddtrace.SpanContext
		ctx  context.Context
	}
)

func (t *Tags) add(key, value string) *Tags {
	if t.tags == nil {
		t.tags = make(map[string]string)
	}

	t.tags[key] = value
	return t
}

func (t *Tags) ToArray() []string {
	tags := make([]string, 0)
	for k, v := range t.tags {
		tags = append(tags, fmt.Sprintf("%s:%s", k, v))
	}
	return tags
}

func (c *Context) WithTags(k string, v interface{}) {
	if v != nil {
		c.tags.add(k, fmt.Sprint(v))
	}
}

func (c *Context) Headers() http.Header {
	return c.r.Header
}

func (c *Context) WithError(err error) {
	tracer.WithError(err)
}

func (c *Context) ParentSpam() ddtrace.SpanContext {
	return c.spam
}

func (c *Context) Context() context.Context {
	return c.ctx
}

func (c *Context) SetErrReason(reason string) *Context {
	c.tags.add(tagReason, reason)
	return c
}

func (c *Context) Send(flow string) {
	if c == nil {
		return
	}

	c.tags.add("flow", flow)
	metric := fmt.Sprintf("%s.stats", appName)
	metrics.IncrementCounter(metric, 1, c.tags.ToArray()...)
}

func NewContext() *Context {
	return &Context{
		r:    http.Request{},
		tags: Tags{},
		ctx:  context.Background(),
	}
}

func NewWebContext(r *http.Request) *Context {
	sctx, err := tracer.Extract(tracer.HTTPHeadersCarrier(r.Header))
	if err != nil {
		tracer.WithError(err)
	}

	return &Context{
		r:    *r.Clone(r.Context()),
		ctx:  r.Context(),
		spam: sctx,
	}
}

func NewBackgroundContext() *Context {
	c := &Context{
		tags: Tags{},
	}

	c.ctx = context.Background()
	return c
}
