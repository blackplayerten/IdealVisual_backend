package api

import (
	"encoding/json"

	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

func (s *Server) handleRequest(ctx *fasthttp.RequestCtx) {
	s.recoverMiddleware(s.accessLogMiddleware(fasthttp.TimeoutHandler(
		s.authMiddleware(s.route),
		s.cfg.HTTP.Timeout, "")))(ctx)
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
		s.handleAccount(ctx)
	case "/post":
		s.handlePost(ctx)
	case "/upload":
		s.upload(ctx)
	default:
		ctx.SetStatusCode(fasthttp.StatusNotFound)
	}
}

func (s *Server) handleSession(ctx *fasthttp.RequestCtx) {
	switch string(ctx.Method()) {
	case fasthttp.MethodPost:
		s.newSession(ctx)
	case fasthttp.MethodDelete:
		s.deleteSession(ctx)
	default:
		ctx.SetStatusCode(fasthttp.StatusMethodNotAllowed)
	}
}

func (s *Server) handleAccount(ctx *fasthttp.RequestCtx) {
	switch string(ctx.Method()) {
	case fasthttp.MethodGet:
		s.getAccount(ctx)
	case fasthttp.MethodPost:
		s.newAccount(ctx)
	case fasthttp.MethodPut:
		s.updateAccount(ctx)
	default:
		ctx.SetStatusCode(fasthttp.StatusMethodNotAllowed)
	}
}

func (s *Server) handlePost(ctx *fasthttp.RequestCtx) {
	switch string(ctx.Method()) {
	case fasthttp.MethodGet:
		s.getPost(ctx)
	case fasthttp.MethodPost:
		s.newPost(ctx)
	case fasthttp.MethodPut:
		s.updatePost(ctx)
	case fasthttp.MethodDelete:
		s.deletePost(ctx)
	default:
		ctx.SetStatusCode(fasthttp.StatusMethodNotAllowed)
	}
}

func (s *Server) writeJSONResponse(ctx *fasthttp.RequestCtx, resp json.Marshaler) {
	j, err := resp.MarshalJSON()
	if err != nil {
		s.l.Error("cannot marshal json", zap.Error(err))
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)

		return
	}

	ctx.SetContentTypeBytes([]byte("application/json"))
	ctx.Write(j) // nolint:errcheck
}
