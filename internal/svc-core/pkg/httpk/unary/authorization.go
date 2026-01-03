package unaryHttpk

import (
	"encoding/base64"
	"strings"

	"github.com/konsultin/project-goes-here/internal/svc-core/pkg/httpk"
	f "github.com/valyala/fasthttp"
)

type BasicAuth struct {
	Username string
	Password string
}

func GetBasicAuth(ctx *f.RequestCtx) *BasicAuth {
	val, ok := ctx.UserValue(httpk.BasicAuth).(BasicAuth)

	if !ok {
		return nil
	}
	return &val
}

func AuthorizationMiddleware(next f.RequestHandler) f.RequestHandler {
	return func(ctx *f.RequestCtx) {
		// Get authorization header value
		authType, authValue := getAuthorizationHeaderValue(ctx, "Authorization")
		if authType == "" || authValue == "" {
			next(ctx)
			return
		}

		switch authType {
		case "Basic":
			parseBasicAuth(ctx, authValue)
		case "Bearer":
			ctx.SetUserValue(httpk.BearerToken, authValue)
		}

		next(ctx)
	}
}

func getAuthorizationHeaderValue(ctx *f.RequestCtx, header string) (authType string, value string) {
	val := ctx.Request.Header.Peek(header)
	if len(val) == 0 {
		return "", ""
	}

	authHeader := string(val)
	// Parse
	tmp := strings.SplitN(authHeader, " ", 2)
	if len(tmp) != 2 {
		return "", ""
	}

	return tmp[0], tmp[1]
}

func parseBasicAuth(ctx *f.RequestCtx, authValue string) {
	c, err := base64.StdEncoding.DecodeString(authValue)
	if err != nil {
		return
	}

	cs := string(c)
	s := strings.IndexByte(cs, ':')
	if s < 0 {
		return
	}

	// Get username and password
	username, password := cs[:s], cs[s+1:]

	// Set to context
	ctx.SetUserValue(httpk.BasicAuth, BasicAuth{username, password})
}
