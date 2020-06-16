package xrequestid

import (
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"net/http"
	"testing"
)

func TestNewUuid(t *testing.T) {
	result := NewUuid()
	assert.Len(t, result, 36)
}

func TestXRequestId_ServeHTTP(t *testing.T) {
	type fields struct {
		logger *zap.Logger
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
		{"Given no X-Request-Id header, when ServeHTTP, then generate one and add to request.",
			fields{
				logger: nil,
			},
			args{
				writer: nil,
				request: &http.Request{
					Header: http.Header{},
				},
				nextHandler: AssertionHandler{
					t: t,
				},
			},
			false,
		},
		{"Given an empty X-Request-Id header, when ServeHTTP, then generate one and add to request.",
			fields{
				logger: nil,
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
				},
			},
			false,
		},
		{"Given an all whitespace X-Request-Id header, when ServeHTTP, then generate one and add to request.",
			fields{
				logger: nil,
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
				},
			},
			false,
		},
		{"Given a non empty X-Request-Id header, when ServeHTTP, then leave it as is.",
			fields{
				logger: nil,
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
					expected: "66b5651c-b01b-11ea-b3de-0242ac130004",
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			xr := XRequestId{
				logger: tt.fields.logger,
			}
			if err := xr.ServeHTTP(tt.args.writer, tt.args.request, tt.args.nextHandler); (err != nil) != tt.wantErr {
				t.Errorf("ServeHTTP() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

type AssertionHandler struct {
	t        *testing.T
	expected string
}

func (a AssertionHandler) ServeHTTP(_ http.ResponseWriter, request *http.Request) error {
	actual := request.Header.Get("X-Request-Id")
	assert.NotEmpty(a.t, actual)
	if len(a.expected) != 0 {
		assert.Equal(a.t, a.expected, actual)
	}
	return nil
}
