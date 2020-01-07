package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/blackplayerten/IdealVisual_backend/account"
	"github.com/blackplayerten/IdealVisual_backend/api"
	"github.com/blackplayerten/IdealVisual_backend/config"
	"github.com/blackplayerten/IdealVisual_backend/database"
	"github.com/blackplayerten/IdealVisual_backend/post"
	"github.com/blackplayerten/IdealVisual_backend/session"
)

// nolint:funlen
// TODO: to application entity
func main() {
	l, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("cannot create logger: %s", err)
	}

	cfg, err := config.NewConfig()
	if err != nil {
		l.Fatal("cannot create config", zap.Error(err))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	sessionSvc := session.New(l.With(zap.Namespace("session")), cfg.Session)
	if err := sessionSvc.ConnectWithCtx(ctx); err != nil {
		cancel()
		l.Fatal("cannot create session service", zap.Error(err))
	}

	cancel()
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)

	ds := database.NewDatabase(cfg.DataSource)
	if err := ds.ConnectToDB(ctx); err != nil {
		cancel()
		l.Fatal("cannot create data source",
			zap.Error(err),
			zap.String("connString", cfg.DataSource.ConnString),
			zap.String("name", cfg.DataSource.Name),
		)
	}

	cancel()
	l.Info("successfully connected to data source",
		zap.String("connString", cfg.DataSource.ConnString),
		zap.String("name", cfg.DataSource.Name),
	)

	n, err := ds.ApplyMigrations()
	if err != nil {
		l.Fatal("cannot apply migrations", zap.Error(err))
	} else if n != 0 {
		l.Info("applied migrations", zap.Int("count", n))
	}

	accountSvc := account.New(ds)
	postSvc := post.New(ds)

	srv, err := api.New(cfg.Server, l.With(zap.Namespace("server")), sessionSvc, accountSvc, postSvc)
	if err != nil {
		l.Fatal("cannot create server", zap.Error(err))
	}

	idleConnsClosed := make(chan struct{})

	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		sig := <-sigs

		l.Info("shutting down server", zap.String("got", sig.String()))

		if err := srv.Shutdown(); err != nil {
			l.Error("server shutdown error", zap.Error(err))
		}

		if err := sessionSvc.Close(); err != nil {
			l.Error("session service close error", zap.Error(err))
		}

		if err := ds.Close(); err != nil {
			l.Error("data source close error", zap.Error(err))
		}

		close(idleConnsClosed)
	}()

	l.Info("starting server at", zap.String("addr", cfg.Server.HTTP.Addr))

	if err := srv.ListenAndServe(); err != nil {
		l.Fatal("server serve error", zap.Error(err))
	}

	<-idleConnsClosed
}
