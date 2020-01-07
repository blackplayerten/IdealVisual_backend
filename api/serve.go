package api

import (
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"

	"github.com/blackplayerten/IdealVisual_backend/account"
	"github.com/blackplayerten/IdealVisual_backend/post"
	"github.com/blackplayerten/IdealVisual_backend/session"
)

type Server struct {
	cfg *Config

	s fasthttp.Server
	l *zap.Logger

	sessionSvc *session.Service
	accountSvc *account.Service
	postSvc    *post.Service
}

func New(
	cfg *Config,
	l *zap.Logger,
	sessionSvc *session.Service,
	accountSvc *account.Service,
	postSvc *post.Service,
) (*Server, error) {
	s := &Server{
		cfg: cfg,

		s: fasthttp.Server{
			NoDefaultContentType:  true,
			NoDefaultServerHeader: true,

			MaxRequestBodySize: cfg.HTTP.BodyLimitMB * 1024 * 1024,
		},
		l: l,

		sessionSvc: sessionSvc,
		accountSvc: accountSvc,
		postSvc:    postSvc,
	}
	s.s.Handler = s.handleRequest
	s.s.ErrorHandler = s.handleError

	if err := checkRoot(s.l, s.cfg.Static.Root); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Server) ListenAndServe() error {
	return s.s.ListenAndServe(s.cfg.HTTP.Addr)
}

func (s *Server) Shutdown() error {
	err := s.s.Shutdown()
	s.l.Sync() // nolint:errcheck

	return err
}
