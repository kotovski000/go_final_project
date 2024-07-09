package task

func GetTask(s TaskStorage, id string) (Task, error) {
	var task Task
	err := s.db.QueryRow("SELECT * FROM scheduler WHERE id = ?", id).Scan(
		&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat,
	)
	return task, err
}
