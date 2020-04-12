package storage

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/saromanov/experiments/chat/pkg/models"
	"github.com/saromanov/experiments/chat/pkg/config"
)

// Storage defines handlling of storage
type Storage struct {
	db *sqlx.DB
}

// New provides definition of
func New(c *config.Project) (*Storage, error) {
	db, err := sqlx.Connect("postgres", "user=foo dbname=bar sslmode=disable")
	if err != nil {
		return nil, err
	}

	return &Storage{
		db: db,
	}, nil
}

func (s *Storage) AddUser(u *models.User) {
	tx := s.db.MustBegin()
    tx.MustExec("INSERT INTO users (first_name, last_name, email, password) VALUES ($1, $2, $3, $4)", u.FirstName, u.LastName, u.Email, u.Password)
    tx.Commit()
}
