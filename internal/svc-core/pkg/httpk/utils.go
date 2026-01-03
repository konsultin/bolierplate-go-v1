package httpk

import (
	"strings"

	f "github.com/valyala/fasthttp"
)

func GetClientIP(ctx *f.RequestCtx, trustProxy string) string {
	var ip string

	ip = string(ctx.Request.Header.Peek("X-Forwarded-For"))
	if ip != "" {
		ips := strings.Split(ip, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	ip = string(ctx.Request.Header.Peek("X-Real-IP"))
	if ip != "" {
		return strings.TrimSpace(ip)
	}

	return ctx.RemoteIP().String()
}

func GetHeaderFromContext(ctx *f.RequestCtx, key string) (string, bool) {
	val := ctx.Request.Header.Peek(key)
	if val == nil {
		return "", false
	}
	return string(val), true
}
