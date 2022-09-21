package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	oops "go-dao-pattern/pkg/errors"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	apictx "go-dao-pattern/pkg/context"
	"go-dao-pattern/pkg/metrics"

	"github.com/pedidosya/go-client/kit/restclient/restclient"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/ext"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

type Actions interface {
	Get(uri string) (*Response, error)
	Post(uri string) (*Response, error)
	Put(endpoint string) (*Response, error)
}

type Body []byte

func (b Body) ToString() string {
	return string(b)
}

type Response struct {
	StatusCode int
	Body       Body
}

type Request struct {
	c        restclient.Client
	ctx      *apictx.Context
	body     io.Reader
	headers  Headers
	resource string
}

func (r *Request) Get(endpoint string) (*Response, error) {
	return execute(r.ctx, r.c, http.MethodGet, endpoint, r.resource, r.body, r.headers)
}

func (r *Request) Post(endpoint string) (*Response, error) {
	return execute(r.ctx, r.c, http.MethodPost, endpoint, r.resource, r.body, r.headers)
}

func (r *Request) Put(endpoint string) (*Response, error) {
	return execute(r.ctx, r.c, http.MethodPut, endpoint, r.resource, r.body, r.headers)
}

func NewRequest(options ...func(request *Request)) *Request {
	defaultHeader := make(map[string]string)
	defaultHeader["Content-Type"] = "application/json"

	r := &Request{
		ctx:     apictx.NewBackgroundContext(context.Background()),
		headers: Headers{headers: defaultHeader},
	}
	for _, option := range options {
		option(r)
	}
	return r
}

func WithClient(c restclient.Client) func(*Request) {
	return func(request *Request) {
		request.c = c
	}
}

func WithBody(body interface{}) func(*Request) {
	if b, ok := body.(string); ok {
		return func(request *Request) {
			request.body = strings.NewReader(b)
		}
	}

	if b, ok := body.([]byte); ok {
		return func(request *Request) {
			request.body = bytes.NewReader(b)
		}
	}

	return func(request *Request) {
		b, _ := json.Marshal(body)
		request.body = bytes.NewReader(b)
	}
}

func WithHeaders(h Headers) func(*Request) {
	return func(request *Request) {
		request.headers = h
	}
}

func WithResource(s string) func(*Request) {
	return func(request *Request) {
		request.resource = s
	}
}

func WithContext(c *apictx.Context) func(*Request) {
	return func(request *Request) {
		request.ctx = c
	}
}

type Headers struct {
	headers map[string]string
}

func (h *Headers) Add(k, v string) *Headers {
	if h.headers == nil {
		h.headers = make(map[string]string)
		h.headers["Content-Type"] = "application/json"
	}
	h.headers[k] = v
	return h
}

func BindResponse(response *Response, err error, in interface{}) error {
	if err != nil {
		return err
	}

	if in == nil {
		return nil
	}

	if err = json.Unmarshal(response.Body, &in); err != nil {
		return oops.Errorf(oops.E5xxINTERNAL, err.Error())
	}

	return nil
}

func toString(body io.ReadCloser) string {
	r, err := ioutil.ReadAll(body)
	if err != nil {
		return ""
	}
	return string(r)
}

func execute(ctx *apictx.Context, client restclient.Client, method, endpoint, resource string, payload io.Reader, h Headers) (*Response, error) {
	req, err := http.NewRequestWithContext(ctx.Context(), method, endpoint, payload)
	if err != nil {
		return nil, oops.Errorf(oops.E5xxINTERNAL, "error building request")
	}

	for k, v := range h.headers {
		req.Header.Add(k, v)
	}

	response := new(http.Response)
	segment := metrics.StartSegment(func() {
		response, err = client.Do(req)
	},
		metrics.WithAction(method),
		metrics.WithResource(resource),
		metrics.WithPlatform(metrics.Http),
		metrics.WithParent(ctx.Context(), ctx.ParentSpam(), ctx.Headers()),
	)

	defer func() {
		if err != nil {
			segment.SetTag(ext.ErrorMsg, err.Error())
		}
		segment.Finish(tracer.WithError(err))
	}()

	if err != nil {
		err = oops.Errorf(oops.E6xxNETWORK, fmt.Sprintf("[execute] [endpoint: %s] [error: %s]", endpoint, err.Error()))
		return nil, err
	}
	defer response.Body.Close()

	var code string
	switch response.StatusCode {
	case http.StatusOK, http.StatusCreated:
		return makeResponse(response)
	case http.StatusBadRequest:
		code = oops.E4xxCLIENTSIDE
	case http.StatusUnauthorized:
		code = oops.E4xxUNAUTHORIZED
	case http.StatusNotFound:
		code = oops.E4xxNOTFOUND
	case http.StatusUnprocessableEntity:
		code = oops.E4xxUNPROCESSABLE
	case http.StatusServiceUnavailable:
		code = oops.E5xxUNAVAILABLE
	default:
		code = oops.E5xxINTERNAL
	}

	err = oops.Errorf(code, toString(response.Body))
	return nil, err
}

func makeResponse(response *http.Response) (*Response, error) {
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, oops.Errorf(oops.E5xxINTERNAL, "error reading response")
	}

	return &Response{
		StatusCode: response.StatusCode,
		Body:       body,
	}, nil
}
