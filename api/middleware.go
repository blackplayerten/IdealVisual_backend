package api

import (
	"bytes"

	"github.com/valyala/fasthttp"
	"go.uber.org/zap"

	"github.com/blackplayerten/IdealVisual_backend/session"
)

func (s *Server) recoverMiddleware(handler fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		defer func() {
			if r := recover(); r != nil {
				s.l.Error("recovered panic", zap.Reflect("value", r), zap.Stack("stacktrace"))
				ctx.SetStatusCode(fasthttp.StatusInternalServerError)
			}
		}()

		handler(ctx)
	}
}

const (
	KeyIsAuthenticated = "auth_is_authenticated"
	KeySessionID       = "auth_session_id"
	KeyUserID          = "auth_user_id"
)

func (s *Server) authMiddleware(handler fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		var bearer = []byte("Bearer ")

		val := ctx.Request.Header.Peek("Authorization")
		if len(val) != 0 {
			var token []byte
			if bytes.HasPrefix(val, bearer) {
				token = bytes.TrimPrefix(val, bearer)
			}

			if len(token) != 0 {
				userID, err := s.sessionSvc.Get(string(token))
				switch err {
				case nil:
					ctx.SetUserValue(KeyIsAuthenticated, true)
					ctx.SetUserValue(KeySessionID, token)
					ctx.SetUserValue(KeyUserID, userID)
				case session.ErrKeyNotFound:
					ctx.SetUserValue(KeyIsAuthenticated, false)
				default:
					s.l.Error("cannot get id from token", zap.ByteString("token", token), zap.Error(err))
					ctx.SetStatusCode(fasthttp.StatusInternalServerError)

					return
				}
			} else {
				ctx.SetUserValue(KeyIsAuthenticated, false)
			}
		} else {
			ctx.SetUserValue(KeyIsAuthenticated, false)
		}

		handler(ctx)
	}
}

func setContentTypeToAppJSON(handler fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		handler(ctx)

		ctx.SetContentTypeBytes([]byte("application/json"))
	}
}
