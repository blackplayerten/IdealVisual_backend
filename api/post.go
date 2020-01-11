package api

import (
	"github.com/google/uuid"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"

	"github.com/blackplayerten/IdealVisual_backend/database"
	"github.com/blackplayerten/IdealVisual_backend/post"
)

func (s *Server) getPost(ctx *fasthttp.RequestCtx) {
	if !ctx.UserValue(KeyIsAuthenticated).(bool) {
		ctx.SetStatusCode(fasthttp.StatusUnauthorized)
		return
	}

	ids := ctx.QueryArgs().PeekMulti("id")

	var uuids []string
	if len(ids) != 0 {
		uuids = make([]string, 0, len(ids))

		for _, id := range ids {
			parsed, err := uuid.ParseBytes(id)
			if err != nil {
				ctx.SetStatusCode(fasthttp.StatusBadRequest)
				return
			}

			uuids = append(uuids, parsed.String())
		}
	}

	userID := ctx.UserValue(KeyUserID).(uint64)

	posts, err := s.postSvc.Get(userID, uuids)
	if err != nil {
		s.l.Error("cannot get posts",
			zap.Error(err),
			zap.Uint64("user_id", userID),
			zap.Strings("ids", uuids),
		)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)

		return
	}

	if len(posts) != 0 {
		jsonPosts := Posts(posts)
		s.writeJSONResponse(ctx, &jsonPosts)
	} else {
		ctx.SetContentTypeBytes([]byte("application/json"))
		ctx.Write([]byte("[]"))
	}
}

func (s *Server) newPost(ctx *fasthttp.RequestCtx) {
	if !ctx.UserValue(KeyIsAuthenticated).(bool) {
		ctx.SetStatusCode(fasthttp.StatusUnauthorized)
		return
	}

	var pst post.Post
	if err := pst.UnmarshalJSON(ctx.PostBody()); err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	if pst.Photo == "" {
		ctx.SetStatusCode(fasthttp.StatusUnprocessableEntity)

		errors := Errors{Errors: []*FieldError{{
			Field:   "photo",
			Reasons: []string{WrongLen},
		}}}
		s.writeJSONResponse(ctx, &errors)

		return
	}

	pst.Acc = ctx.UserValue(KeyUserID).(uint64)

	newP, err := s.postSvc.New(&pst)
	if err != nil {
		if _, ok := err.(database.ForeignKeyViolation); ok {
			ctx.SetStatusCode(fasthttp.StatusNotFound)
			return
		}

		s.l.Error("cannot create new post", zap.Error(err), zap.Reflect("post", pst))
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)

		return
	}

	s.writeJSONResponse(ctx, newP)
}

func (s *Server) updatePost(ctx *fasthttp.RequestCtx) {
	if !ctx.UserValue(KeyIsAuthenticated).(bool) {
		ctx.SetStatusCode(fasthttp.StatusUnauthorized)
		return
	}

	var pst post.Post
	if err := pst.UnmarshalJSON(ctx.PostBody()); err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	pst.Acc = ctx.UserValue(KeyUserID).(uint64)

	upd, err := s.postSvc.Update(&pst)
	if err != nil {
		if err == database.ErrNotFound {
			ctx.SetStatusCode(fasthttp.StatusNotFound)
			return
		}

		s.l.Error("cannot update post", zap.Error(err), zap.Reflect("post", pst))
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)

		return
	}

	s.writeJSONResponse(ctx, upd)
}

func (s *Server) deletePost(ctx *fasthttp.RequestCtx) {
	if !ctx.UserValue(KeyIsAuthenticated).(bool) {
		ctx.SetStatusCode(fasthttp.StatusUnauthorized)
		return
	}

	ids := ctx.QueryArgs().PeekMulti("id")
	if len(ids) == 0 {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	uuids := make([]string, 0, len(ids))

	for _, id := range ids {
		parsed, err := uuid.ParseBytes(id)
		if err != nil {
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			return
		}

		uuids = append(uuids, parsed.String())
	}

	userID := ctx.UserValue(KeyUserID).(uint64)
	if err := s.postSvc.Delete(userID, uuids); err != nil {
		s.l.Error("cannot delete posts",
			zap.Error(err),
			zap.Uint64("user_id", userID),
			zap.Strings("ids", uuids),
		)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
	}
}
