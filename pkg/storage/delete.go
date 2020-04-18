package storage

import (
	"github.com/pkg/errors"
)

// DeleteUser provides removing of the user
func (s *Storage) DeleteUser(id int64) error {
	tx := s.db.MustBegin()
	tx.MustExec("DELETE from users WHERE id = $1", id)
	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "unable to commit changes")
	}
	return nil
}
