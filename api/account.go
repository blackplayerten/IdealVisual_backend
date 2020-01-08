package api

import (
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"

	"github.com/blackplayerten/IdealVisual_backend/account"
	"github.com/blackplayerten/IdealVisual_backend/database"
)

func (s *Server) getAccount(ctx *fasthttp.RequestCtx) {
	if !ctx.UserValue(KeyIsAuthenticated).(bool) {
		ctx.SetStatusCode(fasthttp.StatusUnauthorized)
		return
	}

	userID := ctx.UserValue(KeyUserID).(uint64)

	acc, err := s.accountSvc.GetByID(userID)
	if err != nil {
		if err == database.ErrNotFound {
			ctx.SetStatusCode(fasthttp.StatusNotFound)
			return
		}

		s.l.Error("cannot get user by id", zap.Error(err), zap.Uint64("user_id", userID))

		return
	}

	s.writeJSONResponse(ctx, acc)
}

func (s *Server) newAccount(ctx *fasthttp.RequestCtx) {
	var info account.FullAccount
	if err := info.UnmarshalJSON(ctx.PostBody()); err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	validationErrors := validateAll(&info)

	if len(validationErrors) != 0 {
		ctx.SetStatusCode(fasthttp.StatusUnprocessableEntity)

		errors := Errors{Errors: validationErrors}
		s.writeJSONResponse(ctx, &errors)

		return
	}

	acc, err := s.accountSvc.New(&info)
	if err != nil {
		s.processAccountCreationUpdateError(ctx, err)
		return
	}

	s.loginUserAndSendInfo(ctx, acc)
}

func (s *Server) updateAccount(ctx *fasthttp.RequestCtx) {
	if !ctx.UserValue(KeyIsAuthenticated).(bool) {
		ctx.SetStatusCode(fasthttp.StatusUnauthorized)
		return
	}

	var upd account.FullAccount
	if err := upd.UnmarshalJSON(ctx.PostBody()); err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	validationErrors := make([]*FieldError, 0, 3)

	if upd.Email != "" {
		if validationError := validateEmail(upd.Email); validationError != nil {
			validationErrors = append(validationErrors, validationError)
		}
	}

	if upd.Password != "" {
		if validationError := validatePassword(upd.Password); validationError != nil {
			validationErrors = append(validationErrors, validationError)
		}
	}

	if upd.Username != "" {
		if validationError := validateUsername(upd.Username); validationError != nil {
			validationErrors = append(validationErrors, validationError)
		}
	}

	if len(validationErrors) != 0 {
		ctx.SetStatusCode(fasthttp.StatusUnprocessableEntity)

		errors := Errors{Errors: validationErrors}
		s.writeJSONResponse(ctx, &errors)

		return
	}

	upd.ID = ctx.UserValue(KeyUserID).(uint64)

	acc, err := s.accountSvc.Update(&upd)
	if err != nil {
		if err == database.ErrNotFound {
			ctx.SetStatusCode(fasthttp.StatusNotFound)
			return
		}

		s.processAccountCreationUpdateError(ctx, err)

		return
	}

	s.writeJSONResponse(ctx, acc)
}

func (s *Server) processAccountCreationUpdateError(ctx *fasthttp.RequestCtx, err error) {
	switch tErr := err.(type) {
	case database.UniqueConstraintViolationError:
		ctx.SetStatusCode(fasthttp.StatusUnprocessableEntity)

		errors := Errors{Errors: []*FieldError{{
			Field:   tErr.Field,
			Reasons: []string{AlreadyExists},
		}}}
		s.writeJSONResponse(ctx, &errors)
	default:
		s.l.Error("cannot create new account", zap.Error(err))
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
	}
}
