package xrequestid

import (
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

var (
	// Interface guards to ensure this module implements the following interfaces
	_ caddy.Provisioner           = (*XRequestId)(nil)
	_ caddy.Validator             = (*XRequestId)(nil)
	_ caddyhttp.MiddlewareHandler = (*XRequestId)(nil)
)

func init() {
	caddy.RegisterModule(XRequestId{})
}

type XRequestId struct {
	logger *zap.Logger
}

func (x XRequestId) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.headers.xrequestid",
		New: func() caddy.Module { return new(XRequestId) },
	}
}

// Provision sets up the module.
func (x *XRequestId) Provision(ctx caddy.Context) error {
	x.logger = ctx.Logger(x) // g.logger is a *zap.Logger
	return nil
}

// Validate validates that the module has a usable config.
func (x XRequestId) Validate() error {
	// TODO: validate the module's setup
	return nil
}

// ServeHTTP implements caddyhttp.MiddlewareHandler.
func (x XRequestId) ServeHTTP(writer http.ResponseWriter, request *http.Request, nextHandler caddyhttp.Handler) error {
	requestId := request.Header.Get("X-Request-Id")
	if len(strings.TrimSpace(requestId)) == 0 {
		request.Header.Set("X-Request-Id", NewUuid())
		x.logger.Debug("Adding X-Request-Id request header and generating new value.",
			zap.String("X-Request-Id", request.Header.Get("X-Request-Id")),
		)
	} else {
		x.logger.Debug("Found existing X-Request-Id request header and using it.",
			zap.String("X-Request-Id", request.Header.Get("X-Request-Id")),
		)
	}
	return nextHandler.ServeHTTP(writer, request)
}

func NewUuid() string {
	return uuid.New().String()
}
