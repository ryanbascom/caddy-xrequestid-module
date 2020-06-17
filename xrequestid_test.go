package xrequestid

import (
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
	"net/http"
	"testing"
)

func TestNewId(t *testing.T) {
	result := NewId()
	assert.Len(t, result, 36)
}

func TestXRequestId_ServeHTTP(t *testing.T) {
	type fields struct {
		logger *zap.Logger
		disabled bool
	}
	type args struct {
		writer      http.ResponseWriter
		request     *http.Request
		nextHandler caddyhttp.Handler
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"Given no X-Request-Id header and disabled, when ServeHTTP, then leave as is.",
			fields{
				logger: zaptest.NewLogger(t),
				disabled: true,
			},
			args{
				writer: nil,
				request: &http.Request{
					Header: http.Header{},
				},
				nextHandler: AssertionHandler{
					t: t,
					assertions: func(t *testing.T, r *http.Request) {
						actual := r.Header.Get("X-Request-Id")
						assert.Empty(t, actual)
					},
				},
			},
			false,
		},
		{"Given no X-Request-Id header, when ServeHTTP, then generate one and add to request.",
			fields{
				logger: zaptest.NewLogger(t),
			},
			args{
				writer: nil,
				request: &http.Request{
					Header: http.Header{},
				},
				nextHandler: AssertionHandler{
					t: t,
					assertions: func(t *testing.T, r *http.Request) {
						actual := r.Header.Get("X-Request-Id")
						assert.NotEmpty(t, actual)
					},
				},
			},
			false,
		},
		{"Given an X-Request-Id header value that is empty, when ServeHTTP, then generate one and add to request.",
			fields{
				logger: zaptest.NewLogger(t),
			},
			args{
				writer: nil,
				request: &http.Request{
					Header: http.Header{
						"X-Request-Id": []string{""},
					},
				},
				nextHandler: AssertionHandler{
					t: t,
					assertions: func(t *testing.T, r *http.Request) {
						actual := r.Header.Get("X-Request-Id")
						assert.NotEmpty(t, actual)
					},
				},
			},
			false,
		},
		{"Given an X-Request-Id header value that is all whitespace, when ServeHTTP, then generate one and add to request.",
			fields{
				logger: zaptest.NewLogger(t),
			},
			args{
				writer: nil,
				request: &http.Request{
					Header: http.Header{
						"X-Request-Id": []string{"  "},
					},
				},
				nextHandler: AssertionHandler{
					t: t,
					assertions: func(t *testing.T, r *http.Request) {
						actual := r.Header.Get("X-Request-Id")
						assert.NotEmpty(t, actual)
					},
				},
			},
			false,
		},
		{"Given an X-Request-Id header value, when ServeHTTP, then leave it as is.",
			fields{
				logger: zaptest.NewLogger(t),
			},
			args{
				writer: nil,
				request: &http.Request{
					Header: http.Header{
						"X-Request-Id": []string{"66b5651c-b01b-11ea-b3de-0242ac130004"},
					},
				},
				nextHandler: AssertionHandler{
					t:        t,
					assertions: func(t *testing.T, r *http.Request) {
						actual := r.Header.Get("X-Request-Id")
						assert.Equal(t, "66b5651c-b01b-11ea-b3de-0242ac130004", actual)
					},
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			xr := XRequestId{
				logger: tt.fields.logger,
				Disabled: tt.fields.disabled,
			}
			if err := xr.ServeHTTP(tt.args.writer, tt.args.request, tt.args.nextHandler); (err != nil) != tt.wantErr {
				t.Errorf("ServeHTTP() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

type AssertionHandler struct {
	t        *testing.T
	assertions assertFunc
}

func (a AssertionHandler) ServeHTTP(_ http.ResponseWriter, request *http.Request) error {
	a.assertions(a.t, request)
	return nil
}

type assertFunc func(t *testing.T, r *http.Request)


