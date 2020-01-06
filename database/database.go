package database

import (
	"context"
	"regexp"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // postgres driver
	migrate "github.com/rubenv/sql-migrate"
)

type Database struct {
	cfg *Config

	*sqlx.DB

	keyParenthesesRe *regexp.Regexp
}

func NewDatabase(cfg *Config) *Database {
	return &Database{
		cfg: cfg,

		keyParenthesesRe: regexp.MustCompile(`Key \((.*?)\)=`),
	}
}

func (db *Database) ConnectToDB(ctx context.Context) error {
	var err error
	if db.DB, err = sqlx.Open(
		"postgres",
		"postgresql://"+db.cfg.ConnString+"/"+db.cfg.Name+"?sslmode=disable",
	); err != nil {
		return err
	}

	return db.PingContext(ctx)
}

func (db *Database) ApplyMigrations() (int, error) {
	migrations := &migrate.FileMigrationSource{
		Dir: db.cfg.MigrationsFolder,
	}

	return migrate.Exec(db.DB.DB, "postgres", migrations, migrate.Up)
}
