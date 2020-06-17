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
	_ caddy.Provisioner           = (*Middleware)(nil)
	_ caddy.Validator             = (*Middleware)(nil)
	_ caddyhttp.MiddlewareHandler = (*Middleware)(nil)
)

func init() {
	caddy.RegisterModule(Middleware{})
}

type Middleware struct {
	logger *zap.Logger
	Disabled bool `json:"disabled"`
}

func (m Middleware) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.x_request_id",
		New: func() caddy.Module { return new(Middleware) },
	}
}

// Provision sets up the module.
func (m *Middleware) Provision(ctx caddy.Context) error {
	m.logger = ctx.Logger(m) // g.logger is a *zap.Logger
	return nil
}

// Validate validates that the module has a usable config.
func (m Middleware) Validate() error {
	// TODO: validate the module's setup
	return nil
}

// ServeHTTP implements caddyhttp.MiddlewareHandler.
func (m Middleware) ServeHTTP(writer http.ResponseWriter, request *http.Request, nextHandler caddyhttp.Handler) error {
	if !m.Disabled {
		requestId := request.Header.Get("X-Request-Id")
		if len(strings.TrimSpace(requestId)) == 0 {
			request.Header.Set("X-Request-Id", NewXRequestId())
			m.logger.Debug("Adding X-Request-Id request header and generating new value.",
				zap.String("X-Request-Id", request.Header.Get("X-Request-Id")),
			)
		} else {
			m.logger.Debug("Found existing X-Request-Id request header and using it.",
				zap.String("X-Request-Id", request.Header.Get("X-Request-Id")),
			)
		}
	}
	return nextHandler.ServeHTTP(writer, request)
}

func NewXRequestId() string {
	return uuid.New().String()
}
