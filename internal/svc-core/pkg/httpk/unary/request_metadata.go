package unaryHttpk

import (
	"time"

	"github.com/google/uuid"
	"github.com/konsultin/project-goes-here/internal/svc-core/pkg/httpk"
	f "github.com/valyala/fasthttp"
)

type RequestMetadata struct {
	RequestId string    `json:"requestId"`
	ClientIP  string    `json:"clientIP"`
	UserAgent string    `json:"userAgent"`
	StartedAt time.Time `json:"startedAt"`
}

// InjectRequestMetadata is a middleware that injects request metadata into the context
func InjectRequestMetadata(trustProxy string) func(next f.RequestHandler) f.RequestHandler {
	return func(next f.RequestHandler) f.RequestHandler {
		return func(ctx *f.RequestCtx) {
			meta := newRequestMetadata(ctx, trustProxy)

			// Set middleware context
			ctx.SetUserValue(httpk.RequestMetadata, meta)

			// Set response header
			ctx.Response.Header.Set("X-Request-Id", meta.RequestId)

			next(ctx)
		}
	}
}

func GetRequestMetadata(ctx *f.RequestCtx) *RequestMetadata {
	val, ok := ctx.UserValue(httpk.RequestMetadata).(RequestMetadata)
	if !ok {
		return nil
	}
	return &val
}

func newRequestMetadata(ctx *f.RequestCtx, trustProxy string) RequestMetadata {
	reqId := string(ctx.Request.Header.Peek("X-Request-Id"))
	if reqId == "" {
		reqId = uuid.NewString()
	}

	return RequestMetadata{
		RequestId: reqId,
		ClientIP:  httpk.GetClientIP(ctx, trustProxy),
		UserAgent: string(ctx.UserAgent()),
		StartedAt: time.Now(),
	}
}
