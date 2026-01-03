package unaryHttpk

import (
	"fmt"
	"runtime/debug"
	"time"

	"github.com/valyala/fasthttp"
)

// RecoveryMiddleware handles panics and logs request completion
func RecoveryMiddleware(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		startedAt := time.Now()

		defer func() {
			if r := recover(); r != nil {
				// Log panic
				stack := string(debug.Stack())
				log.Errorf("Panic occurred during HTTP request. Error=%v Stack=%s", r, stack)

				// Return 500
				ctx.Error(fmt.Sprintf("Internal Server Error: %v", r), fasthttp.StatusInternalServerError)
			}

			// Log request completion (access log)
			logAccess(ctx, startedAt)
		}()

		next(ctx)
	}
}

func logAccess(ctx *fasthttp.RequestCtx, startedAt time.Time) {
	reqMeta := GetRequestMetadata(ctx)
	if reqMeta == nil {
		// Try to reconstruct basic metadata if missing (should be injected by metadata middleware)
		reqMeta = &RequestMetadata{
			StartedAt: startedAt,
		}
	}

	// Calculate duration
	duration := time.Since(startedAt)

	// Log details
	log.Infof("Handled HTTP Request. Method=%s Path=%s Status=%d Duration=%s ClientIP=%s UserAgent=%s RequestId=%s",
		string(ctx.Method()),
		string(ctx.Path()),
		ctx.Response.StatusCode(),
		duration,
		reqMeta.ClientIP,
		reqMeta.UserAgent,
		reqMeta.RequestId,
	)
}
