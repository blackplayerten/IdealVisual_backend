package api

import (
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

func (s *Server) handleRequest(ctx *fasthttp.RequestCtx) {
	s.recoverMiddleware(fasthttp.TimeoutHandler(setContentTypeToAppJSON(
		s.authMiddleware(s.route),
	), s.cfg.HTTP.Timeout, ""))(ctx)
}

func (s *Server) handleError(_ *fasthttp.RequestCtx, err error) {
	// TODO: more info
	s.l.Error("server core error", zap.Error(err))
}

func (s *Server) route(ctx *fasthttp.RequestCtx) {
	switch string(ctx.Path()) {
	case "/session":
		s.handleSession(ctx)
	case "/account":

	case "/post":

	case "/upload":
		s.upload(ctx)
	default:
		ctx.SetStatusCode(fasthttp.StatusNotFound)
	}
}

func (s *Server) handleSession(ctx *fasthttp.RequestCtx) {
	switch string(ctx.Method()) {
	case fasthttp.MethodPost:
		// TODO: think about authorization type
	case fasthttp.MethodDelete:

	default:
		ctx.SetStatusCode(fasthttp.StatusMethodNotAllowed)
	}
}
