package storage

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/pressly/goose"
	"github.com/saromanov/chat/pkg/config"
)

// Storage defines handlling of storage
type Storage struct {
	db *sqlx.DB
}

// New provides definition of
func New(c *config.Project) (*Storage, error) {
	url := fmt.Sprintf("sslmode=disable user=%s dbname=%s host=%s port=5432 password=%s", c.DatabaseUser, c.DatabaseName, c.DatabaseHost, c.DatabasePassword)
	fmt.Println(url)
	db, err := sqlx.Connect("postgres", fmt.Sprintf("sslmode=disable user=%s dbname=%s host=%s port=5432 password=%s", c.DatabaseUser, c.DatabaseName, c.DatabaseHost, c.DatabasePassword))
	if err != nil {
		return nil, err
	}

	return &Storage{
		db: db,
	}, nil
}

// Prepare provides applying of helpful methods after init of db
func (s *Storage) Prepare() error {

	if err := s.applyMigrations(); err != nil {
		return errors.Wrap(err, "apply migrations was failed")
	}
	return nil
}

// applyMigrations provides applying of available migrations
func (s *Storage) applyMigrations() error {

	path := os.Getenv("MIGRATIONS_PATH")
	if path == "" {
		return errors.New("variable MIGRATIONS_PATH is not set")
	}

	if err := goose.Up(s.db.DB, path); err != nil {
		return errors.Wrap(err, "unable to apply migrations")
	}

	return nil
}
