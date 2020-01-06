package api

import (
	"errors"
	"io"
	"os"
	"path/filepath"

	uuid "github.com/satori/go.uuid"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

//easyjson:json
type UploadedPath struct {
	Path string `json:"path"`
}

func (s *Server) upload(ctx *fasthttp.RequestCtx) {
	if !ctx.IsPost() {
		ctx.SetStatusCode(fasthttp.StatusMethodNotAllowed)
		return
	}

	if !ctx.UserValue(KeyIsAuthenticated).(bool) {
		ctx.SetStatusCode(fasthttp.StatusUnauthorized)
		return
	}

	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	headers := form.File["file"]
	if len(headers) != 1 {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	header := headers[0]

	name := header.Filename
	if !s.cfg.Static.KeepOriginalName {
		name = uuid.NewV4().String() + filepath.Ext(name)
	}

	file, err := header.Open()
	if err != nil {
		s.l.Error("cannot open file header", zap.Error(err))
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)

		return
	}
	defer file.Close() // nolint:errcheck

	if err = saveFile(file, s.cfg.Static.Root, name); err != nil {
		s.l.Error("cannot save file", zap.Error(err))
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)

		return
	}

	uploadedTo := UploadedPath{Path: name}
	s.writeJSONResponse(ctx, &uploadedTo)
}

func saveFile(file io.Reader, path, filename string) error {
	if err := os.MkdirAll(path, 0755); err != nil { // Checks existence too
		return err
	}

	fullPath := filepath.Join(path, filename)
	if _, err := os.Stat(fullPath); !os.IsNotExist(err) {
		if err == nil {
			return os.ErrExist
		}

		return err
	} // else file doesn't exist

	f, err := os.OpenFile(fullPath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close() // nolint:errcheck

	if _, err = io.Copy(f, file); err != nil {
		return err
	}

	return nil
}

var errRoot = errors.New("root path is invalid")

func checkRoot(l *zap.Logger, root string) error {
	if info, err := os.Stat(root); err != nil {
		if os.IsNotExist(err) {
			if err = os.MkdirAll(root, 0755); err != nil {
				l.Error("cannot create storage root directory", zap.Error(err), zap.String("root", root))
				return errRoot
			}

			return nil
		}

		l.Error("cannot get info about storage root path", zap.Error(err), zap.String("root", root))

		return errRoot
	} else if !info.IsDir() {
		l.Error("storage root path is not a directory", zap.String("root", root))
		return errRoot
	}

	return nil
}
