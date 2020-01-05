package session

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/gomodule/redigo/redis"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
)

var (
	ErrKeyNotFound = errors.New("key not found")
)

type Service struct {
	cfg *Config

	l    *zap.Logger
	pool *redis.Pool
}

func New(l *zap.Logger, cfg *Config) *Service {
	return &Service{
		cfg: cfg,

		l: l,
		pool: &redis.Pool{
			MaxIdle:     500,
			IdleTimeout: 240 * time.Second,
			MaxActive:   1000,
			Wait:        true,
			Dial: func() (redis.Conn, error) {
				return redis.DialURL("redis://" + cfg.ConnString + "/" + strconv.Itoa(cfg.Database))
			},
		},
	}
}

func (s *Service) Connect() error {
	return s.ConnectWithCtx(context.Background())
}

func (s *Service) ConnectWithCtx(ctx context.Context) error {
	conn, err := s.pool.GetContext(ctx)
	if err != nil {
		s.l.Error(
			"cannot open connection",
			zap.Int("db", s.cfg.Database), zap.String("connString", s.cfg.ConnString), zap.Error(err),
		)

		return err
	}
	defer conn.Close()

	if _, err = conn.Do("PING"); err != nil {
		s.l.Error(
			"cannot ping connection",
			zap.Int("db", s.cfg.Database), zap.String("connString", s.cfg.ConnString), zap.Error(err),
		)

		return err
	}

	s.l.Info("successfully connected", zap.Int("db", s.cfg.Database), zap.String("connString", s.cfg.ConnString))

	return nil
}

func (s *Service) Close() {
	s.pool.Close()
}

func (s *Service) Create(userID uint64) (string, error) {
	conn := s.pool.Get()
	defer conn.Close()

	var sID string

	for {
		sID = uuid.NewV4().String()
		if res, err := redis.String(conn.Do(
			"SET", sID, userID, "NX", "EX", int64(s.cfg.Expiration.Round(1*time.Second)))); err != nil {
			return "", err
		} else if res == "OK" {
			break
		}
	}

	s.l.Info("created new session", zap.String("sessionID", sID), zap.Uint64("userID", userID))

	return sID, nil
}

func (s *Service) Get(sessionID string) (uint64, error) {
	conn := s.pool.Get()
	defer conn.Close()

	res, err := redis.Uint64(conn.Do("GET", sessionID))
	if err != nil {
		if err == redis.ErrNil {
			return 0, ErrKeyNotFound
		}

		return 0, err
	}

	return res, nil
}

func (s *Service) Delete(sessionID string) error {
	conn := s.pool.Get()
	defer conn.Close()

	if _, err := conn.Do("DEL", sessionID); err != nil {
		return err
	}

	s.l.Info("session deleted", zap.String("sessionID", sessionID))

	return nil
}
