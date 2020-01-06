package account

import (
	"context"

	"go.uber.org/zap"

	"github.com/blackplayerten/IdealVisual_backend/database"
)

type Service struct {
	cfg *Config

	l  *zap.Logger
	db *database.Database
}

func New(l *zap.Logger, cfg *Config) *Service {
	return &Service{
		cfg: cfg,

		l:  l,
		db: database.NewDatabase(cfg.Database),
	}
}

func (s *Service) Connect() error {
	return s.ConnectWithCtx(context.Background())
}

func (s *Service) ConnectWithCtx(ctx context.Context) error {
	if err := s.db.ConnectToDB(ctx); err != nil {
		return err
	}

	s.l.Info("successfully connected to db",
		zap.String("connString", s.cfg.Database.ConnString),
		zap.String("db", s.cfg.Database.Name),
	)

	n, err := s.db.ApplyMigrations()
	if err == nil && n != 0 {
		s.l.Info("applied migrations", zap.Int("count", n))
	}

	return err
}

func (s *Service) Close() error {
	return s.db.Close()
}
