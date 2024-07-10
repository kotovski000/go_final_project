package task

func PutTask(s TaskStorage, updTask Task) error {
	_, err := s.db.Exec("UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?",
		updTask.Date, updTask.Title, updTask.Comment, updTask.Repeat, updTask.ID)
	if err != nil {
		return err
	}
	return nil
}
