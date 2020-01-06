package api

import (
	"bytes"
	"time"

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
					ctx.SetUserValue(KeySessionID, string(token))
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

// TODO: request ID
func (s *Server) accessLogMiddleware(handler fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		body := ctx.PostBody()

		size := len(body)
		if size > 1024 {
			size = 1024
		}

		// sensitive information
		s.l.Info("access request",
			zap.ByteString("method", ctx.Method()),
			zap.ByteString("path", ctx.RequestURI()),
			zap.ByteString("content-type", ctx.Request.Header.Peek("Content-Type")),
			zap.ByteString("body", body[:size]),

			zap.ByteString("authorization", ctx.Request.Header.Peek("Authorization")),
		)

		started := time.Now()

		handler(ctx)

		ended := time.Since(started)

		body = ctx.Response.Body()

		size = len(body)
		if size > 1024 {
			size = 1024
		}

		// sensitive information
		s.l.Info("access response",
			zap.Int("code", ctx.Response.StatusCode()),
			zap.ByteString("content-type", ctx.Response.Header.Peek("Content-Type")),
			zap.ByteString("body", body[:size]),
			zap.Duration("timing", ended),
		)
	}
}
