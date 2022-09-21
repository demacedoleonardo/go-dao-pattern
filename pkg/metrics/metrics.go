package metrics

import (
	"context"
	"fmt"
	"net/http"

	"go-dao-pattern/pkg/env"

	"github.com/DataDog/datadog-go/statsd"
	log "github.com/pedidosya/peya-go/logs"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/ext"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

const (
	nameSpace = ""

	Http Platform = iota
	Cache
	Database
)

var (
	instance = NewCollector(env.Environment())
)

type (
	Platform int

	segment struct {
		action   string
		resource string
		platform Platform
		ctx      context.Context
		parent   ddtrace.SpanContext
		h        http.Header
	}

	SegmentOption func(s *segment)
)

func (p Platform) String() string {
	switch p {
	case Cache:
		return "cache"
	case Database:
		return "db"
	default:
		return "web"
	}
}

func NewCollector(env string) statsd.ClientInterface {
	agentAddress := fmt.Sprintf("datadog-agent.%s.peja.co:8125", env)
	envTag := fmt.Sprintf("env:%s", env)
	client, err := statsd.New(agentAddress, statsd.WithNamespace(nameSpace),
		statsd.WithTags([]string{envTag}))

	if err != nil {
		return nil
	}

	return client
}

func IncrementCounter(metricName string, value int64, tags ...string) {
	if err := instance.Count(metricName, value, tags, 1); err != nil {
		log.Error("[IncrementCounter] fail sending metrics", err)
	}
}

func StartSegment(f func(), opts ...SegmentOption) ddtrace.Span {
	s := new(segment)
	for _, opt := range opts {
		opt(s)
	}

	span, _ := tracer.StartSpanFromContext(s.ctx, s.action, tracer.ResourceName(s.resource),
		tracer.Measured(), tracer.SpanType(s.platform.String()), tracer.ChildOf(s.parent))

	if Http == s.platform && s.h != nil {
		if err := tracer.Inject(span.Context(), tracer.HTTPHeadersCarrier(s.h)); err != nil {
			span.Finish(tracer.WithError(err))
		}
	}

	f()

	return span
}

func StartStoreSegment(f func() error, opts ...SegmentOption) {
	s := new(segment)
	for _, opt := range opts {
		opt(s)
	}

	span, _ := tracer.StartSpanFromContext(s.ctx, s.action, tracer.ResourceName(s.resource),
		tracer.Measured(), tracer.SpanType(ext.AppTypeDB), tracer.ChildOf(s.parent))

	span.Finish(tracer.WithError(f()))
}

func WithAction(a string) SegmentOption {
	return func(s *segment) {
		s.action = a
	}
}

func WithResource(r string) SegmentOption {
	return func(s *segment) {
		s.resource = r
	}
}

func WithPlatform(p Platform) SegmentOption {
	return func(s *segment) {
		s.platform = p
	}
}

func WithContext(c context.Context) SegmentOption {
	return func(s *segment) {
		s.ctx = c
	}
}

func WithParent(c context.Context, parent ddtrace.SpanContext, h http.Header) SegmentOption {
	return func(s *segment) {
		s.ctx = c
		s.parent = parent
		s.h = h
	}
}
