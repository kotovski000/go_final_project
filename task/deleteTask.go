package task

import (
	"errors"
)

func DeleteTask(s TaskStorage, id string) error {
	res, err := s.db.Exec("DELETE FROM scheduler WHERE id =?", id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("задача не найдена")
	}

	return nil
}
