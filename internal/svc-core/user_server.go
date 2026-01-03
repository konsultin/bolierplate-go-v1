package svcCore

import (
	"github.com/konsultin/project-goes-here/dto"
	unaryHttpk "github.com/konsultin/project-goes-here/internal/svc-core/pkg/httpk/unary"
	httpkPkg "github.com/konsultin/project-goes-here/internal/svc-core/pkg/httpk"
	f "github.com/valyala/fasthttp"
)

func (s *Server) HandleCreateAnonymousUserSession(ctx *f.RequestCtx) (*dto.CreateAnonymousSession_Result, error) {
	basicAuth := unaryHttpk.GetBasicAuth(ctx)
	if basicAuth == nil {
		s.log.Errorf("Basic Auth is not set in header")
		return nil, s.wrapError(ctx, httpkPkg.UnauthorizedError)
	}

	// Init Service
	svc, err := s.NewService(ctx)
	if err != nil {
		s.log.Errorf("Failed to create service: %v", err)
		return nil, err
	}

	// Call Service
	result, err := svc.CreateAnonymousUserSession(ctx, basicAuth.Username, basicAuth.Password)
	if err != nil {
		s.log.Errorf("Failed to create user session: %v", err)
		return nil, err
	}

	return result, nil
}
