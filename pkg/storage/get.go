package storage

import (
	"github.com/pkg/errors"
	"github.com/saromanov/chat/pkg/models"
)

// GetUser provides adding of the new user
func (s *Storage) GetUser(id int64) (*models.User, error) {
	var user *models.User
	if err := s.db.Get(&user, "SELECT * FROM users WHERE id = $1", id); err != nil {
		return nil, errors.Wrap(err, "unable to get user by id")
	}
	return user, nil
}
