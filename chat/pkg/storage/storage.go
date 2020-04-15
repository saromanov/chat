package storage

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/saromanov/experiments/chat/pkg/config"
	"github.com/saromanov/experiments/chat/pkg/models"
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

// GetUserByID provides getting user by id
func (s *Storage) GetUserByID(id int64) (*models.User, error) {
	var user *models.User
	tx := s.db.MustBegin()
	if err := tx.Get(&user, "SELECT * FROM users WHERE id = $1", id); err != nil {
		return nil, errors.Wrap(err, "unable to get users")
	}
	return user, nil
}
