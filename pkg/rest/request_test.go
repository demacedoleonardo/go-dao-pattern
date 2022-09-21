package rest

import (
	"errors"
	oops "github.com/pedidosya/aragorn-service/pkg/errors"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

type errReader int

func (r errReader) Close() error {
	panic("implement me")
}

func (errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("test error")
}

type MockHttpClient struct {
	DoFunc func(r *http.Request) (*http.Response, error)
}

func (m *MockHttpClient) Do(r *http.Request) (*http.Response, error) {
	return m.DoFunc(r)
}

func TestBindResponseErr(t *testing.T) {
	parametrized := []struct {
		test     string
		response *Response
		err      error
		expected string
	}{
		{
			test:     "binding json error",
			response: &Response{StatusCode: http.StatusOK, Body: []byte("success")},
			err:      nil,
			expected: "error: code=internal message=invalid character 's' looking for beginning of value",
		},
		{
			test:     "request fail has error",
			response: nil,
			err:      oops.Errorf(oops.E5xxINTERNAL, "error internal"),
			expected: "error: code=internal message=error internal",
		},
	}

	for _, p := range parametrized {
		t.Run(p.test, func(t *testing.T) {
			var data map[string]interface{}
			err := BindResponse(p.response, p.err, data)

			assert.Nil(t, data)
			assert.NotNil(t, err)
			assert.Equal(t, p.expected, err.Error())
		})
	}
}

func TestBindResponseSuccess(t *testing.T) {
	parametrized := []struct {
		test     string
		response *Response
		in       interface{}
		expected interface{}
	}{
		{
			test:     "bind skipped",
			response: &Response{StatusCode: http.StatusOK, Body: []byte("success")},
			in:       nil,
			expected: nil,
		},
		{
			test:     "bind success",
			response: &Response{StatusCode: http.StatusOK, Body: []byte(`{"body": "fake"}`)},
			in:       new(map[string]interface{}),
			expected: &map[string]interface{}{"body": "fake"},
		},
	}

	for _, p := range parametrized {
		t.Run(p.test, func(t *testing.T) {
			err := BindResponse(p.response, nil, p.in)

			assert.Nil(t, err)
			assert.Equal(t, p.expected, p.in)
		})
	}
}

func TestNewRequest(t *testing.T) {
	headers := Headers{}
	headers.Add("Test", "Fake")

	r := NewRequest(WithClient(&http.Client{}), WithBody("fake body"), WithHeaders(headers))

	assert.NotNil(t, r.c)
	assert.NotNil(t, r.body)
	assert.NotNil(t, r.headers)
	assert.NotNil(t, r)
}

func TestNewRequest_Headers(t *testing.T) {
	headers := Headers{}
	headers.Add("Test", "Fake")

	r := NewRequest(WithHeaders(headers))

	assert.NotNil(t, r.headers)

	value, ok := r.headers.headers["Test"]
	assert.True(t, ok)
	assert.Equal(t, "Fake", value)
}

func TestNewRequest_Headers_Default_ContentTypeJSON(t *testing.T) {
	r := NewRequest()
	assert.NotNil(t, r.headers)

	value, ok := r.headers.headers["Content-Type"]
	assert.True(t, ok)
	assert.Equal(t, "application/json", value)
}

func TestWithBody_AllowString(t *testing.T) {
	body := "Fake Body"
	r := NewRequest(WithBody(body))

	assert.NotNil(t, r.body)

	buf := new(strings.Builder)
	_, err := io.Copy(buf, r.body)

	assert.Nil(t, err)
	assert.Equal(t, body, buf.String())
}

func TestWithBody_AllowSliceByte(t *testing.T) {
	body := []byte("Fake Body")
	r := NewRequest(WithBody(body))

	assert.NotNil(t, r.body)

	buf := new(strings.Builder)
	_, err := io.Copy(buf, r.body)

	assert.Nil(t, err)
	assert.Equal(t, string(body), buf.String())
}

func TestWithBody_AllowStruct(t *testing.T) {
	body := struct {
		Test string
	}{
		Test: "11-09-1983",
	}
	r := NewRequest(WithBody(body))

	buf := new(strings.Builder)
	_, err := io.Copy(buf, r.body)

	assert.Nil(t, err)
	assert.NotNil(t, r.body)
	assert.Equal(t, `{"Test":"11-09-1983"}`, buf.String())
}

func TestHttpResponse_toString_Success(t *testing.T) {
	response := http.Response{
		Status:     "Ok",
		StatusCode: 200,
		Body:       ioutil.NopCloser(strings.NewReader("11-09-1983")),
	}

	assert.Equal(t, "11-09-1983", toString(response.Body))
}

func TestHttpResponse_toString_OnErrorEmpty(t *testing.T) {
	response := http.Response{
		Status:     "Ok",
		StatusCode: 200,
		Body:       errReader(0),
	}

	assert.Equal(t, "", toString(response.Body))
}

func TestHttpResponse_makeResponse_Success(t *testing.T) {
	response := &http.Response{
		Status:     "Ok",
		StatusCode: 200,
		Body:       ioutil.NopCloser(strings.NewReader("11-09-1983")),
	}

	r, err := makeResponse(response)

	assert.Nil(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, http.StatusOK, r.StatusCode)
	assert.Equal(t, "11-09-1983", r.Body.ToString())
}

func TestHttpResponse_makeResponse_ErrorReadingBody(t *testing.T) {
	response := &http.Response{
		Status:     "Ok",
		StatusCode: 200,
		Body:       errReader(0),
	}

	r, err := makeResponse(response)

	assert.Nil(t, r)
	assert.NotNil(t, err)
	assert.Equal(t, "error: code=internal message=error reading response", err.Error())
}

func TestRequest_Post_ErrBindingRequest(t *testing.T) {
	r, err := NewRequest().Post(`{"fake": "url"}`)
	assert.Nil(t, r)
	assert.NotNil(t, err)
	assert.Equal(t, "error: code=internal message=error building request", err.Error())
}

func TestRequest_Post_DoRequestError(t *testing.T) {
	mock := &MockHttpClient{DoFunc: func(r *http.Request) (*http.Response, error) {
		return nil, errors.New("internal error")
	}}

	r, err := NewRequest(WithClient(mock)).Post("http://localhost/frodo")
	assert.Nil(t, r)
	assert.NotNil(t, err)
	assert.Equal(t, "error: code=network message=[execute] [endpoint: http://localhost/frodo] [error: internal error]", err.Error())
}

func TestRequest_Post_DoRequestErrorCodes(t *testing.T) {
	mock := &MockHttpClient{}

	parametrized := []struct {
		Test          string
		StatusCode    int
		ExpectedError error
	}{
		{
			Test: http.StatusText(http.StatusBadRequest),
			StatusCode: http.StatusBadRequest,
			ExpectedError: oops.Errorf(oops.E4xxCLIENTSIDE, http.StatusText(http.StatusBadRequest)),
		},
		{
			Test: http.StatusText(http.StatusUnauthorized),
			StatusCode: http.StatusUnauthorized,
			ExpectedError: oops.Errorf(oops.E4xxUNAUTHORIZED, http.StatusText(http.StatusUnauthorized)),
		},
		{
			Test: http.StatusText(http.StatusUnprocessableEntity),
			StatusCode: http.StatusUnprocessableEntity,
			ExpectedError: oops.Errorf(oops.E4xxUNPROCESSABLE, http.StatusText(http.StatusUnprocessableEntity)),
		},
		{
			Test: http.StatusText(http.StatusNotFound),
			StatusCode: http.StatusNotFound,
			ExpectedError: oops.Errorf(oops.E4xxNOTFOUND, http.StatusText(http.StatusNotFound)),
		},
		{
			Test: http.StatusText(http.StatusInternalServerError),
			StatusCode: http.StatusInternalServerError,
			ExpectedError: oops.Errorf(oops.E5xxINTERNAL, http.StatusText(http.StatusInternalServerError)),
		},
	}

	for _, p := range parametrized {
		t.Run(p.Test, func(t *testing.T) {
			mock.DoFunc = func(r *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: p.StatusCode,
					Body:       ioutil.NopCloser(strings.NewReader(http.StatusText(p.StatusCode))),
				}, nil
			}

			r, err := NewRequest(WithClient(mock)).Post("http://localhost/frodo")

			assert.Nil(t, r)
			assert.NotNil(t, err)
			assert.Equal(t, p.ExpectedError.Error(), err.Error())
		})
	}
}

func TestRequest_Post_DoRequestSuccess(t *testing.T) {
	mock := &MockHttpClient{}

	s := struct {
		Test          string
		StatusCode    int
		ExpectedError error
	}{
			Test: http.StatusText(http.StatusOK),
			StatusCode: http.StatusOK,
			ExpectedError: nil,
	}

	mock.DoFunc = func(r *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: s.StatusCode,
			Body:       ioutil.NopCloser(strings.NewReader(http.StatusText(s.StatusCode))),
		}, nil
	}

	r, err := NewRequest(WithClient(mock)).Post("http://localhost/frodo")

	assert.Nil(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, http.StatusOK, r.StatusCode)
}

func TestRequest_Put_DoRequestSuccess(t *testing.T) {
	mock := &MockHttpClient{}

	s := struct {
		Test          string
		StatusCode    int
		ExpectedError error
	}{
		Test: http.StatusText(http.StatusOK),
		StatusCode: http.StatusOK,
		ExpectedError: nil,
	}

	mock.DoFunc = func(r *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: s.StatusCode,
			Body:       ioutil.NopCloser(strings.NewReader(http.StatusText(s.StatusCode))),
		}, nil
	}

	r, err := NewRequest(WithClient(mock)).Put("http://localhost/frodo")

	assert.Nil(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, http.StatusOK, r.StatusCode)
}

func TestRequest_Get_DoRequestSuccess(t *testing.T) {
	mock := &MockHttpClient{}

	s := struct {
		Test          string
		StatusCode    int
		ExpectedError error
	}{
		Test: http.StatusText(http.StatusOK),
		StatusCode: http.StatusOK,
		ExpectedError: nil,
	}

	mock.DoFunc = func(r *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: s.StatusCode,
			Body:       ioutil.NopCloser(strings.NewReader(http.StatusText(s.StatusCode))),
		}, nil
	}

	r, err := NewRequest(WithClient(mock)).Get("http://localhost/frodo")

	assert.Nil(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, http.StatusOK, r.StatusCode)
}
