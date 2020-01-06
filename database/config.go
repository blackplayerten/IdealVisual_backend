package database

type Config struct {
	ConnString string
	Name       string

	MigrationsFolder string
}
