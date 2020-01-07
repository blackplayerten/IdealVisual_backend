package api

import (
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"

	"github.com/blackplayerten/IdealVisual_backend/account"
	"github.com/blackplayerten/IdealVisual_backend/database"
)

func (s *Server) newSession(ctx *fasthttp.RequestCtx) {
	if ctx.UserValue(KeyIsAuthenticated).(bool) {
		// Send old token
		userID := ctx.UserValue(KeyUserID).(uint64)
		sessionID := ctx.UserValue(KeySessionID).(string)

		acc, err := s.accountSvc.GetByID(userID)
		if err != nil {
			if err == database.ErrNotFound {
				ctx.SetStatusCode(fasthttp.StatusForbidden)
				return
			}

			s.l.Error("cannot get user by id",
				zap.Error(err),
				zap.String("session_id", sessionID),
				zap.Uint64("user_id", userID),
			)

			return
		}

		info := AccWithToken{
			Token:   sessionID,
			Account: acc,
		}
		s.writeJSONResponse(ctx, &info)

		return
	}

	var cre account.Credentials
	if err := cre.UnmarshalJSON(ctx.PostBody()); err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	acc, ok, err := s.accountSvc.CheckCredentials(&cre)
	if err != nil {
		s.l.Error("cannot check credentials", zap.Error(err))
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)

		return
	}

	if ok {
		s.loginUserAndSendInfo(ctx, acc)
	} else {
		ctx.SetStatusCode(fasthttp.StatusForbidden)
	}
}

func (s *Server) deleteSession(ctx *fasthttp.RequestCtx) {
	if !ctx.UserValue(KeyIsAuthenticated).(bool) {
		return
	}

	sessionID := ctx.UserValue(KeySessionID).(string)
	if err := s.sessionSvc.Delete(sessionID); err != nil {
		s.l.Error("cannot delete user session",
			zap.Error(err),
			zap.String("session_id", sessionID),
			zap.Uint64("user_id", ctx.UserValue(KeyUserID).(uint64)),
		)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
	}
}

func (s *Server) loginUserAndSendInfo(ctx *fasthttp.RequestCtx, acc *account.Account) {
	sessionID, err := s.sessionSvc.Create(acc.ID)
	if err != nil {
		s.l.Error("cannot create session", zap.Error(err), zap.Uint64("user_id", acc.ID))
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)

		return
	}

	info := AccWithToken{
		Token:   sessionID,
		Account: acc,
	}
	s.writeJSONResponse(ctx, &info)
}
