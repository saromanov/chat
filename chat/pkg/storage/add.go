package storage

import (
	"github.com/pkg/errors"
	"github.com/saromanov/experiments/chat/pkg/models"
)

// AddUser provides adding of the new user
func (s *Storage) AddUser(u *models.User) error {
	return s.add("INSERT INTO users (first_name, last_name, email, password) VALUES ($1, $2, $3, $4)", u.FirstName, u.LastName, u.Email, u.Password)
}

// AddMessage provides adding of the new message
func (s *Storage) AddMessage(u *models.Message) error {
	return s.add("INSERT INTO messages (body) VALUES ($1)", u.Body)
}

// add is a general method for insert data
func (s *Storage) add(query string, arguments ...interface{}) error {
	tx := s.db.MustBegin()
	tx.MustExec(query, arguments...)
	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "unable to commit changes")
	}
	return nil
}
