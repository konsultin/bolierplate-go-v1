package middleware

import (
	"fmt"

	unaryHttpk "github.com/konsultin/project-goes-here/internal/svc-core/pkg/httpk/unary"
	"github.com/valyala/fasthttp"
)

// Init wires all middlewares and returns a ready-to-use handler.
func Init(cfg Config) (handler func(ctx *fasthttp.RequestCtx), err error) {
	if cfg.Handler == nil {
		return nil, fmt.Errorf("middleware: handler is nil")
	}
	if cfg.Logger == nil {
		return nil, fmt.Errorf("middleware: logger is nil")
	}
	if cfg.OnError == nil {
		return nil, fmt.Errorf("middleware: onError responder is nil")
	}

	rl := NewRateLimiter(cfg.RateLimitRPS, cfg.RateLimitBurst)
	metrics := cfg.Metrics
	if metrics == nil {
		metrics = NewMetrics()
	}

	handler = Chain(cfg.Handler,
		Recovery(cfg.Logger, cfg.OnError),
		RequestID(),
		Logging(cfg.Logger, metrics),
		RateLimit(rl, cfg.Logger, cfg.OnError),
		CORS(cfg.CORSAllowOrigins),
		unaryHttpk.AuthorizationMiddleware,
	)

	return handler, nil
}
