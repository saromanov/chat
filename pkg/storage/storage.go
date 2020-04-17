package storage

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

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
